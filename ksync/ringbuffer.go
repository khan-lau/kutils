package ksync

import (
	"errors"
	"sync/atomic"
	"time"
)

// SPSC Ring Buffer - Single Producer Single Consumer
// 容量必须是 2 的幂次方: 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096, 8192, 16384, 32768, 65536
type RingBuffer[T any] struct {
	_      [8]uint64 // 填充，防止与前面的字段伪共享
	head   uint64
	_      [8]uint64 // 填充，确保 head 和 tail 在不同 Cache Line
	tail   uint64
	_      [8]uint64
	buffer []T
	mask   uint64
}

var (
	ErrRingBufferSize = errors.New("ring buffer size must be a power of 2") // ErrRingBufferSize 无效的缓冲区大小错误
)

// NewRingBuffer 创建 RingBuffer，size 必须是 2 的幂
func NewRingBuffer[T any](size uint64) (*RingBuffer[T], error) {
	if size == 0 || (size&(size-1)) != 0 {
		return nil, ErrRingBufferSize
	}
	return &RingBuffer[T]{
		buffer: make([]T, size),
		mask:   size - 1,
	}, nil
}

// ====================== 非阻塞接口 ======================

// AsyncEnqueue 非阻塞写入单个元素
// 返回 true 表示成功，false 表示队列已满
func (that *RingBuffer[T]) AsyncEnqueue(item T) bool {
	head := atomic.LoadUint64(&that.head)
	tail := atomic.LoadUint64(&that.tail)

	nextTail := (tail + 1) & that.mask
	if nextTail == head { // 满
		return false
	}

	that.buffer[tail] = item
	atomic.StoreUint64(&that.tail, nextTail)
	return true
}

// AsyncEnqueueBatch 非阻塞批量写入
// 返回实际成功写入的数量（可能小于 len(items)）
func (that *RingBuffer[T]) AsyncEnqueueBatch(items ...T) int {
	if len(items) == 0 {
		return 0
	}

	head := atomic.LoadUint64(&that.head)
	tail := atomic.LoadUint64(&that.tail)

	available := (head - tail - 1) & that.mask
	if available == 0 {
		return 0
	}

	n := len(items)
	if uint64(n) > available {
		n = int(available)
	}
	// 批量写入 + 清晰的环形处理
	start := int(tail & that.mask)

	// 情况1: 不跨边界（最常见）
	if start+n <= len(that.buffer) {
		copy(that.buffer[start:start+n], items[:n])
	} else {
		// 情况2: 跨边界（环绕写入）
		n1 := len(that.buffer) - start
		copy(that.buffer[start:], items[:n1])
		copy(that.buffer[0:], items[n1:n])
	}

	atomic.StoreUint64(&that.tail, (tail+uint64(n))&that.mask)
	return n
}

// Dequeue 非阻塞读取单个元素
// 返回 (item, true) 成功，(zero, false) 队列为空
func (that *RingBuffer[T]) AsyncDequeue() (T, bool) {
	head := atomic.LoadUint64(&that.head)
	tail := atomic.LoadUint64(&that.tail)

	if head == tail {
		var zero T
		return zero, false
	}

	item := that.buffer[head]
	atomic.StoreUint64(&that.head, (head+1)&that.mask)
	return item, true
}

// DequeueBatch 非阻塞批量读取
// 返回实际读取的数量（最多不超过 len(dst)）
func (that *RingBuffer[T]) AsyncDequeueBatch(max int) ([]T, int) {
	head := atomic.LoadUint64(&that.head)
	tail := atomic.LoadUint64(&that.tail)

	available := (tail - head) & that.mask
	if available == 0 {
		return nil, 0
	}

	if uint64(max) > available {
		max = int(available)
	}

	// 始终分配新内存并拷贝
	result := make([]T, max)
	start := int(head & that.mask)

	// 不跨边界
	if start+max <= len(that.buffer) {
		copy(result, that.buffer[start:start+max])
	} else {
		// 跨边界：需要 copy
		n1 := len(that.buffer) - start
		copy(result[0:n1], that.buffer[start:])
		copy(result[n1:], that.buffer[0:max-n1])
	}

	atomic.StoreUint64(&that.head, (head+uint64(max))&that.mask)
	return result, max
}

// ====================== 阻塞接口 ======================

// Enqueue 阻塞写入单个元素，直到成功或超时
// timeout <= 0 表示永久阻塞
func (that *RingBuffer[T]) Enqueue(timeout time.Duration, item T) bool {
	deadline := time.Time{}
	if timeout > 0 {
		deadline = time.Now().Add(timeout)
	}

	for {
		if that.AsyncEnqueue(item) {
			return true
		}

		if timeout > 0 && time.Now().After(deadline) {
			return false
		}

		// 短暂让出 CPU，降低忙轮询开销
		time.Sleep(1 * time.Microsecond)
	}
}

// EnqueueBatchBlocking 阻塞批量写入，直到全部写入或超时
// 返回实际写入数量（超时情况下可能 < len(items)）
func (that *RingBuffer[T]) EnqueueBatchBlocking(timeout time.Duration, items ...T) int {
	if len(items) == 0 {
		return 0
	}
	deadline := time.Time{}
	if timeout > 0 {
		deadline = time.Now().Add(timeout)
	}
	written := 0
	for written < len(items) {
		n := that.AsyncEnqueueBatch(items[written:]...)
		written += n

		if written == len(items) {
			return written
		}

		if timeout > 0 && time.Now().After(deadline) {
			return written
		}

		time.Sleep(1 * time.Microsecond)
	}
	return written
}

// DequeueBlocking 阻塞读取单个元素，直到成功或超时
func (that *RingBuffer[T]) DequeueBlocking(timeout time.Duration) (T, bool) {
	deadline := time.Time{}
	if timeout > 0 {
		deadline = time.Now().Add(timeout)
	}

	for {
		if item, ok := that.AsyncDequeue(); ok {
			return item, true
		}

		if timeout > 0 && time.Now().After(deadline) {
			var zero T
			return zero, false
		}

		time.Sleep(1 * time.Microsecond)
	}
}

// DequeueBatchBlocking 阻塞批量读取，直到读到至少1个或超时
// 返回实际读取的数量
func (that *RingBuffer[T]) DequeueBatchBlocking(max int, timeout time.Duration) ([]T, int) {
	if max == 0 {
		return nil, 0
	}

	deadline := time.Time{}
	if timeout > 0 {
		deadline = time.Now().Add(timeout)
	}
	for {
		if data, n := that.AsyncDequeueBatch(max); n > 0 {
			return data, n
		}
		if timeout > 0 && time.Now().After(deadline) {
			return nil, 0
		}
		time.Sleep(1 * time.Microsecond)
	}
}

// ====================== 查询接口 ======================

func (that *RingBuffer[T]) Len() uint64 {
	head := atomic.LoadUint64(&that.head)
	tail := atomic.LoadUint64(&that.tail)
	return (tail - head) & that.mask
}

func (that *RingBuffer[T]) Cap() uint64 {
	return uint64(len(that.buffer))
}

func (that *RingBuffer[T]) IsEmpty() bool {
	return atomic.LoadUint64(&that.head) == atomic.LoadUint64(&that.tail)
}

func (that *RingBuffer[T]) IsFull() bool {
	head := atomic.LoadUint64(&that.head)
	tail := atomic.LoadUint64(&that.tail)
	return ((tail + 1) & that.mask) == head
}
