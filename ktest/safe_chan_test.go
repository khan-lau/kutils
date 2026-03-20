package ktest

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/khan-lau/kutils/katomic"
)

func TestSafeChannel_RaceCondition(t *testing.T) {
	const (
		producerCount = 100
		sendCount     = 1000
		bufferSize    = 100
	)

	sc := katomic.NewSafeChannel[int](bufferSize)
	var wg sync.WaitGroup
	var totalSent, totalOutdated, totalClosed int64

	// 1. 启动 100 个生产者
	for i := 0; i < producerCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < sendCount; j++ {
				err := sc.Send(j)
				switch err {
				case nil:
					atomic.AddInt64(&totalSent, 1)
				case katomic.ErrGenerationOutdated:
					atomic.AddInt64(&totalOutdated, 1)
				case katomic.ErrChannelClosed:
					atomic.AddInt64(&totalClosed, 1)
				}
				// 模拟业务耗时，增加竞态发生的概率
				if j%10 == 0 {
					time.Sleep(time.Microsecond)
				}
			}
		}(i)
	}

	// 2. 启动一个干扰者，频繁关闭并重置
	go func() {
		for i := 0; i < 5; i++ {
			time.Sleep(5 * time.Millisecond)
			sc.Close()
			time.Sleep(2 * time.Millisecond)
			sc.Reset()
		}
	}()

	// 3. 启动一个消费者，排空数据
	go func() {
		for {
			ch, _ := sc.Ch()
			if ch == nil {
				return
			}
			for range ch {
				// 只是排空数据
			}
			time.Sleep(time.Millisecond)
		}
	}()

	wg.Wait()
	sc.Close()

	t.Logf("测试完成:\n成功发送: %d\n版本过期拦截: %d\n关闭拒绝发送: %d",
		totalSent, totalOutdated, totalClosed)
}

func BenchmarkNativeChannel(b *testing.B) {
	ch := make(chan int, 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		select {
		case ch <- i:
		default:
			<-ch
			ch <- i
		}
	}
}

func BenchmarkSafeChannel_Send(b *testing.B) {
	sc := katomic.NewSafeChannel[int](1000)

	// 启动一个后台协程不断消费，防止阻塞
	done := make(chan struct{})
	go func() {
		ch, _ := sc.Ch()
		for range ch {
			// 消费数据
		}
		close(done)
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 这里的 Send 会在缓冲区满时自动阻塞，直到后台协程腾出空间
		sc.Send(i)
	}
	b.StopTimer() // 计时结束

	sc.Close()
	<-done // 等待数据排空
}
