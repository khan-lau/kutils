package ksync

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/khan-lau/kutils/container/kcontext"
)

var (
	ErrWaitTimeout     = errors.New("wait timeout")                      // 定义超时错误
	ErrContextCanceled = errors.New("context canceled before countdown") // 定义上下文取消错误
)

type CountDownLatch struct {
	// ctx 是上下文节点, 当其不为nil时, 用于取消倒计时锁
	ctx *kcontext.ContextNode

	// count 表示剩余的任务数. 使用 int32 是为了兼容 sync/atomic 包的操作
	count int32

	// done 是一个信号通道. 当 count 变为 0 时, 我们会关闭这个通道
	// Go 语言中, 读取一个已关闭的通道会立即返回, 利用这个特性实现“广播唤醒”
	done chan struct{}

	once sync.Once // 新增: 保证 done 只关闭一次
}

// NewCountDownLatch 初始化一个指定计数的倒计时锁
func NewCountDownLatch(delta int) *CountDownLatch {
	return &CountDownLatch{
		ctx:   nil,
		count: int32(delta),
		done:  make(chan struct{}),
	}
}

func NewCountDownLatchWithCtx(ctx *kcontext.ContextNode, delta int) *CountDownLatch {
	return &CountDownLatch{
		ctx:   ctx,
		count: int32(delta),
		done:  make(chan struct{}),
	}
}

// CountDown 减少计数器的值
func (that *CountDownLatch) CountDown() {
	// 1. 原子操作：安全地将 count 减 1。
	// atomic.AddInt32 是线程安全的，即使 1000 个协程同时调用，
	// 每一笔减法都不会丢失，且返回值是减完后的最终结果。
	newCount := atomic.AddInt32(&that.count, -1)

	// 2. 边界检查：如果计数已经归零，则触发唤醒。
	if newCount == 0 {
		// 关闭通道。这是线程安全的核心点：
		// 所有的 Wait() 协程都在阻塞等待读取 that.done。
		// 一旦 close，所有 Wait 的协程都会同时收到信号并解除阻塞。
		that.closeDone()
	} else if newCount < 0 {
		// 防止调用次数超过初始设置的 delta，通常在调试时很有用
		// panic("CountDown called too many times")
	}
}

// AllDown 直接清理计数器，不等待所有任务完成
func (that *CountDownLatch) AllDown() bool {
	// 直接将 count 设为 0，触发所有 Wait() 协程的解除阻塞
	atomic.StoreInt32(&that.count, 0)
	that.closeDone()
	return true
}

// Wait 阻塞等待，直到计数器归零
func (that *CountDownLatch) Wait() {
	// 这里的操作是线程安全的：
	// 在 Go 中，从通道接收数据是阻塞操作。
	// 当通道被 close 时，<-that.done 会立即解除阻塞并返回零值。
	if that.ctx == nil {
		<-that.done
		return
	}

	select {
	case <-that.ctx.Context().Done():
		return
	case <-that.done:
		return
	}
}

// WaitWithTimeout 阻塞等待，直到计数器归零或超时
func (that *CountDownLatch) WaitWithTimeout(timeout time.Duration) error {
	// 创建一个定时器，经过指定时间后会向 timer.C 发送当前时间
	timer := time.NewTimer(timeout)
	defer timer.Stop() // 必须停止定时器以释放资源

	if that.ctx == nil {
		// select 语句会监听多个通道，谁先有数据（或被关闭）就执行谁
		select {
		case <-that.done:
			// 场景 A：计数器先归零，通道被关闭。
			// 返回 nil 表示成功。
			return nil
		case <-timer.C:
			// 场景 B：时间先到了，done 通道还没关。
			// 返回超时错误。
			return ErrWaitTimeout
		}
	} else {
		select {
		case <-that.ctx.Context().Done():
			return ErrContextCanceled
		case <-timer.C:
			return ErrWaitTimeout

		case <-that.done:
			return nil
		}
	}

}

// GetCount 获取当前还有多少个任务没完成
func (that *CountDownLatch) GetCount() int {
	// 使用原子加载，确保读取到的是内存中最新的值
	return int(atomic.LoadInt32(&that.count))
}

// 统一关闭入口, CountDown / AllDown 都走这里
func (that *CountDownLatch) closeDone() {
	that.once.Do(func() { close(that.done) })
}
