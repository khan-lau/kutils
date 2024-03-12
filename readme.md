# KUtils v0.1.0

## kuuid 
实现uuidV1 

## container 
1. klists 实现泛型的list
2. kstrings 实现 StringFormatter
3. kobjs 实现obj to json5 string
4. kslices slice的补充工具库
5. kmaps map的补充工具库

```go
// StringFormatter 范例
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
```


```go
// to json5 string example

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

	str := ""
	// str = ObjectDump(&myInstance)
	// log.Printf("%s", str)

	str = ObjectToJson5(myInstance) //获取对象树信息
	kstrings.Debug("{}", str)

	kstrings.Debug("")
}
```
输出结果为
```bash
[2024-03-12 13:15:44.594] object_test.go:81     {
  Err:null,
  PublicField:42,
  privateField:"private field 0",
  Func:"func(string)(string, error)",
  StrList:[
    "aa",
    "bb",
    "cc",
    "dd"
  ],
  List:[
    "aa",
    "bb",
    "cc",
    "dd"
  ],
  Map:{
    key01:"value01",
    key02:"value02"
  },
  strs:[
    "str1",
    "str2",
    "str3"
  ],
  func:{
    PublicConstMethod0:"func(string)",
    PublicMethod1:"func(string)(string, error)"
  }
}
```

## data
1. gzip deflate br 的压缩与解压
2. Generator 自增原子数

## datetime
将时间段按自然周期分组
例如, `自然天` `自然周` `自然月` `自然年`

## file_format
efile 格式解析

## filesystem
文件系统补充工具库

## kcrypto
chacha20算法 加密解密

## db
数据库相关

### kredis
基于`go-redis/v9`的一些简单封装

## logger
基于zap 与 file-rotatelogs 的日志库简单封装

```go

logger := LoggerInstanceOnlyConsole(-1)
	logger.D("string {} fuck off", []string{"0", "1", "2", "3", "4"})

	cmp := complex(4, 4)
	cmp64 := complex64(cmp)
	logger.D("complex128 {} complex64 {} fuck off", cmp, cmp64)

	logger.D("complex64 {} fuck off", []complex64{complex(4, 0), complex(4, 1), complex(4, 2), complex(4, 3), complex(4, 4)})
	logger.D("complex128 {} fuck off", []complex128{complex(4, 0), complex(4, 1), complex(4, 2), complex(4, 3), complex(4, 4)})


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


```