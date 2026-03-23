package katomic

import (
	"errors"
	"sync/atomic"
)

var (
	ErrChannelNotClosed   = errors.New("reset allowed only in closed state")
	ErrChannelClosed      = errors.New("send on closed channel")
	ErrGenerationOutdated = errors.New("generation outdated, send rejected")
	ErrChannelNil         = errors.New("channel is nil")
	ErrChannelFull        = errors.New("channel is full")
)

// SafeChannel 是一个高性能、Panic 免疫的并发通道封装。
// 它使用原子状态机替代了传统的读写锁，以换取接近原生 Channel 的吞吐量。
//
// 设计哲学：
//  1. 100% 杜绝向已关闭通道发送或重复关闭引发的 Panic。
//  2. 仅允许在已关闭状态下进行重置(Reset)，确保生命周期的严谨性。
//
// 注意事项：

type SafeChannel[T any] struct {
	ch     chan T
	closed uint32 // 0: 运行中, 1: 已关闭
	gen    uint64 // 代数/版本号，每次 Reset 自增
	size   int
}

// NewAtomicSafeChannel 创建高性能安全通道
func NewSafeChannel[T any](size int) *SafeChannel[T] {
	return &SafeChannel[T]{
		ch:   make(chan T, size),
		size: size,
		gen:  1, // 初始代数为 1
	}
}

// Reset 仅在已关闭状态下工作 (无锁)
func (that *SafeChannel[T]) Reset() error {
	if that == nil {
		return ErrChannelNil
	}

	// 核心逻辑：使用 CAS 保证只有一个人能重置成功
	// 只有当前状态是 1 (Closed) 时，才将其改为 0 (Running)
	if !atomic.CompareAndSwapUint32(&that.closed, 1, 0) {
		// 如果当前是 0，说明还没关，不允许重置
		return ErrChannelNotClosed
	}

	// 1. 先更新代数，让所有旧周期的发送者感知到变化并失效
	atomic.AddUint64(&that.gen, 1)

	// 2. 替换 Channel, 此时已经成功将状态切回 0，所有 Send 的原子检查都会通过。
	// 重新创建内部 channel。
	that.ch = make(chan T, that.size)
	return nil
}

// Send 保持极致的原子性能 (无锁)
func (that *SafeChannel[T]) Send(data T) (err error) {
	if that == nil || that.ch == nil {
		return ErrChannelNil
	}

	// 1. 记录发起发送时的版本号
	currentGen := atomic.LoadUint64(&that.gen)

	// 2. 检查关闭状态
	if atomic.LoadUint32(&that.closed) == 1 {
		return ErrChannelClosed
	}

	// 2. 拦截 Panic
	defer func() {
		if r := recover(); r != nil {
			err = ErrChannelClosed
		}
	}()

	// 3. 核心：二次校验代数
	// 如果在检查状态和执行发送之间发生了 Reset，这里会拦截
	if atomic.LoadUint64(&that.gen) != currentGen {
		return ErrGenerationOutdated
	}

	that.ch <- data
	return nil
}

// AsyncSend 保持极致的原子性能 (无锁)，但为非阻塞, 可能返回 ErrChannelFull 错误。
func (that *SafeChannel[T]) AsyncSend(data T) (err error) {
	if that == nil || that.ch == nil {
		return ErrChannelNil
	}

	// 1. 记录发起发送时的版本号
	currentGen := atomic.LoadUint64(&that.gen)

	// 2. 检查关闭状态
	if atomic.LoadUint32(&that.closed) == 1 {
		return ErrChannelClosed
	}

	// 3. 拦截 Panic（仍然保留）
	defer func() {
		if r := recover(); r != nil {
			err = ErrChannelClosed
		}
	}()

	// 4. 核心：二次校验代数
	if atomic.LoadUint64(&that.gen) != currentGen {
		return ErrGenerationOutdated
	}

	// 改为 select + default，非阻塞
	select {
	case that.ch <- data:
		return nil
	default:
		return ErrChannelFull // 或自定义错误
	}
}

// 安全的Close函数, 确保只关闭一次，且是原子操作
func (that *SafeChannel[T]) Close() error {
	if that == nil || that.ch == nil {
		return ErrChannelNil
	}

	// 使用 Swap 确保只关闭一次，且是原子操作
	if atomic.SwapUint32(&that.closed, 1) == 0 {
		close(that.ch)
	}
	return nil
}

// Status 检查是否已关闭 (原子操作，无锁)
func (that *SafeChannel[T]) Status() bool {
	if that == nil {
		return true
	}
	return atomic.LoadUint32(&that.closed) == 1
}

// Ch 获取原生 channel 和当前代数用于 select 或 range
//
// 返回: 原生 channel 和当前代数,代数在每次 Close 时递增
func (that *SafeChannel[T]) Ch() (<-chan T, uint64) {
	if that == nil {
		return nil, 0
	}
	return that.ch, atomic.LoadUint64(&that.gen)
}
