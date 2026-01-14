package kstrings

import (
	"strconv"
	"strings"
	"time"

	"github.com/khan-lau/kutils/datetime"
)

func ToInt8NoThrow(str string) (int8, bool) {
	str = strings.TrimSpace(str)
	if val, err := strconv.ParseInt(str, 10, 8); err == nil {
		return int8(val), true
	}
	return 0, false
}

func ToInt16NoThrow(str string) (int16, bool) {
	str = strings.TrimSpace(str)
	if val, err := strconv.ParseInt(str, 10, 16); err == nil {
		return int16(val), true
	}
	return 0, false
}

func ToInt32NoThrow(str string) (int32, bool) {
	str = strings.TrimSpace(str)
	if val, err := strconv.ParseInt(str, 10, 32); err == nil {
		return int32(val), true
	}
	return 0, false
}

func ToIntNoThrow(str string) (int, bool) {
	str = strings.TrimSpace(str)
	if val, err := strconv.Atoi(str); err == nil {
		return val, true
	}
	return 0, false
}

func ToInt64NoThrow(str string) (int64, bool) {
	str = strings.TrimSpace(str)
	if val, err := strconv.ParseInt(str, 10, 64); err == nil {
		return val, true
	}
	return 0, false
}

func ToUint8NoThrow(str string) (uint8, bool) {
	str = strings.TrimSpace(str)
	if val, err := strconv.ParseUint(str, 10, 8); err == nil {
		return uint8(val), true
	}
	return 0, false
}

func ToUint16NoThrow(str string) (uint16, bool) {
	str = strings.TrimSpace(str)
	if val, err := strconv.ParseUint(str, 10, 16); err == nil {
		return uint16(val), true
	}
	return 0, false
}

func ToUint32NoThrow(str string) (uint32, bool) {
	str = strings.TrimSpace(str)
	if val, err := strconv.ParseUint(str, 10, 32); err == nil {
		return uint32(val), true
	}
	return 0, false
}

func ToUintNoThrow(str string) (uint, bool) {
	str = strings.TrimSpace(str)
	if val, err := strconv.ParseUint(str, 10, 32); err == nil {
		return uint(val), true
	}
	return 0, false
}

func ToUint64NoThrow(str string) (uint64, bool) {
	str = strings.TrimSpace(str)
	if val, err := strconv.ParseUint(str, 10, 64); err == nil {
		return val, true
	}
	return 0, false
}

func ToFloat64NoThrow(str string) (float64, bool) {
	str = strings.TrimSpace(str)
	if val, err := strconv.ParseFloat(str, 64); err == nil {
		return val, true
	}
	return 0, false
}

func ToFloat32NoThrow(str string) (float32, bool) {
	str = strings.TrimSpace(str)
	if val, err := strconv.ParseFloat(str, 32); err == nil {
		return float32(val), true
	}
	return 0, false
}

func ToBoolNoThrow(str string) (bool, bool) {
	str = strings.TrimSpace(str)
	if val, err := strconv.ParseBool(str); err == nil {
		return val, true
	}
	return false, false
}

func ToDateTime(str string) (*time.Time, bool) {
	str = strings.TrimSpace(str)
	if val, err := time.Parse(datetime.DATETIME_FORMATTER, str); err == nil {
		return &val, true
	}
	return nil, false
}

func ToDateTimeByFormat(str string, format string) (*time.Time, bool) {
	str = strings.TrimSpace(str)
	if val, err := time.Parse(format, str); err == nil {
		return &val, true
	}

	return nil, false
}

func ToTimestamp(str string) (uint64, bool) {
	str = strings.TrimSpace(str)
	if val, err := time.Parse(datetime.DATETIME_FORMATTER, str); err == nil {
		return uint64(val.Unix()), true
	}
	return 0, false
}

func ToTimestampForMilli(str string) (uint64, bool) {
	str = strings.TrimSpace(str)
	if val, err := time.Parse(datetime.DATETIME_FORMATTER, str); err == nil {
		return uint64(val.UnixMilli()), true
	}
	return 0, false
}

// ToDateTimeByFormatWithLocation 根据指定的格式和时区偏移量将字符串转换为时间对象
//
// 该函数会先去除字符串两端的空白字符，然后使用指定的时区偏移量来解析时间字符串。
// 如果解析成功，返回时间对象和true；如果解析失败，返回nil和false。
//
//	@param str 要解析的时间字符串
//	@param time.ParseInLocation函数使用的格式字符串
//	@param offset 时区偏移量，单位为秒（东八区为+28800）
//	@return 解析成功返回*time.Time对象，失败返回nil
//	@return 解析成功返回true，失败返回false
func ToDateTimeByFormatWithLocation(str string, format string, offset int) (*time.Time, bool) {
	str = strings.TrimSpace(str)
	loc := time.FixedZone("", offset)
	if val, err := time.ParseInLocation(format, str, loc); err == nil {
		return &val, true
	}

	return nil, false
}
