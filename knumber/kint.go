package knumber

// ToInt 将任意类型转换为int，转换失败返回0和false。
func ToInt(val interface{}) (int, bool) {
	switch v := val.(type) {
	case int:
		return v, true
	case int8:
		return int(v), true
	case int16:
		return int(v), true
	case int32:
		return int(v), true
	case int64:
		return int(v), true
	case uint:
		return int(v), true
	case uint8:
		return int(v), true
	case uint16:
		return int(v), true
	case uint32:
		return int(v), true
	case uint64:
		return int(v), true
	case float32:
		return int(v), true
	case float64:
		return int(v), true
	default:
		return 0, false
	}
}

// ToInt8 将任意类型转换为int8，转换失败返回0和false。
func ToInt8(val interface{}) (int8, bool) {
	switch v := val.(type) {
	case int:
		return int8(v), true
	case int8:
		return v, true
	case int16:
		return int8(v), true
	case int32:
		return int8(v), true
	case int64:
		return int8(v), true
	case uint:
		return int8(v), true
	case uint8:
		return int8(v), true
	case uint16:
		return int8(v), true
	case uint32:
		return int8(v), true
	case uint64:
		return int8(v), true
	case float32:
		return int8(v), true
	case float64:
		return int8(v), true
	default:
		return 0, false
	}
}

// ToInt16 将任意类型转换为int16，转换失败返回0和false。
func ToInt16(val interface{}) (int16, bool) {
	switch v := val.(type) {
	case int:
		return int16(v), true
	case int8:
		return int16(v), true
	case int16:
		return v, true
	case int32:
		return int16(v), true
	case int64:
		return int16(v), true
	case uint:
		return int16(v), true
	case uint8:
		return int16(v), true
	case uint16:
		return int16(v), true
	case uint32:
		return int16(v), true
	case uint64:
		return int16(v), true
	case float32:
		return int16(v), true
	case float64:
		return int16(v), true
	default:
		return 0, false
	}
}

// ToInt32 将任意类型转换为int32，转换失败返回0和false。
func ToInt32(val interface{}) (int32, bool) {
	switch v := val.(type) {
	case int:
		return int32(v), true
	case int8:
		return int32(v), true
	case int16:
		return int32(v), true
	case int32:
		return v, true
	case int64:
		return int32(v), true
	case uint:
		return int32(v), true
	case uint8:
		return int32(v), true
	case uint16:
		return int32(v), true
	case uint32:
		return int32(v), true
	case uint64:
		return int32(v), true
	case float32:
		return int32(v), true
	case float64:
		return int32(v), true
	default:
		return 0, false
	}
}

// ToInt64 将任意类型转换为int64，转换失败返回0和false。
func ToInt64(val interface{}) (int64, bool) {
	switch v := val.(type) {
	case int:
		return int64(v), true
	case int8:
		return int64(v), true
	case int16:
		return int64(v), true
	case int32:
		return int64(v), true
	case int64:
		return v, true
	case uint:
		return int64(v), true
	case uint8:
		return int64(v), true
	case uint16:
		return int64(v), true
	case uint32:
		return int64(v), true
	case uint64:
		return int64(v), true
	case float32:
		return int64(v), true
	case float64:
		return int64(v), true
	default:
		return 0, false
	}
}

// ToUint 将任意类型转换为uint，转换失败返回0和false。
func ToUint(val interface{}) (uint, bool) {
	switch v := val.(type) {
	case int:
		return uint(v), true
	case int8:
		return uint(v), true
	case int16:
		return uint(v), true
	case int32:
		return uint(v), true
	case int64:
		return uint(v), true
	case uint:
		return v, true
	case uint8:
		return uint(v), true
	case uint16:
		return uint(v), true
	case uint32:
		return uint(v), true
	case uint64:
		return uint(v), true
	case float32:
		return uint(v), true
	case float64:
		return uint(v), true
	default:
		return 0, false
	}
}

// ToUint8 将任意类型转换为uint8，转换失败返回0和false。
func ToUint8(val interface{}) (uint8, bool) {
	switch v := val.(type) {
	case int:
		return uint8(v), true
	case int8:
		return uint8(v), true
	case int16:
		return uint8(v), true
	case int32:
		return uint8(v), true
	case int64:
		return uint8(v), true
	case uint:
		return uint8(v), true
	case uint8:
		return v, true
	case uint16:
		return uint8(v), true
	case uint32:
		return uint8(v), true
	case uint64:
		return uint8(v), true
	case float32:
		return uint8(v), true
	case float64:
		return uint8(v), true
	default:
		return 0, false
	}
}

// ToUint16 将任意类型转换为uint16，转换失败返回0和false。
func ToUint16(val interface{}) (uint16, bool) {
	switch v := val.(type) {
	case int:
		return uint16(v), true
	case int8:
		return uint16(v), true
	case int16:
		return uint16(v), true
	case int32:
		return uint16(v), true
	case int64:
		return uint16(v), true
	case uint:
		return uint16(v), true
	case uint8:
		return uint16(v), true
	case uint16:
		return v, true
	case uint32:
		return uint16(v), true
	case uint64:
		return uint16(v), true
	case float32:
		return uint16(v), true
	case float64:
		return uint16(v), true
	default:
		return 0, false
	}
}

// ToUint32 将任意类型转换为uint32，转换失败返回0和false。
func ToUint32(val interface{}) (uint32, bool) {
	switch v := val.(type) {
	case int:
		return uint32(v), true
	case int8:
		return uint32(v), true
	case int16:
		return uint32(v), true
	case int32:
		return uint32(v), true
	case int64:
		return uint32(v), true
	case uint:
		return uint32(v), true
	case uint8:
		return uint32(v), true
	case uint16:
		return uint32(v), true
	case uint32:
		return v, true
	case uint64:
		return uint32(v), true
	case float32:
		return uint32(v), true
	case float64:
		return uint32(v), true
	default:
		return 0, false
	}
}

// ToUint64 将任意类型转换为uint64，转换失败返回0和false。
func ToUint64(val interface{}) (uint64, bool) {
	switch v := val.(type) {
	case int:
		return uint64(v), true
	case int8:
		return uint64(v), true
	case int16:
		return uint64(v), true
	case int32:
		return uint64(v), true
	case int64:
		return uint64(v), true
	case uint:
		return uint64(v), true
	case uint8:
		return uint64(v), true
	case uint16:
		return uint64(v), true
	case uint32:
		return uint64(v), true
	case uint64:
		return v, true
	case float32:
		return uint64(v), true
	case float64:
		return uint64(v), true
	default:
		return 0, false
	}
}
