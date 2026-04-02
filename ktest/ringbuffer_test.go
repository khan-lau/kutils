package ktest

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/khan-lau/kutils/ksync"
)

// ============================================================
// 吞吐量测试：测量每秒处理的消息数量
// ============================================================

// TestThroughput_RingBuffer 测试无锁 RingBuffer 吞吐量
func TestThroughput_RingBuffer_Safe(t *testing.T) {
	rb, _ := ksync.NewRingBuffer[int](65536)
	const testDuration = 3 * time.Second

	var produced, consumed atomic.Int64
	var stopProducer atomic.Bool
	var wg sync.WaitGroup

	// 1. 启动生产者
	wg.Add(1)
	go func() {
		defer wg.Done()
		for !stopProducer.Load() {
			// 防御：如果队列满，必须重试直到成功，确保 produced 计数真实
			for !rb.AsyncEnqueue(1) {
				if stopProducer.Load() {
					return
				}
				runtime.Gosched() // 让出 CPU，避免死循环锁死
			}
			produced.Add(1)
		}
	}()

	// 2. 启动消费者
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if _, ok := rb.AsyncDequeue(); ok {
				consumed.Add(1)
			} else {
				// 如果生产者已经停止，且队列已经读空，则退出
				if stopProducer.Load() && rb.IsEmpty() {
					return
				}
				runtime.Gosched()
			}
		}
	}()

	// 3. 运行指定时间
	time.Sleep(testDuration)

	// 4. 优雅关闭流程
	stopProducer.Store(true) // 首先通知生产者停止写入
	wg.Wait()                // 等待生产者退出，且等待消费者排空缓冲区

	producedCount := produced.Load()
	consumedCount := consumed.Load()

	// 5. 最终校验
	if producedCount != consumedCount {
		t.Errorf("数据不一致! 生产: %d, 消费: %d, 差异: %d", producedCount, consumedCount, producedCount-consumedCount)
	}

	throughput := consumedCount / int64(testDuration.Seconds())
	t.Logf("[RingBuffer] 吞吐量: %d 条/秒, 总处理: %d", throughput, consumedCount)
}

// TestThroughput_LockedRingBuffer 测试有锁 LockedRingBuffer 吞吐量
func TestThroughput_LockedRingBuffer_Safe(t *testing.T) {
	rb, _ := ksync.NewLockedRingBuffer[int](65536)
	const testDuration = 3 * time.Second

	var produced, consumed atomic.Int64
	var stopProducer atomic.Bool // 明确只控制生产者
	var wg sync.WaitGroup

	// 1. 生产者：必须保证数据“实打实”存入
	wg.Add(1)
	go func() {
		defer wg.Done()
		for !stopProducer.Load() {
			// 防御点：在高频测试中，不要只试一次。
			// 如果 TryLock 失败或队列满，应持续重试直到 stop 信号转真
			for !rb.AsyncEnqueue(1) {
				if stopProducer.Load() {
					return
				}
				// 稍微让出 CPU，避免 TryLock 疯狂自旋导致锁更加难以获取
				runtime.Gosched()
			}
			produced.Add(1)
		}
	}()

	// 2. 消费者：必须排空缓冲区
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if _, ok := rb.AsyncDequeue(); ok {
				consumed.Add(1)
			} else {
				// 防御点：只有当生产者停了，且队列确实空了，才能退出
				if stopProducer.Load() && rb.Len() == 0 {
					return
				}
				runtime.Gosched()
			}
		}
	}()

	// 3. 运行测试
	time.Sleep(testDuration)

	// 4. 优雅关闭流程
	stopProducer.Store(true) // 先通知停止生产
	wg.Wait()                // 等待所有协程按序处理完数据退出

	producedCount := produced.Load()
	consumedCount := consumed.Load()

	// 5. 最终对账
	if producedCount != consumedCount {
		t.Errorf("[LockedRingBuffer] 数据丢失! 生产: %d, 消费: %d, 差额: %d",
			producedCount, consumedCount, producedCount-consumedCount)
	}

	throughput := consumedCount / int64(testDuration.Seconds())
	t.Logf("[LockedRingBuffer] 吞吐量: %d 条/秒, 总处理: %d", throughput, consumedCount)
}

// TestThroughput_Channel 测试原生 Channel 吞吐量
func TestThroughput_Channel(t *testing.T) {
	ch := make(chan int, 65536)
	const testDuration = 3 * time.Second

	var produced, consumed atomic.Int64
	var stop atomic.Bool

	// 生产者
	go func() {
		for !stop.Load() {
			select {
			case ch <- 1:
				produced.Add(1)
			default:
				// 队列满，跳过
			}
		}
		close(ch)
	}()

	// 消费者
	go func() {
		for {
			select {
			case _, ok := <-ch:
				if !ok {
					return
				}
				consumed.Add(1)
			default:
				if stop.Load() {
					return
				}
			}
		}
	}()

	// 运行指定时间
	time.Sleep(testDuration)
	stop.Store(true)
	time.Sleep(100 * time.Millisecond)

	producedCount := produced.Load()
	consumedCount := consumed.Load()
	throughput := consumedCount / int64(testDuration.Seconds())

	t.Logf("[Channel] 吞吐量: %d 条/秒, 生产: %d, 消费: %d", throughput, producedCount, consumedCount)
}

// // TestThroughput_MPSC_RingBuffer MPSC 场景吞吐量测试
// func TestThroughput_MPSC_RingBuffer(t *testing.T) {
// 	rb, _ := ksync.NewRingBuffer[int](65536)
// 	const testDuration = 3 * time.Second
// 	const producers = 4
//
// 	var produced, consumed atomic.Int64
// 	var stop atomic.Bool
// 	var mu sync.Mutex
//
// 	// 多个生产者
// 	for p := 0; p < producers; p++ {
// 		go func() {
// 			for !stop.Load() {
// 				mu.Lock()
// 				ok := rb.AsyncEnqueue(1)
// 				mu.Unlock()
// 				if ok {
// 					produced.Add(1)
// 				}
// 			}
// 		}()
// 	}
//
// 	// 单消费者
// 	go func() {
// 		for !stop.Load() {
// 			if _, ok := rb.AsyncDequeue(); ok {
// 				consumed.Add(1)
// 			}
// 		}
// 	}()
//
// 	time.Sleep(testDuration)
// 	stop.Store(true)
// 	time.Sleep(100 * time.Millisecond)
//
// 	producedCount := produced.Load()
// 	consumedCount := consumed.Load()
// 	throughput := consumedCount / int64(testDuration.Seconds())
//
// 	t.Logf("[RingBuffer MPSC] 吞吐量: %d 条/秒, 生产: %d, 消费: %d", throughput, producedCount, consumedCount)
// }

// TestThroughput_MPSC_LockedRingBuffer MPSC 场景吞吐量测试
func TestThroughput_MPSC_LockedRingBuffer(t *testing.T) {
	rb, _ := ksync.NewLockedRingBuffer[int](65536)
	const testDuration = 3 * time.Second
	const producers = 4

	var produced, consumed atomic.Int64
	var stop atomic.Bool

	// 多个生产者
	for p := 0; p < producers; p++ {
		go func() {
			for !stop.Load() {
				if rb.AsyncEnqueue(1) {
					produced.Add(1)
				}
			}
		}()
	}

	// 单消费者
	go func() {
		for !stop.Load() {
			if _, ok := rb.AsyncDequeue(); ok {
				consumed.Add(1)
			}
		}
	}()

	time.Sleep(testDuration)
	stop.Store(true)
	time.Sleep(100 * time.Millisecond)

	producedCount := produced.Load()
	consumedCount := consumed.Load()
	throughput := consumedCount / int64(testDuration.Seconds())

	t.Logf("[LockedRingBuffer MPSC] 吞吐量: %d 条/秒, 生产: %d, 消费: %d", throughput, producedCount, consumedCount)
}

// TestThroughput_MPSC_Channel MPSC 场景吞吐量测试
func TestThroughput_MPSC_Channel(t *testing.T) {
	ch := make(chan int, 65536)
	const testDuration = 3 * time.Second
	const producers = 4

	var produced, consumed atomic.Int64
	var stop atomic.Bool

	// 多个生产者
	for p := 0; p < producers; p++ {
		go func() {
			for !stop.Load() {
				select {
				case ch <- 1:
					produced.Add(1)
				default:
				}
			}
		}()
	}

	// 单消费者
	go func() {
		for !stop.Load() {
			select {
			case _, ok := <-ch:
				if ok {
					consumed.Add(1)
				}
			default:
			}
		}
	}()

	time.Sleep(testDuration)
	stop.Store(true)
	time.Sleep(100 * time.Millisecond)

	producedCount := produced.Load()
	consumedCount := consumed.Load()
	throughput := consumedCount / int64(testDuration.Seconds())

	t.Logf("[Channel MPSC] 吞吐量: %d 条/秒, 生产: %d, 消费: %d", throughput, producedCount, consumedCount)
}

// TestThroughput_Comparison 综合对比测试
func TestThroughput_Comparison(t *testing.T) {
	t.Log("========== SPSC 吞吐量对比 ==========")

	t.Run("RingBuffer", TestThroughput_RingBuffer_Safe)
	t.Run("LockedRingBuffer", TestThroughput_LockedRingBuffer_Safe)
	t.Run("Channel", TestThroughput_Channel)

	t.Log("\n========== MPSC 吞吐量对比 (4生产者) ==========")

	t.Run("LockedRingBuffer_MPSC", TestThroughput_MPSC_LockedRingBuffer)
	t.Run("Channel_MPSC", TestThroughput_MPSC_Channel)
}

// ============================================================
// 批量操作吞吐量测试
// ============================================================

// TestThroughput_Batch_RingBuffer 批量操作吞吐量测试
func TestThroughput_Batch_RingBuffer(t *testing.T) {
	rb, _ := ksync.NewRingBuffer[int](65536)
	const testDuration = 3 * time.Second
	const batchSize = 64

	items := make([]int, batchSize)
	for i := range items {
		items[i] = i
	}

	var produced, consumed atomic.Int64
	var stopProducer atomic.Bool
	var wg sync.WaitGroup

	// 生产者 - 批量写入
	wg.Add(1)
	go func() {
		defer wg.Done()
		for !stopProducer.Load() {
			n := rb.AsyncEnqueueBatch(items...)
			if n > 0 {
				produced.Add(int64(n))
			}
			if n < batchSize {
				runtime.Gosched()
			}
		}
	}()

	// 消费者 - 批量读取
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			_, n := rb.AsyncDequeueBatch(batchSize)
			if n > 0 {
				consumed.Add(int64(n))
			} else {
				if stopProducer.Load() && rb.IsEmpty() {
					return
				}
				runtime.Gosched()
			}
		}
	}()

	time.Sleep(testDuration)
	stopProducer.Store(true)
	wg.Wait()

	producedCount := produced.Load()
	consumedCount := consumed.Load()

	if producedCount != consumedCount {
		t.Errorf("[RingBuffer Batch] 数据不一致! 生产: %d, 消费: %d", producedCount, consumedCount)
	}

	throughput := consumedCount / int64(testDuration.Seconds())
	t.Logf("[RingBuffer Batch] 吞吐量: %d 条/秒, 总处理: %d, 批次大小: %d", throughput, consumedCount, batchSize)
}

// TestThroughput_Batch_LockedRingBuffer 批量操作吞吐量测试
func TestThroughput_Batch_LockedRingBuffer(t *testing.T) {
	rb, _ := ksync.NewLockedRingBuffer[int](65536)
	const testDuration = 3 * time.Second
	const batchSize = 64

	items := make([]int, batchSize)
	for i := range items {
		items[i] = i
	}

	var produced, consumed atomic.Int64
	var stopProducer atomic.Bool
	var wg sync.WaitGroup

	// 生产者 - 批量写入
	wg.Add(1)
	go func() {
		defer wg.Done()
		for !stopProducer.Load() {
			n := rb.AsyncEnqueueBatch(items...)
			if n > 0 {
				produced.Add(int64(n))
			}
			if n < batchSize {
				runtime.Gosched()
			}
		}
	}()

	// 消费者 - 批量读取
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			_, n := rb.AsyncDequeueBatch(batchSize)
			if n > 0 {
				consumed.Add(int64(n))
			} else {
				if stopProducer.Load() && rb.Len() == 0 {
					return
				}
				runtime.Gosched()
			}
		}
	}()

	time.Sleep(testDuration)
	stopProducer.Store(true)
	wg.Wait()

	producedCount := produced.Load()
	consumedCount := consumed.Load()

	if producedCount != consumedCount {
		t.Errorf("[LockedRingBuffer Batch] 数据不一致! 生产: %d, 消费: %d", producedCount, consumedCount)
	}

	throughput := consumedCount / int64(testDuration.Seconds())
	t.Logf("[LockedRingBuffer Batch] 吞吐量: %d 条/秒, 总处理: %d, 批次大小: %d", throughput, consumedCount, batchSize)
}

// TestThroughput_Batch_Channel 批量操作吞吐量测试
func TestThroughput_Batch_Channel(t *testing.T) {
	ch := make(chan int, 65536)
	const testDuration = 3 * time.Second
	const batchSize = 64

	items := make([]int, batchSize)
	for i := range items {
		items[i] = i
	}

	var produced, consumed atomic.Int64
	var stopProducer atomic.Bool
	var wg sync.WaitGroup

	// 生产者 - 模拟批量写入（逐条发送）
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(ch)
		for !stopProducer.Load() {
			sent := 0
			for _, item := range items {
				select {
				case ch <- item:
					sent++
				default:
					break
				}
			}
			if sent > 0 {
				produced.Add(int64(sent))
			}
			if sent < batchSize {
				runtime.Gosched()
			}
		}
	}()

	// 消费者 - 批量读取
	wg.Add(1)
	go func() {
		defer wg.Done()
		batch := make([]int, 0, batchSize)
		for {
			select {
			case item, ok := <-ch:
				if !ok {
					if len(batch) > 0 {
						consumed.Add(int64(len(batch)))
					}
					return
				}
				batch = append(batch, item)
				if len(batch) >= batchSize {
					consumed.Add(int64(len(batch)))
					batch = batch[:0]
				}
			default:
				if stopProducer.Load() && len(ch) == 0 {
					if len(batch) > 0 {
						consumed.Add(int64(len(batch)))
					}
					return
				}
				runtime.Gosched()
			}
		}
	}()

	time.Sleep(testDuration)
	stopProducer.Store(true)
	wg.Wait()

	// producedCount := produced.Load()
	consumedCount := consumed.Load()

	throughput := consumedCount / int64(testDuration.Seconds())
	t.Logf("[Channel Batch] 吞吐量: %d 条/秒, 总处理: %d, 批次大小: %d", throughput, consumedCount, batchSize)
}

// TestThroughput_Batch_Comparison 批量操作综合对比测试
func TestThroughput_Batch_Comparison(t *testing.T) {
	t.Log("========== 批量操作吞吐量对比 (batchSize=64) ==========")

	t.Run("RingBuffer_Batch", TestThroughput_Batch_RingBuffer)
	t.Run("LockedRingBuffer_Batch", TestThroughput_Batch_LockedRingBuffer)
	t.Run("Channel_Batch", TestThroughput_Batch_Channel)
}

// TestThroughput_Batch_DifferentSizes 不同批次大小对比测试
func TestThroughput_Batch_DifferentSizes(t *testing.T) {
	batchSizes := []int{16, 32, 64, 128, 256}

	t.Log("========== 不同批次大小吞吐量对比 ==========")

	for _, size := range batchSizes {
		t.Logf("\n--- 批次大小: %d ---", size)

		// RingBuffer
		t.Run(fmt.Sprintf("RingBuffer_Batch_%d", size), func(t *testing.T) {
			testBatchThroughputRingBuffer(t, size)
		})

		// LockedRingBuffer
		t.Run(fmt.Sprintf("LockedRingBuffer_Batch_%d", size), func(t *testing.T) {
			testBatchThroughputLockedRingBuffer(t, size)
		})

		// Channel
		t.Run(fmt.Sprintf("Channel_Batch_%d", size), func(t *testing.T) {
			testBatchThroughputChannel(t, size)
		})
	}
}

func testBatchThroughputRingBuffer(t *testing.T, batchSize int) {
	rb, _ := ksync.NewRingBuffer[int](65536)
	const testDuration = 2 * time.Second

	items := make([]int, batchSize)
	for i := range items {
		items[i] = i
	}

	var produced, consumed atomic.Int64
	var stopProducer atomic.Bool
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for !stopProducer.Load() {
			n := rb.AsyncEnqueueBatch(items...)
			if n > 0 {
				produced.Add(int64(n))
			}
			if n < batchSize {
				runtime.Gosched()
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			_, n := rb.AsyncDequeueBatch(batchSize)
			if n > 0 {
				consumed.Add(int64(n))
			} else {
				if stopProducer.Load() && rb.IsEmpty() {
					return
				}
				runtime.Gosched()
			}
		}
	}()

	time.Sleep(testDuration)
	stopProducer.Store(true)
	wg.Wait()

	throughput := consumed.Load() / int64(testDuration.Seconds())
	t.Logf("[RingBuffer] 批次=%d, 吞吐量: %d 条/秒", batchSize, throughput)
}

func testBatchThroughputLockedRingBuffer(t *testing.T, batchSize int) {
	rb, _ := ksync.NewLockedRingBuffer[int](65536)
	const testDuration = 2 * time.Second

	items := make([]int, batchSize)
	for i := range items {
		items[i] = i
	}

	var produced, consumed atomic.Int64
	var stopProducer atomic.Bool
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for !stopProducer.Load() {
			n := rb.AsyncEnqueueBatch(items...)
			if n > 0 {
				produced.Add(int64(n))
			}
			if n < batchSize {
				runtime.Gosched()
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			_, n := rb.AsyncDequeueBatch(batchSize)
			if n > 0 {
				consumed.Add(int64(n))
			} else {
				if stopProducer.Load() && rb.Len() == 0 {
					return
				}
				runtime.Gosched()
			}
		}
	}()

	time.Sleep(testDuration)
	stopProducer.Store(true)
	wg.Wait()

	throughput := consumed.Load() / int64(testDuration.Seconds())
	t.Logf("[LockedRingBuffer] 批次=%d, 吞吐量: %d 条/秒", batchSize, throughput)
}

func testBatchThroughputChannel(t *testing.T, batchSize int) {
	ch := make(chan int, 65536)
	const testDuration = 2 * time.Second

	items := make([]int, batchSize)
	for i := range items {
		items[i] = i
	}

	var produced, consumed atomic.Int64
	var stopProducer atomic.Bool
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(ch)
		for !stopProducer.Load() {
			sent := 0
			for _, item := range items {
				select {
				case ch <- item:
					sent++
				default:
					break
				}
			}
			if sent > 0 {
				produced.Add(int64(sent))
			}
			if sent < batchSize {
				runtime.Gosched()
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		batch := make([]int, 0, batchSize)
		for {
			select {
			case item, ok := <-ch:
				if !ok {
					if len(batch) > 0 {
						consumed.Add(int64(len(batch)))
					}
					return
				}
				batch = append(batch, item)
				if len(batch) >= batchSize {
					consumed.Add(int64(len(batch)))
					batch = batch[:0]
				}
			default:
				if stopProducer.Load() && len(ch) == 0 {
					if len(batch) > 0 {
						consumed.Add(int64(len(batch)))
					}
					return
				}
				runtime.Gosched()
			}
		}
	}()

	time.Sleep(testDuration)
	stopProducer.Store(true)
	wg.Wait()

	throughput := consumed.Load() / int64(testDuration.Seconds())
	t.Logf("[Channel] 批次=%d, 吞吐量: %d 条/秒", batchSize, throughput)
}
