package kstrings

import (
	"regexp"
	"strings"

	"github.com/khan-lau/kutils/container/klists"
	"github.com/khan-lau/kutils/container/kmaps"
)

type TYPE uint

const (
	NONE_TYPE TYPE = 0
	VAR_TYPE  TYPE = 1
	FUNC_TYPE TYPE = 2
)

// 字符串参数解析与处理
//   - 例如: "hello word ${param1}! i'm $(param3),  fuck off ${param2}..."
//   - - ${param1} ${} 代表变量
//   - - $(param3) $() 代表方法
type StringParams struct {
	srcStr string
	params map[string]string
}

// 解析表达式
func Parse(str string) *StringParams {
	l := klists.New[string]()

	// pattern := `\$\{([^}]+)\}` //1个捕获组
	// pattern := `\$\(([^)]+)\)` //1个捕获组
	pattern := `\$\{([^}]+)\}|\$\(([^)]+)\)` //两个捕获组
	regex := regexp.MustCompile(pattern)
	matches := regex.FindAllStringSubmatch(str, -1)
	for _, match := range matches {
		// match 格式为 {"${param1}", "param1", ""}
		// 第一个元素是完整的匹配项，例如 ${param1} 或 $(param3)
		// 第二个元素是第一个捕获组的内容，例如 param1 或 param3。
		// 第三个元素是第二个捕获组的内容，如果没有匹配到，则为空字符串。
		l.PushBack(match[0])
	}
	paramMap := make(map[string]string)
	for e := l.Front(); e != nil; e = e.Next() {
		paramMap[e.Value] = ""
	}

	return &StringParams{srcStr: str, params: paramMap}
}

// 指定一个方法的值, 如方法不存在, 不做任何操作
func (that *StringParams) SetFunc(param string, value string) *StringParams {
	key := "$(" + param + ")"
	if kmaps.HasKey[string, string](that.params, key) {
		that.params[key] = value
	}
	return that
}

// 指定一个变量的值, 如变量不存在, 不做任何操作
func (that *StringParams) Set(param string, value string) *StringParams {
	return that.SetVar(param, value)
}

// 指定一个变量的值, 如变量不存在, 不做任何操作
func (that *StringParams) SetVar(param string, value string) *StringParams {
	key := "${" + param + "}"
	if kmaps.HasKey[string, string](that.params, key) {
		that.params[key] = value
	}
	return that
}

// 将包含变量的表达式编译成字符串
func (that *StringParams) Build() string {
	str := that.srcStr
	for key, val := range that.params {
		str = strings.ReplaceAll(str, key, val)
	}
	return str
}

func (that *StringParams) Get() *klists.KList[string] {
	l := klists.New[string]()
	for key := range that.params {
		l.PushBack(key)
	}
	return l
}

func (that *StringParams) GetVarName() *klists.KList[string] {
	l := klists.New[string]()
	for key := range that.params {
		if strings.HasPrefix(key, "${") {
			str := key[2 : len(key)-1]
			l.PushBack(str)
		}
	}
	return l
}

func (that *StringParams) GetFuncName() *klists.KList[string] {
	l := klists.New[string]()
	for key := range that.params {
		if strings.HasPrefix(key, "$(") {
			str := key[2 : len(key)-1]
			l.PushBack(str)
		}
	}
	return l
}

type KParameter string

func (that KParameter) Type() TYPE {
	if strings.HasPrefix(string(that), "${") {
		return VAR_TYPE
	} else if strings.HasPrefix(string(that), "$(") {
		return FUNC_TYPE
	} else {
		return NONE_TYPE
	}
}

func (that KParameter) TypeString() string {
	if strings.HasPrefix(string(that), "${") {
		return "VAR_TYPE"
	} else if strings.HasPrefix(string(that), "$(") {
		return "FUNC_TYPE"
	} else {
		return "NONE_TYPE"
	}
}
