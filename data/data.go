package data

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"io"

	// "github.com/google/brotli/go/cbrotli"
	"github.com/andybalholm/brotli"

	"khan-lau/kutil/katomic"
)

////////////////////////////////////////////////////////////////

// 生成器结构体，包含一个原子整数
type Generator struct {
	counter *katomic.Uint32
}

func NewGenerator(val uint32) *Generator {
	return &Generator{counter: katomic.NewUint32(val)}
}

func (g *Generator) SetCounter(v uint32) {
	g.counter.Store(v)
}

// 生成方法，返回当前计数器的值，并原子地加一
func (g *Generator) Generate() uint32 {
	return g.counter.Add(1)
}

////////////////////////////////////////////////////////////////

////////////////////////////////////////////////////////////////

// 压缩数据
func Compress(data []uint8) ([]uint8, error) {
	buf := bytes.NewBuffer([]uint8{})
	// 创建一个新的写入器，使用默认的压缩级别
	w, err := flate.NewWriter(buf, flate.DefaultCompression)
	if err != nil {
		return nil, err
	}

	// 写入数据并关闭写入器
	_, err = w.Write(data)
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// 解压缩数据块
func Uncompress(data []uint8) ([]uint8, error) {
	buf := bytes.NewBuffer(data)
	// 创建一个新的读取器，用于解压缩数据
	r := flate.NewReader(buf)
	// 读取解压缩后的数据并关闭读取器
	uncompressed, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	err = r.Close()
	if err != nil {
		return nil, err
	}
	return uncompressed, nil
}

func CompressBr(data []uint8) ([]uint8, error) {
	buf := bytes.NewBuffer([]byte{})
	// 创建一个新的读取器，用于解压缩数据
	// w := brotli.NewWriter(buf)
	w := brotli.NewWriterOptions(buf, brotli.WriterOptions{Quality: 5})

	_, err := w.Write(data)
	if err != nil {
		w.Close()
		return nil, err
	}

	w.Close()
	// 将读取器中的数据复制到新的缓冲区
	return buf.Bytes(), nil
}

func UncompressBr(data []uint8) ([]uint8, error) {
	buf := bytes.NewBuffer(data)
	// 创建一个新的读取器，用于解压缩数据
	r := brotli.NewReader(buf)

	out := bytes.NewBuffer([]byte{})
	// 将读取器中的数据复制到新的缓冲区
	io.Copy(out, r)
	return out.Bytes(), nil
}

func GZip(data []uint8) ([]uint8, error) {
	buf := bytes.NewBuffer([]byte{})
	// 创建一个新的读取器，用于解压缩数据
	w := gzip.NewWriter(buf)

	_, err := w.Write(data)
	if err != nil {
		w.Close()
		return nil, err
	}

	w.Close()
	// 将读取器中的数据复制到新的缓冲区
	return buf.Bytes(), nil
}

func UnGZip(data []uint8) ([]uint8, error) {
	buf := bytes.NewBuffer(data)
	// 创建一个新的读取器，用于解压缩数据
	r, err := gzip.NewReader(buf)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	out := bytes.NewBuffer([]byte{})
	// 将读取器中的数据复制到新的缓冲区
	io.Copy(out, r)
	return out.Bytes(), nil
}

func Zip(data []uint8) ([]uint8, error) {
	buf := bytes.NewBuffer([]byte{})
	// 创建一个新的读取器，用于解压缩数据
	w := zlib.NewWriter(buf)

	_, err := w.Write(data)
	if err != nil {
		w.Close()
		return nil, err
	}

	w.Close()
	// 将读取器中的数据复制到新的缓冲区
	return buf.Bytes(), nil
}

func UnZip(data []uint8) ([]uint8, error) {
	buf := bytes.NewBuffer(data)
	// 创建一个新的读取器，用于解压缩数据
	r, err := zlib.NewReader(buf)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	out := bytes.NewBuffer([]byte{})
	// 将读取器中的数据复制到新的缓冲区
	io.Copy(out, r)
	return out.Bytes(), nil
}

////////////////////////////////////////////////////////////////

////////////////////////////////////////////////////////////////

// 校验和计算
func CheckSum(data []byte) uint16 {
	var (
		sum    uint32
		length int = len(data)
		index  int
	)
	//以每16位为单位进行求和，直到所有的字节全部求完或者只剩下一个8位字节（如果剩余一个8位字节说明字节数为奇数个）
	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}
	//如果字节数为奇数个，要加上最后剩下的那个8位字节
	if length > 0 {
		sum += uint32(data[index])
	}
	//加上高16位进位的部分
	sum += (sum >> 16)
	//别忘了返回的时候先求反
	return uint16(^sum)
}

////////////////////////////////////////////////////////////////
