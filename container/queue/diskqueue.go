package queue

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// 日志工具从 github.com/nsqio/nsq/internal/lg 移植

const (
	DISKQUEUE_TAG = "diskqueue"

	DEFAULT_READ_BUFFER_SIZE  = 256 * 1024      // 256KB 优化：增大缓冲区，减少 read 系统调用
	DEFAULT_WRITE_BUFFER_SIZE = 1 * 1024 * 1024 // 写缓冲大小
)

///////////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////////

type LogLevel int

const (
	DEBUG = LogLevel(1)
	INFO  = LogLevel(2)
	WARN  = LogLevel(3)
	ERROR = LogLevel(4)
	FATAL = LogLevel(5)
)

type AppLogFunc func(lvl LogLevel, tag string, f string, args ...any)

// String 返回日志级别的字符串表示。
//
// 返回值：
//   - string: 日志级别名称（DEBUG/INFO/WARNING/ERROR/FATAL）
func (l LogLevel) String() string {
	switch l {
	case 1:
		return "DEBUG"
	case 2:
		return "INFO"
	case 3:
		return "WARNING"
	case 4:
		return "ERROR"
	case 5:
		return "FATAL"
	}
	panic("invalid LogLevel")
}

///////////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////////

// DiskQueue 基于文件系统的 FIFO 队列，提供持久化的消息存储能力。
//
// 数据格式：
//
//	每个消息以 4 字节大端序长度头 + 消息体存储，文件按代际命名：{name}.diskqueue.{seq}.dat
//	元数据文件 {name}.diskqueue.meta.dat 记录深度和读写位置，用于重启后恢复。
//
// 使用范例：
//
//	q, err := New("myqueue", "/data/queue", 100*1024*1024, 1, 1024*1024,
//		time.Second, 2500, 2*time.Second, logf)
//	if err != nil { panic(err) }
//	defer q.Close()
//
//	// 生产者
//	go func() {
//		for i := 0; i < 1000; i++ {
//			q.Put([]byte(fmt.Sprintf("msg-%d", i)))
//		}
//	}()
//
//	// 消费者
//	for msg := range q.ReadChan() {
//		process(msg)
//	}
type Interface interface {
	// Put 写入消息到队列尾部。
	//
	// 参数：
	//   - data: 待写入的消息体
	//
	// 返回值：
	//   - error: 写入失败的错误信息
	Put([]byte) error
	// ReadChan 返回只读 channel，接收队列头部消息（消费即移除）。
	//
	// 返回值：
	//   - <-chan []byte: 只读消息 channel
	ReadChan() <-chan []byte
	// PeekChan 返回只读 channel，接收队列头部消息但不移除（可重复读取）。
	//
	// 返回值：
	//   - <-chan []byte: 只读窥视 channel
	PeekChan() <-chan []byte
	// Close 关闭队列并持久化元数据。
	//
	// 返回值：
	//   - error: 关闭失败的错误信息
	Close() error
	// Delete 删除队列及所有数据文件。
	//
	// 返回值：
	//   - error: 删除失败的错误信息
	Delete() error
	// Depth 返回当前队列深度（待处理消息数）。
	//
	// 返回值：
	//   - int64: 队列中尚未被消费的消息数量
	Depth() int64
	// Empty 清空队列所有数据。
	//
	// 返回值：
	//   - error: 清空失败的错误信息
	Empty() error
}

// diskQueue 实现了一个基于文件系统的 FIFO 队列
type diskQueue struct {
	// 64位原子变量需放在结构体最前面，确保32位平台上的内存对齐
	// 运行时状态（也会持久化到磁盘）
	readPos      int64
	writePos     int64
	readFileNum  int64
	writeFileNum int64
	depth        int64 // 队列长度

	sync.RWMutex

	// 实例化时确定的元数据
	name                string
	dataPath            string
	maxBytesPerFile     int64 // 创建后不可更改
	maxBytesPerFileRead int64
	minMsgSize          int32
	maxMsgSize          int32
	syncEvery           int64         // 每次 fsync 之间的写入次数
	syncTimeout         time.Duration // fsync 的时间间隔
	exitFlag            int32
	needSync            bool

	// 记录已从磁盘读取但尚未通过 readChan 投递给消费者的位置
	nextReadPos     int64
	nextReadFileNum int64

	readFile    *os.File
	writeFile   *os.File
	reader      *bufio.Reader
	writeWriter *bufio.Writer

	// 通过 ReadChan() 暴露
	readChan chan []byte

	// 通过 PeekChan() 暴露
	peekChan chan []byte

	// 内部 channel
	depthChan         chan int64
	writeChan         chan []byte
	writeResponseChan chan error
	emptyChan         chan int
	emptyResponseChan chan error
	exitChan          chan int
	exitSyncChan      chan int

	logf AppLogFunc
}

// New 创建并初始化一个基于文件系统的 FIFO 队列。
//
// 从磁盘恢复上一次的运行状态（读写位置、队列深度），
// 并启动后台事件循环 goroutine（ioLoop）统一管理读写。
//
// 参数：
//
//	@param name: 队列名称，用于命名数据文件和元数据文件
//	@param dataPath: 数据文件存放目录
//	@param maxBytesPerFile: 单个数据文件的最大字节数，超出后自动滚动到新文件
//	@param minMsgSize: 允许的最小消息体大小（字节）
//	@param maxMsgSize: 允许的最大消息体大小（字节）
//	@param syncEvery: 每写入多少次后执行一次 fsync
//	@param syncTimeout: fsync 的超时时间间隔
//	@param logf: 日志记录回调函数
//
// 返回值：
//
//	@returns Interface: 队列操作接口
//	@returns error: 创建失败的错误信息
func New(name string, dataPath string, maxBytesPerFile int64,
	minMsgSize int32, maxMsgSize int32,
	syncEvery int64, syncTimeout time.Duration, logf AppLogFunc,
) (Interface, error) {

	if err := os.MkdirAll(dataPath, os.ModePerm); err != nil {
		return nil, err
	}
	d := diskQueue{
		name:              name,
		dataPath:          dataPath,
		maxBytesPerFile:   maxBytesPerFile,
		minMsgSize:        minMsgSize,
		maxMsgSize:        maxMsgSize,
		readChan:          make(chan []byte),
		peekChan:          make(chan []byte),
		depthChan:         make(chan int64),
		writeChan:         make(chan []byte, 1024),
		writeResponseChan: make(chan error),
		emptyChan:         make(chan int),
		emptyResponseChan: make(chan error),
		exitChan:          make(chan int),
		exitSyncChan:      make(chan int),
		syncEvery:         syncEvery,
		syncTimeout:       syncTimeout,
		logf:              logf,
	}

	// 此时无需加锁，因为没有任何其他代码能触及此实例
	err := d.retrieveMetaData()
	if err != nil && !os.IsNotExist(err) {
		if d.logf != nil {
			d.logf(ERROR, DISKQUEUE_TAG, "DISKQUEUE(%s) failed to retrieveMetaData - %v", d.name, err)
		}
	}

	go d.ioLoop()
	return &d, nil
}

// Depth 返回当前队列深度（待处理消息数）。
//
// 如果 ioLoop 已退出，返回最后一次记录的值。
//
// 返回值：
//   - int64: 队列中尚未被消费的消息数量
func (d *diskQueue) Depth() int64 {
	depth, ok := <-d.depthChan
	if !ok {
		// ioLoop 已退出
		depth = d.depth
	}
	return depth
}

// Read 以回调方式读取消息，消费即移除。
//
// 在独立的 goroutine 中运行，回调内可阻塞，不影响 ioLoop。
// 当队列关闭或回调返回非 nil 错误时退出。
func (d *diskQueue) Read(callback func([]byte) error) {
	for msg := range d.readChan {
		if err := callback(msg); err != nil {
			return
		}
	}
}

// ReadChan 返回用于读取数据的只读 channel。
//
// 消息由 ioLoop 从磁盘读取后通过此 channel 投递。
// 读取后消息会从队列中移除。
//
// 注意：
//   - 该 channel 有缓冲（1024），消费速度慢时可能阻塞 ioLoop
//
// 返回值：
//   - <-chan []byte: 只读消息 channel
func (d *diskQueue) ReadChan() <-chan []byte {
	return d.readChan
}

// PeekChan 返回用于窥视数据的只读 channel。
//
// 与 ReadChan 不同，通过 PeekChan 读取的数据不会从队列中移除，
// 下次 ReadChan 或 PeekChan 仍可读到相同的消息。
//
// 返回值：
//   - <-chan []byte: 只读窥视 channel
func (d *diskQueue) PeekChan() <-chan []byte {
	return d.peekChan
}

// Put 向队列尾部写入一条消息。
//
// 写入操作由 ioLoop 异步处理，当队列正在关闭时返回 "exiting" 错误。
//
// 参数：
//   - data: 待写入的消息体
//
// 返回值：
//   - error: 写入失败的错误信息
func (d *diskQueue) Put(data []byte) error {
	d.RLock()
	if d.exitFlag == 1 {
		d.RUnlock()
		return errors.New("exiting")
	}
	d.RUnlock()

	select {
	case d.writeChan <- data:
		select {
		case err := <-d.writeResponseChan:
			return err
		case <-d.exitChan:
			return errors.New("exiting")
		}
	case <-d.exitChan:
		return errors.New("exiting")
	}
}

// Close 优雅关闭队列。
//
// 先通知 ioLoop 退出，等待其完成当前操作后，
// 再对当前写入文件执行 fsync 并持久化元数据。
//
// 返回值：
//   - error: 关闭失败的错误信息
func (d *diskQueue) Close() error {
	err := d.exit(false)
	if err != nil {
		return err
	}
	return d.sync()
}

// Delete 删除队列及所有数据文件。
//
// 调用 exit 后，ioLoop 退出时清理元数据文件和数据文件。
//
// 返回值：
//   - error: 删除失败的错误信息
func (d *diskQueue) Delete() error {
	return d.exit(true)
}

// exit 执行队列关闭或删除操作。
//
// 设置退出标志，关闭 exitChan 通知 ioLoop 终止，
// 等待 ioLoop 确认退出后清理文件句柄。
//
// 参数：
//   - deleted: true 表示删除队列（删除后不保留数据），
//     false 表示关闭队列（保留数据文件以便下次启动恢复）
func (d *diskQueue) exit(deleted bool) error {
	d.Lock()
	defer d.Unlock()

	if d.exitFlag == 1 {
		return nil
	}
	d.exitFlag = 1

	if d.logf != nil {
		if deleted {
			d.logf(INFO, DISKQUEUE_TAG, "DISKQUEUE(%s) deleting", d.name)
		} else {
			d.logf(INFO, DISKQUEUE_TAG, "DISKQUEUE(%s) closing", d.name)
		}
	}

	close(d.exitChan)
	// 确保 ioLoop 已退出
	<-d.exitSyncChan

	close(d.depthChan)

	if d.readFile != nil {
		d.readFile.Close()
		d.readFile = nil
	}

	if d.writeWriter != nil {
		d.writeWriter = nil
	}
	if d.writeFile != nil {
		d.writeFile.Close()
		d.writeFile = nil
	}

	return nil
}

// Empty 破坏性清空队列中所有待处理数据。
//
// 通过快进读取位置并移除中间文件来实现。
// 操作由 ioLoop 异步处理，当队列正在关闭时返回 "exiting" 错误。
//
// 返回值：
//   - error: 清空失败的错误信息
func (d *diskQueue) Empty() error {
	d.RLock()
	if d.exitFlag == 1 {
		d.RUnlock()
		return errors.New("exiting")
	}
	d.RUnlock()

	if d.logf != nil {
		d.logf(INFO, DISKQUEUE_TAG, "DISKQUEUE(%s) emptying", d.name)
	}

	select {
	case d.emptyChan <- 1:
		select {
		case err := <-d.emptyResponseChan:
			return err
		case <-d.exitChan:
			return errors.New("exiting")
		}
	case <-d.exitChan:
		return errors.New("exiting")
	}
}

// deleteAllFiles 删除所有数据文件和元数据文件。
//
// 先通过 skipToNextRWFile 快进读写位置并清理数据文件，
// 再删除元数据文件。
//
// 返回值：
//   - error: 删除失败的错误信息
func (d *diskQueue) deleteAllFiles() error {
	err := d.skipToNextRWFile()

	innerErr := os.Remove(d.metaDataFileName())
	if innerErr != nil && !os.IsNotExist(innerErr) {
		if d.logf != nil {
			d.logf(ERROR, DISKQUEUE_TAG, "DISKQUEUE(%s) failed to remove metadata file - %v", d.name, innerErr)
		}
		return innerErr
	}

	return err
}

// skipToNextRWFile 将读写位置快进到下一个新文件。
//
// 关闭当前读写文件句柄，删除当前读写文件号之间的所有数据文件，
// 然后将读写位置重置到新的文件号起始处，队列深度归零。
//
// 返回值：
//   - error: 删除文件过程中的错误（累积返回）
func (d *diskQueue) skipToNextRWFile() error {
	var err error

	if d.readFile != nil {
		d.readFile.Close()
		d.readFile = nil
	}

	if d.writeWriter != nil {
		d.writeWriter = nil // 文件将被删除，无需刷盘
	}
	if d.writeFile != nil {
		d.writeFile.Close()
		d.writeFile = nil
	}

	for i := d.readFileNum; i <= d.writeFileNum; i++ {
		fn := d.fileName(i)
		innerErr := os.Remove(fn)
		if innerErr != nil && !os.IsNotExist(innerErr) {
			if d.logf != nil {
				d.logf(ERROR, DISKQUEUE_TAG, "DISKQUEUE(%s) failed to remove data file - %v", d.name, innerErr)
			}
			err = innerErr
		}
	}

	d.writeFileNum++
	d.writePos = 0
	d.readFileNum = d.writeFileNum
	d.readPos = 0
	d.nextReadFileNum = d.writeFileNum
	d.nextReadPos = 0
	d.depth = 0
	return err
}

// readOne 执行单个 []byte 的底层文件系统读取。
//
// 同时按需推进读取位置和滚动文件。若当前文件未打开则自动打开，
// 读取完成后推进 nextReadPos/nextReadFileNum（实际消费时由 moveForward 推进 readPos/readFileNum）。
//
// 返回值：
//   - []byte: 读取到的消息体
//   - error: 读取失败的错误信息
func (d *diskQueue) readOne() ([]byte, error) {
	var err error
	var msgSize int32

	// 如果读写同一文件且有未刷盘的写缓冲数据，先 flush 确保读取可见
	if d.readFileNum == d.writeFileNum && d.writeWriter != nil && d.writeWriter.Buffered() > 0 {
		err = d.writeWriter.Flush()
		if err != nil {
			d.writeWriter = nil
			d.writeFile.Close()
			d.writeFile = nil
			return nil, err
		}
	}

	if d.readFile == nil {
		curFileName := d.fileName(d.readFileNum)
		d.readFile, err = os.OpenFile(curFileName, os.O_RDONLY, 0600)
		if err != nil {
			return nil, err
		}

		if d.logf != nil {
			d.logf(INFO, DISKQUEUE_TAG, "DISKQUEUE(%s) readOne() opened %s", d.name, curFileName)
		}

		if d.readPos > 0 {
			_, err = d.readFile.Seek(d.readPos, 0)
			if err != nil {
				d.readFile.Close()
				d.readFile = nil
				return nil, err
			}
		}

		// 对于"已完成"的文件（即非"当前"文件），maxBytesPerFileRead
		// 应初始化为文件的实际大小，或默认使用 maxBytesPerFile
		d.maxBytesPerFileRead = d.maxBytesPerFile
		if d.readFileNum < d.writeFileNum {
			stat, statErr := d.readFile.Stat()
			if statErr == nil {
				d.maxBytesPerFileRead = stat.Size()
			}
		}

		// // d.reader = bufio.NewReader(d.readFile)
		// 优化：增大缓冲区，减少 read 系统调用
		d.reader = bufio.NewReaderSize(d.readFile, DEFAULT_READ_BUFFER_SIZE) // 256KB
	}

	err = binary.Read(d.reader, binary.BigEndian, &msgSize)
	if err != nil {
		d.readFile.Close()
		d.readFile = nil
		return nil, err
	}

	if msgSize < d.minMsgSize || msgSize > d.maxMsgSize {
		// 此文件已损坏，无法合理确定新消息的起始位置
		d.readFile.Close()
		d.readFile = nil
		return nil, fmt.Errorf("invalid message read size (%d)", msgSize)
	}

	readBuf := make([]byte, msgSize)
	_, err = io.ReadFull(d.reader, readBuf)
	if err != nil {
		d.readFile.Close()
		d.readFile = nil
		return nil, err
	}

	totalBytes := int64(4 + msgSize)

	// 我们只推进 next*，因为尚未将数据投递给消费者
	// （readFileNum, readPos 将在投递后才推进）
	d.nextReadPos = d.readPos + totalBytes
	d.nextReadFileNum = d.readFileNum

	// 只有读取"已完成"的文件时才考虑滚动，
	// 由于我们不知道文件滚动时的大小，用 maxBytesPerFileRead 而非 maxBytesPerFile
	if d.readFileNum < d.writeFileNum && d.nextReadPos >= d.maxBytesPerFileRead {
		if d.readFile != nil {
			d.readFile.Close()
			d.readFile = nil
		}

		d.nextReadFileNum++
		d.nextReadPos = 0
	}

	return readBuf, nil
}

// writeOne 执行单个 []byte 的底层文件系统写入。
//
// 同时按需推进写入位置和滚动文件。若当前文件未打开则自动创建。
// 当文件大小超过 maxBytesPerFile 时自动滚动到下一个文件。
//
// 参数：
//   - data: 待写入的消息体
//
// 返回值：
//   - error: 写入失败的错误信息
func (d *diskQueue) writeOne(data []byte) error {
	var err error

	dataLen := int32(len(data))
	totalBytes := int64(4 + dataLen)

	if dataLen < d.minMsgSize || dataLen > d.maxMsgSize {
		return fmt.Errorf("invalid message write size (%d) minMsgSize=%d maxMsgSize=%d", dataLen, d.minMsgSize, d.maxMsgSize)
	}

	if d.writeFile == nil {
		curFileName := d.fileName(d.writeFileNum)
		d.writeFile, err = os.OpenFile(curFileName, os.O_RDWR|os.O_CREATE, 0600)
		if err != nil {
			return err
		}
		if d.logf != nil {
			d.logf(INFO, DISKQUEUE_TAG, "DISKQUEUE(%s) writeOne() opened %s", d.name, curFileName)
		}

		if d.writePos > 0 {
			_, err = d.writeFile.Seek(d.writePos, 0)
			if err != nil {
				d.writeFile.Close()
				d.writeFile = nil
				return err
			}
		}

		d.writeWriter = bufio.NewWriterSize(d.writeFile, DEFAULT_WRITE_BUFFER_SIZE)
	}

	// 直接写入 writeWriter，跳过中间缓冲，减少一次拷贝 + 消除 binary.Write 反射开销
	var header [4]byte
	binary.BigEndian.PutUint32(header[:], uint32(dataLen))
	_, err = d.writeWriter.Write(header[:])
	if err != nil {
		d.writeWriter = nil
		d.writeFile.Close()
		d.writeFile = nil
		return err
	}

	_, err = d.writeWriter.Write(data)
	if err != nil {
		d.writeWriter = nil
		d.writeFile.Close()
		d.writeFile = nil
		return err
	}

	d.writePos += totalBytes
	d.depth += 1

	// 只要 maxBytesPerFile + maxMsgSize < Int64Max 就不会溢出
	if d.writePos > 0 && d.writePos > d.maxBytesPerFile {
		if d.readFileNum == d.writeFileNum {
			d.maxBytesPerFileRead = d.writePos
		}

		d.writeFileNum++
		d.writePos = 0

		// 每次新开文件时执行 sync
		err = d.sync()
		if err != nil {
			if d.logf != nil {
				d.logf(ERROR, DISKQUEUE_TAG, "DISKQUEUE(%s) failed to sync - %v", d.name, err)
			}
		}

		if d.writeFile != nil {
			d.writeFile.Close()
			d.writeFile = nil
		}
		d.writeWriter = nil
	}

	return err
}

// sync 对当前写入文件执行 fsync 并持久化元数据。
//
// 先对 writeFile 调用 Sync() 确保数据落盘，
// 再调用 persistMetaData() 将运行时状态写入元数据文件。
// 操作成功后清除 needSync 标志。
//
// 返回值：
//   - error: 同步失败的错误信息
func (d *diskQueue) sync() error {
	if d.writeWriter != nil {
		err := d.writeWriter.Flush()
		if err != nil {
			d.writeWriter = nil
			d.writeFile.Close()
			d.writeFile = nil
			return err
		}
	}

	if d.writeFile != nil {
		err := d.writeFile.Sync()
		if err != nil {
			d.writeFile.Close()
			d.writeFile = nil
			return err
		}
	}

	err := d.persistMetaData()
	if err != nil {
		return err
	}

	d.needSync = false
	return nil
}

// retrieveMetaData 从文件系统恢复运行时状态。
//
// 读取元数据文件，恢复深度、读写位置等状态。
// 若实际文件大小大于 writePos（上次未正常 sync），
// 则跳过当前文件，从下一个新文件开始写入。
//
// 返回值：
//   - error: 读取或解析失败的错误信息
func (d *diskQueue) retrieveMetaData() error {
	var f *os.File
	var err error

	fileName := d.metaDataFileName()
	f, err = os.OpenFile(fileName, os.O_RDONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	var depth int64
	_, err = fmt.Fscanf(f, "%d\n%d,%d\n%d,%d\n",
		&depth,
		&d.readFileNum, &d.readPos,
		&d.writeFileNum, &d.writePos)
	if err != nil {
		return err
	}
	d.depth = depth

	d.nextReadFileNum = d.readFileNum
	d.nextReadPos = d.readPos

	// 如果上次 nsqd 关闭时元数据未 sync，
	// 实际文件大小可能大于 writePos。
	// 此时最安全的做法是跳转到下一个文件进行写入，
	// 让读取端尽力从 diskqueue 中恢复元数据中可能已过时的 readPos 之后的剩余消息
	fileName = d.fileName(d.writeFileNum)
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		return err
	}
	fileSize := fileInfo.Size()
	if d.writePos < fileSize {
		if d.logf != nil {
			d.logf(WARN, DISKQUEUE_TAG, "DISKQUEUE(%s) %s metadata writePos %d < file size of %d, skipping to new file", d.name, fileName, d.writePos, fileSize)
		}

		d.writeFileNum += 1
		d.writePos = 0
		if d.writeFile != nil {
			d.writeFile.Close()
			d.writeFile = nil
		}
	}

	return nil
}

// persistMetaData 原子性地将运行时状态写入文件系统。
//
// 先将状态写入临时文件，再通过原子重命名替换原元数据文件，
// 避免写入过程中崩溃导致元数据损坏。
//
// 返回值：
//   - error: 写入失败的错误信息
func (d *diskQueue) persistMetaData() error {
	var f *os.File
	var err error

	fileName := d.metaDataFileName()
	tmpFileName := fmt.Sprintf("%s.%d.tmp", fileName, rand.Int())

	// 写入临时文件
	f, err = os.OpenFile(tmpFileName, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	defer os.Remove(tmpFileName) // 清理临时文件

	_, err = fmt.Fprintf(f, "%d\n%d,%d\n%d,%d\n",
		d.depth,
		d.readFileNum, d.readPos,
		d.writeFileNum, d.writePos)
	if err != nil {
		f.Close()
		return err
	}
	_ = f.Sync()
	f.Close()

	// 原子性重命名
	err = os.Rename(tmpFileName, fileName)
	if err != nil {
		return err
	}
	return nil
}

// metaDataFileName 返回元数据文件的完整路径。
//
// 返回值：
//   - string: 元数据文件路径
func (d *diskQueue) metaDataFileName() string {
	return filepath.Join(d.dataPath, fmt.Sprintf("%s.diskqueue.meta.dat", d.name))
}

// fileName 返回指定代际的数据文件完整路径。
//
// 参数：
//   - fileNum: 文件代际编号
//
// 返回值：
//   - string: 数据文件路径
func (d *diskQueue) fileName(fileNum int64) string {
	return filepath.Join(d.dataPath, fmt.Sprintf("%s.diskqueue.%06d.dat", d.name, fileNum))
}

// checkTailCorruption 检查队列尾部是否损坏，必要时重置深度。
//
// 当读写位置相遇（无待处理消息）但深度不为 0，或读写位置不一致时，
// 判定为元数据损坏，重置深度或跳转到下一个文件。
//
// 参数：
//   - depth: 当前记录的消息深度
func (d *diskQueue) checkTailCorruption(depth int64) {
	if d.readFileNum < d.writeFileNum || d.readPos < d.writePos {
		return
	}

	// 已到达 diskqueue 尾部，如果深度不为 0 说明有异常
	if depth != 0 {
		if depth < 0 {
			if d.logf != nil {
				d.logf(WARN, DISKQUEUE_TAG, "DISKQUEUE(%s) negative depth at tail (%d), metadata corruption, resetting 0...", d.name, depth)
			}
		} else if depth > 0 {
			if d.logf != nil {
				d.logf(WARN, DISKQUEUE_TAG, "DISKQUEUE(%s) positive depth at tail (%d), data loss, resetting 0...", d.name, depth)
			}
		}
		// 强制将深度设为 0
		d.depth = 0
		d.needSync = true
	}

	if d.readFileNum != d.writeFileNum || d.readPos != d.writePos {
		if d.readFileNum > d.writeFileNum {
			if d.logf != nil {
				d.logf(ERROR, DISKQUEUE_TAG, "DISKQUEUE(%s) readFileNum > writeFileNum (%d > %d), corruption, skipping to next writeFileNum and resetting 0...", d.name, d.readFileNum, d.writeFileNum)
			}
		}

		if d.readPos > d.writePos {
			if d.logf != nil {
				d.logf(ERROR, DISKQUEUE_TAG, "DISKQUEUE(%s) readPos > writePos (%d > %d), skipping to next writeFileNum and resetting 0...", d.name, d.readPos, d.writePos)
			}
		}

		_ = d.skipToNextRWFile()
		d.needSync = true
	}
}

// moveForward 在消息被消费后推进读取位置。
//
// 将 readFileNum/readPos 更新为 nextReadFileNum/nextReadPos，
// 深度减 1。若切换到新文件则删除旧文件并设置 needSync。
// 最后检查队列尾部是否损坏。
func (d *diskQueue) moveForward() {
	oldReadFileNum := d.readFileNum
	d.readFileNum = d.nextReadFileNum
	d.readPos = d.nextReadPos
	d.depth -= 1
	// 检查是否需要清理旧文件
	if oldReadFileNum != d.nextReadFileNum {
		// 每次开始读取新文件时 sync
		d.needSync = true

		fn := d.fileName(oldReadFileNum)
		err := os.Remove(fn)
		if err != nil {
			if d.logf != nil {
				d.logf(ERROR, DISKQUEUE_TAG, "DISKQUEUE(%s) failed to Remove %s - %v", d.name, fn, err)
			}
		}
	}

	d.checkTailCorruption(d.depth)
}

// handleReadError 处理读取错误，跳转到下一个读取文件。
//
// 若当前读取文件与写入文件相同，先跳过该写入文件。
// 将异常文件重命名为 .bad 后缀，然后推进到下一个文件继续读取。
func (d *diskQueue) handleReadError() {
	// 跳转到下一个读取文件，将当前异常文件重命名为 .bad
	if d.readFileNum == d.writeFileNum {
		// 如果无法从当前写入文件读取，跳过该文件
		if d.writeWriter != nil {
			d.writeWriter = nil // 丢弃未刷盘数据，文件可能已损坏
		}
		if d.writeFile != nil {
			d.writeFile.Close()
			d.writeFile = nil
		}
		d.writeFileNum++
		d.writePos = 0
	}

	badFn := d.fileName(d.readFileNum)

	badRenameFn := badFn + ".bad"
	if d.logf != nil {
		d.logf(WARN, DISKQUEUE_TAG, "DISKQUEUE(%s) jump to next file and saving bad file as %s", d.name, badRenameFn)
	}

	err := os.Rename(badFn, badRenameFn)
	if err != nil {
		if d.logf != nil {
			d.logf(ERROR, DISKQUEUE_TAG, "DISKQUEUE(%s) failed to rename bad diskqueue file %s to %s", d.name, badFn, badRenameFn)
		}
	}

	d.readFileNum++
	d.readPos = 0
	d.nextReadFileNum = d.readFileNum
	d.nextReadPos = 0

	// 状态发生重大变化，在下一次循环中执行 sync
	d.needSync = true

	d.checkTailCorruption(d.depth)
}

// ioLoop 是 diskQueue 的后台事件循环，通过 go channel（ReadChan）向消费者投递消息。
//
// 工作原理：根据队列是否有可读数据在读写分支间循环切换，
// 并通过 select 阻塞在对应的 go channel 上，等待读或写事件。
// 使用 pendingRead 标记确保 moveForward 仅在消费者实际读取时调用，
// 防止无消费者时文件位置错误推进导致数据丢失。
//
// 这样做的好处是文件的读取与写入都是异步的。
func (d *diskQueue) ioLoop() {
	var dataRead []byte
	var err error
	var count int64
	var r chan []byte
	var p chan []byte
	pendingRead := false // 是否有已读取但尚未投递给消费者的数据

	syncTicker := time.NewTicker(d.syncTimeout)

	for {
		// 不必每次都 sync :)
		if count >= d.syncEvery {
			d.needSync = true
		}

		if d.needSync {
			err = d.sync()
			if err != nil {
				if d.logf != nil {
					d.logf(ERROR, DISKQUEUE_TAG, "DISKQUEUE(%s) failed to sync - %v", d.name, err)
				}
			}
			count = 0
		}

		// 没有待投递的数据时，尝试从文件读取
		if !pendingRead && ((d.readFileNum < d.writeFileNum) || (d.readPos < d.writePos)) {
			if d.nextReadPos == d.readPos {
				dataRead, err = d.readOne()
				if err != nil {
					if d.logf != nil {
						d.logf(ERROR, DISKQUEUE_TAG, "DISKQUEUE(%s) reading at %d of %s - %v", d.name, d.readPos, d.fileName(d.readFileNum), err)
					}
					d.handleReadError()
					continue
				}
				pendingRead = true
			} else {
				// nextReadPos != readPos → readOne 已读取但未投递
				pendingRead = true
			}
		}

		if pendingRead {
			r = d.readChan
			p = d.peekChan
		} else {
			r = nil
			p = nil
		}

		select {
		// Go channel 规范规定：在 select 中对 nil channel 的读写操作会被跳过。
		// 我们仅在队列有数据时才将 r 设为 d.readChan（否则为 nil），
		// 从而在无数据时 select 不会进入读取分支。
		case p <- dataRead:
			// peek: 发送数据但不推进 readPos，下次可再次读取同一数据
		case r <- dataRead:
			count++
			// 每次读取新文件时 moveForward 会设置 needSync 标志
			d.moveForward()
			pendingRead = false
		case d.depthChan <- d.depth:
		case <-d.emptyChan:
			d.emptyResponseChan <- d.deleteAllFiles()
			pendingRead = false
			count = 0
		case dataWrite := <-d.writeChan:
			count++
			d.writeResponseChan <- d.writeOne(dataWrite)
			// 写入后立即尝试读取（仅当无待投递数据时）
			if !pendingRead && ((d.readFileNum < d.writeFileNum) || (d.readPos < d.writePos)) {
				if d.nextReadPos == d.readPos {
					dataRead, err = d.readOne()
					if err == nil {
						pendingRead = true
						r = d.readChan
						p = d.peekChan
					} else {
						if d.logf != nil {
							d.logf(ERROR, DISKQUEUE_TAG, "DISKQUEUE(%s) reading at %d of %s - %v", d.name, d.readPos, d.fileName(d.readFileNum), err)
						}
						d.handleReadError()
					}
				}
			}
		case <-syncTicker.C:
			if count == 0 {
				// 无写入活动时不做 sync
				continue
			}
			d.needSync = true
		case <-d.exitChan:
			goto exit
		}
	}

exit:
	if d.logf != nil {
		d.logf(INFO, DISKQUEUE_TAG, "DISKQUEUE(%s): closing ... ioLoop", d.name)
	}
	syncTicker.Stop()
	d.exitSyncChan <- 1
}
