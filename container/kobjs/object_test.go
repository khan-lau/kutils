package kobjs

import (
	"container/list"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/khan-lau/kutils/container/klists"
	"github.com/khan-lau/kutils/container/kstrings"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

type MyStruct struct {
	Err          error
	PublicField  int
	privateField string
	Func         func(aa string) (string, error)
	StrList      *klists.KList[string]
	List         *list.List
	Map          map[string]string
	strs         []string
}

func (s MyStruct) PublicConstMethod0(str string) {
	fmt.Println("Public const method called")
}

func (s *MyStruct) PublicMethod1(str string) (string, error) {
	fmt.Println("Public method called, private field: ", s.privateField)
	return "", nil
}

func (s MyStruct) privateCibstMethod0() {
	fmt.Println("private method called, private field: ", s.privateField)
}

func Test_ObjectDump(t *testing.T) {
	tok := klists.New[string]()
	tok.PushBack("aa")
	tok.PushBack("bb")
	tok.PushBack("cc")
	tok.PushBack("dd")

	to := list.New()

	to.PushBack("aa")
	to.PushBack("bb")
	to.PushBack("cc")
	to.PushBack("dd")

	myInstance := MyStruct{
		PublicField:  42,
		privateField: "private field 0",
		Func: func(aa string) (string, error) {
			return "", nil
		},
		StrList: tok,
		List:    to,
		Map:     make(map[string]string, 10),
		strs:    []string{"str1", "str2", "str3"},
	}

	myInstance.Map["key01"] = "value01"
	myInstance.Map["key02"] = "value02"

	kstrings.Debug("{}", strings.Repeat("a", 2))

	str := ""
	// str = ObjectDump(&myInstance)
	// log.Printf("%s", str)

	// log.Printf("")

	str = ObjectToJson5(myInstance)
	kstrings.Debug("{}", str)

	kstrings.Debug("")

	// kstrings.FormatJson5(str)
}

func Test_FormatJson5(t *testing.T) {
	tok := klists.New[string]()
	tok.PushBack("aa")
	tok.PushBack("bb")
	tok.PushBack("cc")
	tok.PushBack("dd")

	// log.Printf("%T\n", tok)
	// log.Printf("%v\n", tok)
	// log.Printf("%#v\n", tok)

	to := list.New()

	to.PushBack("aa")
	to.PushBack("bb")
	to.PushBack("cc")
	to.PushBack("dd")

	map01 := make(map[string]string, 10)
	map01["key01"] = "value01"
	map01["key02"] = "value02"

	str := ObjectDump(map01)
	kstrings.Debug("{}", str)
}

func Test_StringFormatter(t *testing.T) {
	log.Printf("%s \n", kstrings.FormatArgs("Set {} is not equal to {}.", 1, 2).Message())

	log.Println(kstrings.ArrayFormat("{\"user\":\"{}\", \"token\":\"{}\"}", []any{"test", "test001"}).Message())

	log.Println(kstrings.SliceFormat("{\"code\":{},  \"user\":\"{}\", \"token\":\"{}\"}", 0, "test", "test001").Message())

	log.Println(kstrings.SliceFormat("int:{}, double:{}, bool:{}, string:{}, array[int]:{}, array[bool]:{}",
		0, 0.1, false, "test001", []any{1, 2, 3}, []any{false, true, true}).Message())

	log.Println(kstrings.FormatString("int:{}, double:{}, bool:{}, string:{}, array[int]:{}, array[bool]:{}",
		0, 0.1, false, "test001", []any{1, 2, 3}, []any{false, true, true}))

	mapTmp := make(map[string]int, 5)
	mapTmp["1"] = 1
	mapTmp["2"] = 2
	log.Println(kstrings.FormatString("map[string]int : {}", mapTmp))

	kstrings.Println("int:{}, double:{}, bool:{}, string:{}, array[int]:{}, array[bool]:{}",
		0, 0.1, false, "test001", []any{1, 2, 3}, []any{false, true, true})
}
