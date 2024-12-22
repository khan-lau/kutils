package katomic

import (
	"encoding/json"
	"strconv"
	"sync/atomic"
	"unsafe"
)

type NumberConstraints interface {
	int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64
}

type IntNumber[T NumberConstraints] struct {
	_ nocmp // disallow non-atomic comparison
	v T
}

// NewStatus 创建一个新的 Status 实例，并设置其值为 val
func NewStatus[T NumberConstraints](val T) *IntNumber[T] {
	return &IntNumber[T]{v: val}
}

// 加载 Status 实例中存储的值。
//
//	返回值:
//		@return 返回一个类型为 T 的值。如果 T 是 int8, int16, int32, int64, uint8, uint16, uint32 或 uint64 之一，它将使用 atomic 包中的 Load 函数来安全地加载值。如果 T 不是这些类型之一，则会导致 panic。
func (that *IntNumber[T]) Load() T {
	var result T
	switch any(result).(type) {
	case int8, int16, int32:
		return T(atomic.LoadInt32((*int32)(unsafe.Pointer(&that.v))))
	case int64:
		return T(atomic.LoadInt64((*int64)(unsafe.Pointer(&that.v))))
	case uint8, uint16, uint32:
		return T(atomic.LoadUint32((*uint32)(unsafe.Pointer(&that.v))))
	case uint64:
		return T(atomic.LoadUint64((*uint64)(unsafe.Pointer(&that.v))))
	default:
		panic("unsupported type")
	}
}

// 向 Status 结构体的值添加指定的增量，并返回新的值。
//
//	参数:
//		@param delta 是要添加的增量，类型为泛型 T。
//	返回值:
//		@return 是添加增量后的新值，类型为泛型 T。
func (that *IntNumber[T]) Add(delta T) T {
	var result T
	switch any(result).(type) {
	case int8, int16, int32:
		return T(atomic.AddInt32((*int32)(unsafe.Pointer(&that.v)), int32(delta)))
	case int64:
		return T(atomic.AddInt64((*int64)(unsafe.Pointer(&that.v)), int64(delta)))
	case uint8, uint16, uint32:
		return T(atomic.AddUint32((*uint32)(unsafe.Pointer(&that.v)), uint32(delta)))
	case uint64:
		return T(atomic.AddUint64((*uint64)(unsafe.Pointer(&that.v)), uint64(delta)))
	default:
		panic("unsupported type")
	}
}

// Sub 从当前状态值中减去 delta 并返回新的状态值。
//
//	参数：
//		@param delta: 要减去的值。
//	返回值：
//		@return 返回减去 delta 后的新状态值。
func (that *IntNumber[T]) Sub(delta T) T {
	var result T
	switch any(result).(type) {
	case int8, int16, int32:
		return T(atomic.AddInt32((*int32)(unsafe.Pointer(&that.v)), -int32(delta)))
	case int64:
		return T(atomic.AddInt64((*int64)(unsafe.Pointer(&that.v)), -int64(delta)))
	case uint8, uint16, uint32:
		return T(atomic.AddUint32((*uint32)(unsafe.Pointer(&that.v)), ^uint32(delta-1)))
	case uint64:
		return T(atomic.AddUint64((*uint64)(unsafe.Pointer(&that.v)), ^uint64(delta-1)))
	default:
		panic("unsupported type")
	}
}

// Inc 方法对 Status 结构体中的值进行自增操作
//
//	参数:
//		that: *Status[T] 类型的指针，表示要进行递增操作的状态对象
//	返回值:
//		T: 递增后的值
func (that *IntNumber[T]) Inc() T {
	return that.Add(1)
}

// Dec 方法用于对 Status 类型的值进行递减操作
//
//	参数:
//		that: *Status[T] 类型的指针，表示要进行递减操作的状态对象
//	返回值:
//		T: 递减后的值
func (that *IntNumber[T]) Dec() T {
	return that.Sub(1)
}

// CAS 函数尝试以原子方式将 IntNumber 的值从旧值更新为新值。
// 如果当前值等于旧值，则将其设置为新值，并返回 true，表示成功交换。
// 如果当前值不等于旧值，则不进行修改，并返回 false，表示未进行交换。
//
//	参数：
//		@param old：期望的旧值
//		@param new：要设置的新值
//	返回值：
//		@return 如果成功交换了值，则返回 true；否则返回 false。
func (that *IntNumber[T]) CAS(old, new T) (swapped bool) {
	return that.CompareAndSwap(old, new)
}

// CompareAndSwap 方法尝试将 IntNumber 的值从 old 更新为 new。
// 如果当前值为 old，则更新为 new 并返回 true；否则，不执行更新并返回 false。
//
//	参数：
//		@param old：期望的旧值
//		@param new：要设置的新值
//	返回值：
//		@return 如果成功交换了值，则返回 true；否则返回 false。
func (that *IntNumber[T]) CompareAndSwap(old, new T) bool {
	var swapped bool
	switch any(that.v).(type) {
	case int8, int16, int32:
		swapped = atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&that.v)), int32(old), int32(new))
	case int64:
		swapped = atomic.CompareAndSwapInt64((*int64)(unsafe.Pointer(&that.v)), int64(old), int64(new))
	case uint8, uint16, uint32:
		swapped = atomic.CompareAndSwapUint32((*uint32)(unsafe.Pointer(&that.v)), uint32(old), uint32(new))
	case uint64:
		swapped = atomic.CompareAndSwapUint64((*uint64)(unsafe.Pointer(&that.v)), uint64(old), uint64(new))
	default:
		panic("unsupported type")
	}
	return swapped
}

// Store 将给定的值存储到 IntNumber 实例中。
//
//	参数:
//		@param val: 要存储的值。
//
//	说明:
//		该方法使用 atomic 包提供的原子操作来保证对 IntNumber 实例中值的修改是线程安全的。
//		它首先使用 type switch 判断 that.v 的具体类型，然后根据不同的类型使用相应的 atomic.Store 函数进行存储。
//		如果 that.v 的类型不是 int8, int16, int32, int64, uint8, uint16, uint32, uint64 之一，
//		则会引发 panic，并抛出 "unsupported type" 错误信息。
func (that *IntNumber[T]) Store(val T) {
	switch any(that.v).(type) {
	case int8, int16, int32:
		atomic.StoreInt32((*int32)(unsafe.Pointer(&that.v)), int32(val))
	case int64:
		atomic.StoreInt64((*int64)(unsafe.Pointer(&that.v)), int64(val))
	case uint8, uint16, uint32:
		atomic.StoreUint32((*uint32)(unsafe.Pointer(&that.v)), uint32(val))
	case uint64:
		atomic.StoreUint64((*uint64)(unsafe.Pointer(&that.v)), uint64(val))
	default:
		panic("unsupported type")
	}
}

// Swap 方法用于原子性地交换 IntNumber 实例中的值。
//
//	参数:
//		@param val 是要交换进去的新值。
//	返回值:
//		@retuan 返回原来的值 old。如果 IntNumber 实例中存储的值类型不支持原子操作，则触发 panic。
func (that *IntNumber[T]) Swap(val T) (old T) {
	switch any(that.v).(type) {
	case int8, int16, int32:
		return T(atomic.SwapInt32((*int32)(unsafe.Pointer(&that.v)), int32(val)))
	case int64:
		return T(atomic.SwapInt64((*int64)(unsafe.Pointer(&that.v)), int64(val)))
	case uint8, uint16, uint32:
		return T(atomic.SwapUint32((*uint32)(unsafe.Pointer(&that.v)), uint32(val)))
	case uint64:
		return T(atomic.SwapUint64((*uint64)(unsafe.Pointer(&that.v)), uint64(val)))
	default:
		panic("unsupported type")
	}
}

func (that *IntNumber[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(that.Load())
}

func (that *IntNumber[T]) UnmarshalJSON(b []byte) error {
	var v T
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	that.Store(v)
	return nil
}

func (that *IntNumber[T]) String() string {
	v := that.Load()
	return strconv.FormatUint(uint64(v), 10)
}
