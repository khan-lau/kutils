package condexpr

// 仿三目运算符(条件表达式), expr ? a : b.
//
// @example
//     condexpr.CondExpr[string](isServer, "Server", "Client")
//
// expr为true时返回a, 否则返回b
func CondExpr[T comparable](expr bool, a, b T) T {
	if expr {
		return a
	}
	return b
}
