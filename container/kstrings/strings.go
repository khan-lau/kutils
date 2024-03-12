package kstrings

import (
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

// func FormatJson5(json5str string) (string, error) {
// 	trimStr := TrimSpace(json5str)
// 	ident := "  "

// 	first := trimStr[:1]
// 	if first != "[" && first != "{" {
// 		return "", fmt.Errorf("first charset must be '[' or '{'")
// 	}

// 	var sb strings.Builder

// 	tok := klists.New[string]()
// 	tok.PushBack(first)
// 	sb.WriteString(FormatString("{}\n{}", first, strings.Repeat(ident, 1)))
// 	trimStr = trimStr[1:]

// 	pos := strings.IndexAny(trimStr, "[{")
// 	sub := ""

// 	length := len(trimStr)
// 	remain := length

// 	for remain < 1 {
// 		// 取字段名
// 		pos := strings.IndexAny(trimStr, ":")
// 		if pos > -1 {
// 			key := trimStr[:pos]
// 			key = TrimSpace(key)
// 			sb.WriteString(FormatString("{} : {}", strings.Repeat(ident, tok.Len()), key))
// 			remain = remain - (pos + 1)
// 			trimStr = trimStr[(pos + 1):] // `:`后的内容

// 			// 判断`,`前有没有`"` 或`'` 字符串起止符, 要排除 \' \" 转义符
// 		}
// 	}

// 	if pos != -1 {
// 		sub = trimStr[:pos+1]
// 		trimStr = trimStr[(pos + 1):]
// 		tok.PushBack(sub)
// 	}

// 	Println("pos:{}, sub:{}", pos, sub)
// 	return sb.String(), nil
// }
