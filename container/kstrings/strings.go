package kstrings

import (
	"strings"
)

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
	return pos + index
}

func TrimSpace(str string) string {
	return strings.Trim(str, "\x00 \b\t\n\r")
}
