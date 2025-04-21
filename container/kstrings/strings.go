package kstrings

import (
	"strconv"
	"strings"
)

// const (
// 	tokOpenString       = `"`
// 	tokCloseString      = `"`
// 	tokMapTypeSeparator = ":"

// 	tokOpenArray  = "["
// 	tokCloseArray = "]"
// 	tokOpenMap    = "{"
// 	tokCloseMap   = "}"
// )

// 从字符串str的index位置开始查找第一个出现的subStr位置
func IndexOf(str, substr string, index int) int {
	if index < 0 || index >= len(str) {
		return -1
	}

	// 从指定索引位置开始查找子串
	pos := strings.Index(str[index:], substr)
	if pos == -1 {
		return -1
	}

	// 考虑偏移量，返回在原始字符串中的位置
	return index + pos
}

func TrimSpace(str string) string {
	return strings.Trim(str, "\x00 \b\t\n\r")
}

func ToFloat32(str string) (float32, error) {
	ret, err := strconv.ParseFloat(str, 32)
	return float32(ret), err
}

func ToFloat64(str string) (float64, error) {
	return strconv.ParseFloat(str, 64)
}

func ToInt(str string) (int, error) {
	return strconv.Atoi(str)
}

func ToInt8(str string) (int8, error) {
	ret, err := strconv.ParseInt(str, 10, 8)
	return int8(ret), err
}

func ToInt16(str string) (int16, error) {
	ret, err := strconv.ParseInt(str, 10, 16)
	return int16(ret), err
}

func ToInt32(str string) (int32, error) {
	ret, err := strconv.ParseInt(str, 10, 32)
	return int32(ret), err
}

func ToInt64(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 64)
}

func ToUint(str string) (uint, error) {
	ret, err := strconv.ParseUint(str, 10, 32)
	return uint(ret), err
}

func ToUint8(str string) (uint8, error) {
	ret, err := strconv.ParseUint(str, 10, 8)
	return uint8(ret), err
}

func ToUint16(str string) (uint16, error) {
	ret, err := strconv.ParseUint(str, 10, 16)
	return uint16(ret), err
}

func ToUint32(str string) (uint32, error) {
	ret, err := strconv.ParseUint(str, 10, 32)
	return uint32(ret), err
}

func ToUint64(str string) (uint64, error) {
	return strconv.ParseUint(str, 10, 64)
}

func ToBool(str string) (bool, error) {
	return strconv.ParseBool(str)
}
