package condexpr

type Provider[T any] func() T

// CondExprLazy: 仿三目运算符(条件表达式), expr ? a() : b(), 性能较低, 适用于表达式较为复杂时使用
//   - @param expr 条件表达式
//   - @param a 当expr为true时返回a()的返回值
//   - @param b 当expr为false时返回b()的返回值
//   - @return 当expr为true时返回a()的返回值, 否则返回b()的返回值
//
// @example
//
//	res := CondExprLazy(tempValue == nil, Val(0.0), Ptr(tempValue, 0.0))
//	res := CondExprLazy(tempValue == nil, func() float64 { return 64}, func() float64 { return ExpensiveCall() })
func CondExprLazy[T any](expr bool, a, b Provider[T]) T {
	if expr {
		return a()
	}
	return b()
}

// 定义快捷包装
func Val[T any](v T) Provider[T] { return func() T { return v } }
func Ptr[T any](p *T, def T) Provider[T] {
	return func() T {
		if p == nil {
			return def
		}
		return *p
	}
}

//////////////////////////////////////////////////////

//////////////////////////////////////////////////////

// 仿三目运算符(条件表达式), expr ? a : b, 性能较高, 适用于表达式较为简单时使用
//   - @param expr 条件表达式
//   - @param a 当expr为true时返回a, 必须为常量或常量表达式, 否则非常危险
//   - @param b 当expr为false时返回b, 必须为常量或常量表达式, 否则非常危险
//   - @return 当expr为true时返回a, 否则返回b
//
// @example
//
//	condexpr.CondExpr[string](isServer, "Server", "Client")
//
// expr为true时返回a, 否则返回b
func CondExpr[T comparable](expr bool, a, b T) T {
	if expr {
		return a
	}
	return b
}

//////////////////////////////////////////////////////

//////////////////////////////////////////////////////

// SafeGet: 如果指针不为空且满足条件，返回其值，否则返回备选值
//   - @param ptr 指针
//   - @param predicate 当ptr不为nil时，*ptr值的判断callback
//   - @param fallback 当ptr为空或predicate不满足条件时返回的备选值
//   - @return 当ptr非空且predicate满足条件时返回其指向的值, 否则返回fallback
//
// @example
//
//	condexpr.SafeGet(ptr, func(v string) bool { return v != "" }, "default")
//	condexpr.SafeGet(tempValue, func(v float64) bool { return v == 34 }, condexpr.Unwrap(tempValue1, 0))
func SafeGet[T any](ptr *T, predicate func(T) bool, fallback T) T {
	if ptr != nil && predicate(*ptr) {
		return *ptr
	}
	return fallback
}

// Unwrap: 基础的解引用助手, 如果ptr为nil, 返回defaultVal, 否则返回*ptr指向内容
//   - @param ptr 指针
//   - @param defaultVal 当ptr为nil时返回defaultVal
//   - @return 当ptr为nil时返回defaultVal, 否则返回*ptr指向内容
//
// @example
//
//	condexpr.Unwrap(ptr, "default")
func Unwrap[T any](ptr *T, defaultVal T) T {
	if ptr == nil {
		return defaultVal
	}
	return *ptr
}
