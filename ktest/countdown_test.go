package ktest

import (
	"sync"
	"testing"
	"time"

	"github.com/khan-lau/kutils/ksync"
)

// TestCountDownLatch_Wait 测试基础的等待/唤醒功能
func TestCountDownLatch_Wait(t *testing.T) {
	count := 5
	latch := ksync.NewCountDownLatch(count)
	var wg sync.WaitGroup

	// 启动 5 个协程分别执行 CountDown
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			time.Sleep(time.Duration(id*10) * time.Millisecond) // 模拟业务耗时
			latch.CountDown()
		}(i)
	}

	// 记录开始时间
	start := time.Now()

	// 主协程等待
	latch.Wait()

	duration := time.Since(start)
	t.Logf("所有任务完成，耗时: %v", duration)

	if latch.GetCount() != 0 {
		t.Errorf("期望计数为 0，实际为 %d", latch.GetCount())
	}
}

// TestCountDownLatch_Timeout 测试超时机制
func TestCountDownLatch_Timeout(t *testing.T) {
	latch := ksync.NewCountDownLatch(1) // 初始计数为 1

	// 情况 1：在计数归零前超时
	err := latch.WaitWithTimeout(100 * time.Millisecond)
	if err != ksync.ErrWaitTimeout {
		t.Errorf("期望得到超时错误，实际得到: %v", err)
	}

	// 情况 2：在超时时间内计数归零
	latch2 := ksync.NewCountDownLatch(1)
	go func() {
		time.Sleep(50 * time.Millisecond)
		latch2.CountDown()
	}()

	err = latch2.WaitWithTimeout(200 * time.Millisecond)
	if err != nil {
		t.Errorf("期望不超时，实际得到错误: %v", err)
	}
}

// TestCountDownLatch_Concurrency 测试并发安全性（Race Detection）
// 运行测试时建议增加 -race 参数：go test -race -v
func TestCountDownLatch_Concurrency(t *testing.T) {
	numGoroutines := 100
	latch := ksync.NewCountDownLatch(numGoroutines)

	// 多个协程同时调用 Wait 和 CountDown
	for i := 0; i < numGoroutines; i++ {
		go func() {
			latch.CountDown()
		}()
		go func() {
			_ = latch.WaitWithTimeout(500 * time.Millisecond)
		}()
	}

	latch.Wait()
}
