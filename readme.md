# KUtils v0.2.3

## knumber
interface{} 数字类型 to float32 float64 int int8 int16 int32 int64 uint uint8 uint16 uint32 uint64

```go
anymap := make(map[string]interface{})
anymap["key0"] = 10

if val0, ok := kunumber.ToInt32(anymap["key0"]); ok{ // 返回 int32
	kstrings.Println("int32 val0:{}", val0)
}
```

## kuuid 
实现uuidV1 

## container 
1. klists 实现泛型的list
2. kstrings 实现 StringFormatter 与 StringParams
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
// StringParams 范例
func TestStringParams(t *testing.T) {
	// $(func_name) 暂且认为是函数
	// ${var_name} 暂且认为是变量 
	str := "hello word ${param1}! i'm $(param3),  test off ${param2}...${param2}.."
	params := kstrings.Parse(str)
	glog.D(params.Set("param1", "var001").Set("param2", "var002").SetFunc("param3", "fun3").Build())
	glog.D("{}", params.Get())
	glog.D("{}", params.GetVarName())
	glog.D("{}", params.GetFuncName())
}
```
输出结果
```log
2024-03-21 11:51:24.038 DEBUG   logger/logger_test.go:24        hello word var001! i'm fun3,  test off var002...var002..
2024-03-21 11:51:24.047 DEBUG   logger/logger_test.go:25        ["${param1}","$(param3)","${param2}"]
2024-03-21 11:51:24.047 DEBUG   logger/logger_test.go:26        ["param1","param2"]
2024-03-21 11:51:24.047 DEBUG   logger/logger_test.go:27        ["param3"]
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

```go
// 获取指定时间所在月份的天数
func MonthDays(time Time.Time) uint

// 获取指定时间所在年份的天数
func YearDays(time Time.Time) uint


// 当前时间所在 分 的起始秒时间戳
func MinuteFirst(time Time.Time) uint64

// 取1分钟的开始时间与结束时间
func FirstAndLastMinute(time Time.Time) (uint64, uint64)


// 当前时间所在 小时 的起始秒时间戳
func HourFirst(time Time.Time) uint64

// 取1小时的开始时间与结束时间
func FirstAndLastHour(time Time.Time) (uint64, uint64)


// 当前时间所在 天 的起始秒时间戳
func DayFirst(time Time.Time) uint64 

// 取1天的开始时间与结束时间
func FirstAndLastDay(time Time.Time) (uint64, uint64)


// 当前时间所在 周 的起始秒时间戳
func WeekFirst(time Time.Time) uint64

// 取1周的开始时间与结束时间
func FirstAndLastWeek(time Time.Time) (uint64, uint64)


// 当前时间所在 月 的起始秒时间戳
func MonthFirst(time Time.Time) uint64

// 取1个月的开始时间与结束时间
func FirstAndLastMonth(time Time.Time) (uint64, uint64)


// 当前时间所在 年 的起始秒时间戳
func YearFirst(time Time.Time) uint64

// 取1年的开始时间与结束时间
func FirstAndLastYear(time Time.Time) (uint64, uint64)


// 将起止时间按指定周期分割, 返回每个周期的起止时间
//   - @param time.Time start    开始时间
//   - @param time.Time end      结束时间
//   - @param Duration  duration 分割周期
//   - @return []*TimeSlice  每个分段的起止时间
func SplitDuration(start, end Time.Time, duration Duration) []*TimeSlice
```

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

## klogger
基于zap 与 file-rotatelogs 的日志库简单封装

```go

logger := LoggerInstanceOnlyConsole(int8(DebugLevel))
	logger.D("{}", fmt.Errorf("test error"))
	logger.D("string {} fuck off", []string{"0", "1", "2", "3", "4"})

	cmp := complex(4, 4)
	cmp64 := complex64(cmp)
	logger.D("complex128 {} complex64 {} test off", cmp, cmp64)

	logger.D("complex64 {} test off", []complex64{complex(4, 0), complex(4, 1), complex(4, 2), complex(4, 3), complex(4, 4)})
	logger.D("complex128 {} test off", []complex128{complex(4, 0), complex(4, 1), complex(4, 2), complex(4, 3), complex(4, 4)})


	type AA struct {
		A int
		B string
		C complex128
	}
	aa := AA{A: 12, B: "string", C: complex(4, -1)}
	logger.D("obj {} test off", aa)
	logger.D("*obj {} test off", &aa)
	logger.D("obj {} test off", []AA{aa, aa})
	logger.D("obj {} test off", []*AA{&aa, &aa})


```
