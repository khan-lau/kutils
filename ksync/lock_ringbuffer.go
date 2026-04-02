package ksync

import (
	"sync"
)

// LockedRingBuffer 有锁版本的 SPSC Ring Buffer
// 使用 Mutex + Cond 实现高效阻塞，适合中低吞吐或需要规范阻塞语义的场景
type LockedRingBuffer[T any] struct {
	mu        sync.Mutex
	condFull  *sync.Cond // 队列满时，生产者在此等待
	condEmpty *sync.Cond // 队列空时，消费者在此等待
	buffer    []T        // 数据存储区
	mask      uint64     // 容量减一，用于取模
	head      uint64     // 头部offset
	tail      uint64     // 尾部offset
	closed    bool       // 闭塞状态
}

func NewLockedRingBuffer[T any](size uint64) (*LockedRingBuffer[T], error) {
	if size == 0 || (size&(size-1)) != 0 {
		return nil, ErrRingBufferSize
	}

	rb := &LockedRingBuffer[T]{
		buffer: make([]T, size),
		mask:   size - 1,
	}
	rb.condFull = sync.NewCond(&rb.mu)  // 生产者等待：队列满时等待空间
	rb.condEmpty = sync.NewCond(&rb.mu) // 消费者等待：队列空时等待数据
	return rb, nil
}

// ====================== 非阻塞接口 ======================

// AsyncEnqueue 非阻塞单条写入
func (that *LockedRingBuffer[T]) AsyncEnqueue(item T) bool {
	if !that.mu.TryLock() {
		return false
	}
	defer that.mu.Unlock()

	// 检查队列是否已满或已关闭
	if that.closed || ((that.tail+1)&that.mask) == that.head {
		return false
	}

	// 直接使用 mask 实现环形索引
	that.buffer[that.tail&that.mask] = item // ← 显式 & mask
	that.tail = (that.tail + 1) & that.mask // ← 环形递增
	that.condEmpty.Signal()                 // 通知消费者有新数据
	return true
}

// AsyncEnqueueBatch 非阻塞批量写入
func (that *LockedRingBuffer[T]) AsyncEnqueueBatch(items ...T) int {
	if len(items) == 0 {
		return 0
	}

	if !that.mu.TryLock() {
		return 0
	}
	defer that.mu.Unlock()

	if that.closed {
		return 0
	}

	available := (that.head - that.tail - 1) & that.mask
	if available == 0 {
		return 0
	}

	n := len(items)
	if uint64(n) > available {
		n = int(available)
	}

	// 批量写入 + 清晰的环形处理
	start := int(that.tail & that.mask)

	// 情况1: 不跨边界（最常见）
	if start+n <= len(that.buffer) {
		copy(that.buffer[start:start+n], items[:n])
	} else {
		// 情况2: 跨边界（环绕写入）
		n1 := len(that.buffer) - start
		copy(that.buffer[start:], items[:n1])
		copy(that.buffer[0:], items[n1:n])
	}

	// 更新 tail（环形递增）
	that.tail = (that.tail + uint64(n)) & that.mask
	that.condEmpty.Signal() // 通知消费者有新数据
	return n
}

// AsyncDequeue 非阻塞单条读取
func (that *LockedRingBuffer[T]) AsyncDequeue() (T, bool) {
	if !that.mu.TryLock() {
		var zero T
		return zero, false
	}
	defer that.mu.Unlock()

	if that.closed || that.head == that.tail {
		var zero T
		return zero, false
	}

	item := that.buffer[that.head]
	that.head = (that.head + 1) & that.mask
	that.condFull.Signal() // 通知生产者空间已释放
	return item, true
}

// AsyncDequeueBatch 非阻塞批量读取
//
//   - 参数
//
//   - max: 最大读取数量
//
//   - 返回
//
//   - 读取到的数据
//
//   - 读取到的数量
func (that *LockedRingBuffer[T]) AsyncDequeueBatch(max int) ([]T, int) {
	if !that.mu.TryLock() {
		return nil, 0
	}
	defer that.mu.Unlock()

	available := (that.tail - that.head) & that.mask
	if available == 0 {
		return nil, 0
	}

	if uint64(max) > available {
		max = int(available)
	}

	start := int(that.head & that.mask)

	// 不跨边界：零拷贝返回子切片
	if start+max <= len(that.buffer) {
		result := that.buffer[start : start+max]
		that.head = (that.head + uint64(max)) & that.mask
		that.condFull.Signal() // 通知生产者空间已释放
		return result, max
	}

	// 跨边界：需要 copy
	result := make([]T, max)
	n1 := len(that.buffer) - start
	copy(result[0:n1], that.buffer[start:])
	copy(result[n1:], that.buffer[0:max-n1])

	that.head = (that.head + uint64(max)) & that.mask
	that.condFull.Signal() // 通知生产者空间已释放
	return result, max
}

// ====================== 阻塞接口（正常 Cond 实现） ======================

// Enqueue 阻塞写入单个元素，直到成功或超时
func (that *LockedRingBuffer[T]) Enqueue(item T) bool {
	that.mu.Lock()
	defer that.mu.Unlock()

	if that.closed {
		return false
	}

	for ((that.tail + 1) & that.mask) == that.head { // 队列已满时等待
		if that.closed { // 超时或已关闭则退出
			return false
		}
		that.condFull.Wait() // 等待空间释放
	}

	that.buffer[that.tail&that.mask] = item // 存入元素
	that.tail = (that.tail + 1) & that.mask // 更新尾指针
	that.condEmpty.Signal()                 // 通知消费者有新数据
	return true
}

// EnqueueBatch 阻塞批量写入，直到全部写入或超时
func (that *LockedRingBuffer[T]) EnqueueBatch(items ...T) int {
	if len(items) == 0 {
		return 0
	}

	that.mu.Lock()
	defer that.mu.Unlock()

	if that.closed {
		return 0
	}

	written := 0
	for written < len(items) {
		available := (that.head - that.tail - 1) & that.mask

		if available > 0 {
			canWrite := int(available)
			remain := len(items) - written
			if canWrite > remain {
				canWrite = remain
			}

			start := int(that.tail & that.mask)

			// 优化点：使用 copy 替代循环赋值
			if start+canWrite <= len(that.buffer) {
				copy(that.buffer[start:start+canWrite], items[written:written+canWrite])
			} else {
				n1 := len(that.buffer) - start
				copy(that.buffer[start:], items[written:written+n1])
				copy(that.buffer[0:], items[written+n1:written+canWrite])
			}

			that.tail = (that.tail + uint64(canWrite)) & that.mask
			written += canWrite
			that.condEmpty.Signal()
			continue
		}

		// 队列关闭或超时, 直接退出, 否则等待
		if that.closed {
			return written
		}
		that.condFull.Wait()
	}

	return written
}

// Dequeue 阻塞读取单个元素
func (that *LockedRingBuffer[T]) Dequeue() (T, bool) {
	that.mu.Lock()
	defer that.mu.Unlock()

	// 若队列为空
	for that.head == that.tail {
		if that.closed {
			var zero T
			return zero, false
		}
		that.condEmpty.Wait() // 等待生产者入队数据
	}

	item := that.buffer[that.head]
	that.head = (that.head + 1) & that.mask
	that.condFull.Signal() // 通知生产者空间已释放
	return item, true
}

// DequeueBatch 从环形缓冲区中批量取出元素
//
// 如果缓冲区中没有数据，会阻塞等待直到有数据可用、超时或缓冲区关闭。
// 最多读取 max 个元素，实际读取数量取决于缓冲区当前可用元素数量。
//
// 参数:
//   - max: 最大读取数量，必须大于0
//   - timeout: 超时等待时间，<=0 表示无限等待
//
// 返回:
//   - []T: 取出的元素数组
//   - int: 实际取出的元素数量
func (that *LockedRingBuffer[T]) DequeueBatch(max int) ([]T, int) {
	if max <= 0 {
		return nil, 0
	}

	that.mu.Lock()
	defer that.mu.Unlock()

	for {
		if that.closed {
			return nil, 0
		}

		available := (that.tail - that.head) & that.mask // 计算可用元素数
		if available > 0 {
			if uint64(max) > available {
				max = int(available) // 调整读取数量不超过可用数
			}
			return that.dequeueBatchInternal(max)
		}
		that.condEmpty.Wait() // 等待数据
	}
}

// ====================== 控制接口 ======================

// Close 关闭队列并唤醒所有等待者
func (that *LockedRingBuffer[T]) Close() {
	that.mu.Lock()
	that.closed = true
	that.condFull.Broadcast()  // 唤醒所有等待队列满的生产者
	that.condEmpty.Broadcast() // 唤醒所有等待队列空的消费者
	that.mu.Unlock()
}

func (that *LockedRingBuffer[T]) Len() int {
	that.mu.Lock()
	defer that.mu.Unlock()
	return int((that.tail - that.head) & that.mask)
}

func (that *LockedRingBuffer[T]) Cap() uint64 {
	return uint64(len(that.buffer))
}

func (that *LockedRingBuffer[T]) IsClosed() bool {
	that.mu.Lock()
	defer that.mu.Unlock()
	return that.closed
}

func (that *LockedRingBuffer[T]) dequeueBatchInternal(max int) ([]T, int) {
	start := int(that.head & that.mask) // 计算读取起始位置
	if start+max <= len(that.buffer) {
		// 无需跨边界，直接切片返回
		result := that.buffer[start : start+max]
		that.head = (that.head + uint64(max)) & that.mask
		that.condFull.Signal()
		return result, max
	}

	// 跨边界情况：需要从尾部和头部两段拼接
	result := make([]T, max)
	n1 := len(that.buffer) - start           // 尾部可取长度
	copy(result[0:n1], that.buffer[start:])  // 复制尾部数据
	copy(result[n1:], that.buffer[0:max-n1]) // 复制头部数据

	that.head = (that.head + uint64(max)) & that.mask
	that.condFull.Signal()
	return result, max
}
