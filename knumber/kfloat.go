package knumber

import "math"

// 野鸡四舍五入法
func Round(val float64) int {
	return int(math.Floor(val + 0.5))
}

// 浮点数相等比较, 8位精度, 最大可以精确到15~17位
func EqualF64(a, b float64) bool {
	return math.Abs(a-b) < 1e-8
}

// 浮点数相等比较, 8位精度, 最大可以精确到6~7位
func EqualF32(a, b float32) bool {
	return math.Abs(float64(a-b)) < 1e-4
}

// any 数字类型 转 float32
func ToFloat32(val interface{}) (float32, bool) {
	if val == nil {
		return 0, false
	}
	switch v := val.(type) {
	case float32:
		return v, true
	case float64:
		return float32(v), true
	case int:
		return float32(v), true
	case int8:
		return float32(v), true
	case int16:
		return float32(v), true
	case int32:
		return float32(v), true
	case int64:
		return float32(v), true
	case uint:
		return float32(v), true
	case uint8:
		return float32(v), true
	case uint16:
		return float32(v), true
	case uint32:
		return float32(v), true
	case uint64:
		return float32(v), true
	}
	return 0, false
}

// any 数字类型 转 float64
func ToFloat64(val interface{}) (float64, bool) {
	if val == nil {
		return 0, false
	}
	switch v := val.(type) {
	case float32:
		return float64(v), true
	case float64:
		return v, true
	case int:
		return float64(v), true
	case int8:
		return float64(v), true
	case int16:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint:
		return float64(v), true
	case uint8:
		return float64(v), true
	case uint16:
		return float64(v), true
	case uint32:
		return float64(v), true
	case uint64:
		return float64(v), true
	}
	return 0, false
}
