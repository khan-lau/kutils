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