package logger

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/khan-lau/kutils/container/klists"
	"github.com/khan-lau/kutils/container/kobjs"
	"github.com/khan-lau/kutils/container/kstrings"
)

var (
	glog *Logger
)

func init() {
	glog = LoggerInstanceOnlyConsole(int8(DebugLevel))
}

func TestStringParams(t *testing.T) {
	str := "hello word ${param1}! i'm $(param3),  test off ${param2}...${param2}.."
	params := kstrings.Parse(str)
	glog.D(params.Set("param1", "var001").Set("param2", "var002").SetFunc("param3", "fun3").Build())
	glog.D("{}", params.Get())
	glog.D("{}", params.GetVarName())
	glog.D("{}", params.GetFuncName())
}

func Test_Param(t *testing.T) {
	l := klists.New[string]()

	str := "hello word ${param1}! i'm $(param3),  fuck off ${param2}..."

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

	glog.D("{} {} {}", fmt.Errorf("error message"), "aaa", fmt.Errorf("error message2"))
	glog.D("{}", l)
	glog.D("{}", kobjs.ObjectDump(l))

}

func TestLog(t *testing.T) {
	logger := LoggerInstanceOnlyConsole(-1)

	logger.D("fuck off")
	logger.D("{} fuck off", "maybe")

	logger.D("")

	logger.D("int8 {} fuck off", int8(8))
	logger.D("uint8 {} fuck off", uint8(8))
	logger.D("int16 {} fuck off", int16(16))
	logger.D("uint16 {} fuck off", uint16(16))
	logger.D("int {} fuck off", int(10))
	logger.D("uint {} fuck off", uint(10))
	logger.D("int32 {} fuck off", int32(32))
	logger.D("uint32 {} fuck off", uint32(32))
	logger.D("int32 {} fuck off", int32(64))
	logger.D("uint32 {} fuck off", uint32(64))
	logger.D("float32 {} fuck off", float32(4.45))
	logger.D("float64 {} fuck off", float64(2.1))

	logger.D("")

	logger.D("int8 {} fuck off", []int8{0, 1, 2, 3, 4})
	logger.D("uint8 {} fuck off", []uint8{0, 1, 2, 3, 4})
	logger.D("int16 {} fuck off", []int16{0, 1, 2, 3, 4})
	logger.D("uint16 {} fuck off", []uint16{0, 1, 2, 3, 4})
	logger.D("int {} fuck off", []int{0, 1, 2, 3, 4})
	logger.D("uint {} fuck off", []uint{0, 1, 2, 3, 4})
	logger.D("int32 {} fuck off", []int32{0, 1, 2, 3, 4, 0x44, 0x38})
	logger.D("uint32 {} fuck off", []uint32{0, 1, 2, 3, 4})
	logger.D("int64 {} fuck off", []int64{0, 1, 2, 3, 4})
	logger.D("uint64 {} fuck off", []uint64{0, 1, 2, 3, 4})

	logger.D("float32 {} fuck off", []float32{0, 1, 2, 3, 4})
	logger.D("float64 {} fuck off", []float64{0, 1, 2, 3, 4})

	logger.D("")

	logger.D("string {} fuck off", []string{"0", "1", "2", "3", "4"})

	logger.D("")

	cmp := complex(4, 4)
	cmp64 := complex64(cmp)
	logger.D("complex128 {} complex64 {} fuck off", cmp, cmp64)

	logger.D("complex64 {} fuck off", []complex64{complex(4, 0), complex(4, 1), complex(4, 2), complex(4, 3), complex(4, 4)})
	logger.D("complex128 {} fuck off", []complex128{complex(4, 0), complex(4, 1), complex(4, 2), complex(4, 3), complex(4, 4)})

	logger.D("")

	type AA struct {
		A int
		B string
		C complex128
	}

	aa := AA{A: 12, B: "string", C: complex(4, -1)}
	logger.D("obj {} fuck off", aa)
	logger.D("*obj {} fuck off", &aa)
	logger.D("obj {} fuck off", []AA{aa, aa})
	logger.D("obj {} fuck off", []*AA{&aa, &aa})

	logger.D("")
}
