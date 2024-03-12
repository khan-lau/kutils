package kobjs

import (
	"container/list"
	"fmt"
	"log"
	"reflect"
	"strings"
	"unsafe"
)

type flag uintptr // reflect/value.go:flag

type flagROTester struct {
	A   int
	a   int // reflect/value.go:flagStickyRO   // 用于检测hack代码支持情况, 勿删
	int     // reflect/value.go:flagEmbedRO    // 用于检测hack代码支持情况, 勿删
	// Note: flagRO = flagStickyRO | flagEmbedRO
}

var (
	flagOffset               uintptr
	maskFlagRO               flag
	hasExpectedReflectStruct bool
)

// 检测当前golang环境是否支持hack, 并初始化各种flag
func initUnsafe() {
	if field, ok := reflect.TypeOf(reflect.Value{}).FieldByName("flag"); ok {
		flagOffset = field.Offset
	} else {
		log.Println("go-describe: exposeInterface() is disabled because the " +
			"reflect.Value struct no longer has a flag field. Please open an " +
			"issue at https://github.com/kstenerud/go-describe/issues")
		hasExpectedReflectStruct = false
		return
	}

	rv := reflect.ValueOf(flagROTester{})
	getFlag := func(v reflect.Value, name string) flag {
		return flag(reflect.ValueOf(v.FieldByName(name)).FieldByName("flag").Uint())
	}
	flagRO := (getFlag(rv, "a") | getFlag(rv, "int")) ^ getFlag(rv, "A")
	maskFlagRO = ^flagRO

	if flagRO == 0 {
		log.Println("go-describe: exposeInterface() is disabled because the " +
			"reflect flag type no longer has a flagEmbedRO or flagStickyRO bit. " +
			"Please open an issue at https://github.com/kstenerud/go-describe/issues")
		hasExpectedReflectStruct = false
		return
	}

	hasExpectedReflectStruct = true
}

func canExposeInterface() bool {
	return hasExpectedReflectStruct
}

// 如果支持hack, 导出私有变量
func exposeInterface(v reflect.Value) interface{} {
	if canExposeInterface() {
		pFlag := (*flag)(unsafe.Pointer(uintptr(unsafe.Pointer(&v)) + flagOffset))
		*pFlag &= maskFlagRO
		return v.Interface()
	} else {
		return v.Interface()
	}
}

// 将对象的成员变量与成员函数 通过beautiful json5的格式打印
//   - @return string 格式化的json5字符串
func ObjectToJson5(obj interface{}) string {
	initUnsafe()

	var sb strings.Builder
	objectToJson5("", obj, &sb, 0, "  ")

	return sb.String()
}

// 将对象树导出为格式化过的JSON5字符串
//   - @param string           key   对象变量名称
//   - @param any              obj   对象实例, 也可以是指针
//   - @param *strings.Builder sb    字符串缓冲
//   - @param uint             level 对象深度
//   - @param string           ident 格式化前导字符串, 通常为 "  " 或 "\t"
func objectToJson5(key string, obj interface{}, sb *strings.Builder, level uint, ident string) {
	val := mustBeValue(obj)
	if len(key) > 0 {
		sb.WriteString(fmt.Sprintf("\n%s%s:", strings.Repeat(ident, int(level)), key))
	}
	var ptr *reflect.Value
	if ptr == nil {
		tmpNewVal := reflect.New(val.Type()) // 根据val的类型, 新建一个ptr类型的Value
		ptr = &tmpNewVal
	}
	switch val.Kind() {
	case reflect.Struct:
		{
			rt := val.Type()
			rtStr := fmt.Sprintf("%s.%s \n", rt.PkgPath(), rt.Name())
			// fmt.Printf("%s\n", rtStr)
			if strings.Contains(rtStr, "kutils/container/klists.KList[") {

				getType := val.Type() //reflect.TypeOf(ptr)
				met, ok := getType.MethodByName("ToJson5")
				if !ok {
					sb.WriteString(fmt.Sprintf("%v", val.Interface()))
					//panic("method not exist.")
				} else {
					sb.WriteString("[\n")
					level++
					strTmp := strings.Repeat(ident, int(level))
					var args = []reflect.Value{val, reflect.ValueOf(strTmp)}
					retv := met.Func.Call(args)
					reti := retv[0].Interface() // interface{} 类型
					ret, _ := reti.(string)

					sb.WriteString(ret)
					level--
					sb.WriteString(fmt.Sprintf("%s]", strings.Repeat(ident, int(level))))
				}

				return
			} else if strings.Contains(rtStr, "container/list.List") {
				v := val.Interface()
				l, _ := v.(list.List)
				var next *list.Element
				i := 0
				sb.WriteString("[\n")
				level++
				for e := l.Front(); e != nil; e = next {
					next = e.Next()
					sb.WriteString(fmt.Sprintf("%s%#v", strings.Repeat(ident, int(level)), e.Value))
					if i < l.Len()-1 {
						sb.WriteString(",")
					}
					i++
					sb.WriteString("\n")
				}
				level--
				sb.WriteString(fmt.Sprintf("%s]", strings.Repeat(ident, int(level))))
				return
			}

			sb.WriteString("{")
			first := true
			publicFieldNums := 0 // public成员遍历数量
			fieldLen := val.NumField()
			level++
			for i := 0; i < fieldLen; i++ {
				subObj := val.Field(i)
				fieldName := val.Type().Field(i).Name

				if !canExposeInterface() && !fieldPublic(fieldName) {
					continue
				}

				if first {
					first = false
				} else {
					sb.WriteString(",")
				}
				publicFieldNums++

				objectToJson5(fieldName, subObj, sb, level, ident)
			}

			// 遍历成员方法, 暂时只支持Public方法
			if ptr.IsValid() {
				if publicFieldNums > 0 {
					sb.WriteString(",")
				}
				sb.WriteString(fmt.Sprintf("\n%s%s:{\n", strings.Repeat(ident, int(level)), "func"))
				methodMap := getMethodMap(*ptr)
				itemLen := len(methodMap)
				i := 0
				level++

				for k, v := range methodMap {
					sb.WriteString(fmt.Sprintf("%s%s:\"%s\"", strings.Repeat(ident, int(level)), k, v))
					if i < itemLen-1 {
						sb.WriteString(",")
					}
					sb.WriteString("\n")
					i++
				}
				level--
				sb.WriteString(fmt.Sprintf("%s}\n", strings.Repeat(ident, int(level))))
			}
			level--
			sb.WriteString(fmt.Sprintf("%s}", strings.Repeat(ident, int(level))))
		}
	case reflect.Array, reflect.Slice:
		{
			sb.WriteString("[\n")
			itemLen := val.Len()
			level++
			for i := 0; i < itemLen; i++ {
				objectToJson5("", val.Index(i), sb, level, ident)
				if i < itemLen-1 {
					sb.WriteString(",")
				}
				sb.WriteString("\n")
			}
			level--
			sb.WriteString(fmt.Sprintf("%s]", strings.Repeat(ident, int(level))))
		}
	case reflect.String:
		str := fmt.Sprintf("%s", exposeInterface(val))
		str = strings.ReplaceAll(str, "\\", "\\\\")
		str = strings.ReplaceAll(str, "\n", "\\n")
		if len(key) == 0 {
			sb.WriteString(fmt.Sprintf("%s\"%v\"", strings.Repeat(ident, int(level)), str))
		} else {
			sb.WriteString(fmt.Sprintf("\"%v\"", str))
		}
	case reflect.Func:
		str := funcVarToJson5(key, *ptr)
		sb.WriteString(fmt.Sprintf("\"func%s\"", str))
	case reflect.Map:
		sb.WriteString("{")
		keys := val.MapKeys()
		length := len(keys)
		i := 0
		level++
		for _, key := range keys {
			strct := val.MapIndex(key)
			objectToJson5(fmt.Sprintf("%v", key.Interface()), strct, sb, level, ident)
			if i < length-1 {
				sb.WriteString(",")
			}
			i++
		}
		level--
		sb.WriteString(fmt.Sprintf("\n%s}", strings.Repeat(ident, int(level))))
	default:
		v1 := exposeInterface(val)
		if nil == v1 {
			sb.WriteString(fmt.Sprintf("%v", "null"))
		} else {
			sb.WriteString(fmt.Sprintf("%v", exposeInterface(val)))
		}
	}
}

// 将对象的成员变量与成员函数 json5的格式打印
//   - @return string json5字符串
func ObjectDump(obj interface{}) string {
	initUnsafe()

	var sb strings.Builder
	objectDump("", obj, &sb)

	return sb.String()
}

// 将对象树导出为JSON5字符串
//   - @param string           key   对象变量名称
//   - @param any              obj   对象实例, 也可以是指针
//   - @param *strings.Builder sb    字符串缓冲
func objectDump(key string, obj interface{}, sb *strings.Builder) {
	val := mustBeValue(obj)
	if len(key) > 0 {
		sb.WriteString(fmt.Sprintf("%s:", key))
	}
	var ptr *reflect.Value
	if ptr == nil {
		tmpNewVal := reflect.New(val.Type()) // 根据val的类型, 新建一个ptr类型的Value
		ptr = &tmpNewVal
	}
	switch val.Kind() {
	case reflect.Struct:
		{
			rt := val.Type()
			rtStr := fmt.Sprintf("%s.%s \n", rt.PkgPath(), rt.Name())
			//fmt.Printf("%s\n", rtStr)
			if strings.Contains(rtStr, "kutils/container/klists.KList[") {

				getType := val.Type() //reflect.TypeOf(ptr)
				met, ok := getType.MethodByName("ToString")
				if !ok {
					sb.WriteString(fmt.Sprintf("%v", val.Interface()))
					//panic("method not exist.")
				} else {
					var args = []reflect.Value{val}
					retv := met.Func.Call(args)
					reti := retv[0].Interface() // interface{} 类型
					ret, _ := reti.(string)
					sb.WriteString(ret)
				}

				return
			} else if strings.Contains(rtStr, "container/list.List") {
				v := val.Interface()
				l, _ := v.(list.List)
				var next *list.Element
				i := 0
				sb.WriteString("[")
				for e := l.Front(); e != nil; e = next {
					next = e.Next()
					sb.WriteString(fmt.Sprintf("%#v", e.Value))
					if i < l.Len()-1 {
						sb.WriteString(",")
					}
					i++
				}
				sb.WriteString("]")
				return
			}

			sb.WriteString("{")
			first := true
			publicFieldNums := 0 // public成员遍历数量
			fieldLen := val.NumField()
			for i := 0; i < fieldLen; i++ {
				subObj := val.Field(i)
				fieldName := val.Type().Field(i).Name

				if !canExposeInterface() && !fieldPublic(fieldName) {
					continue
				}

				if first {
					first = false
				} else {
					sb.WriteString(",")
				}
				publicFieldNums++
				objectDump(fieldName, subObj, sb)
			}

			// 遍历成员方法, 暂时只支持Public方法
			if ptr.IsValid() {
				if publicFieldNums > 0 {
					sb.WriteString(",")
				}
				sb.WriteString(fmt.Sprintf("%s:", "func"))
				sb.WriteString("{")
				getMethod(*ptr, sb)
				sb.WriteString("}")
			}
			sb.WriteString("}")
		}
	case reflect.Array, reflect.Slice:
		{
			sb.WriteString("[")
			itemLen := val.Len()
			for i := 0; i < itemLen; i++ {
				objectDump("", val.Index(i), sb)
				if i < itemLen-1 {
					sb.WriteString(",")
				}
			}
			sb.WriteString("]")
		}
	case reflect.String:
		str := fmt.Sprintf("%s", exposeInterface(val))
		str = strings.ReplaceAll(str, "\\", "\\\\")
		str = strings.ReplaceAll(str, "\n", "\\n")
		sb.WriteString(fmt.Sprintf("\"%v\"", str))
	case reflect.Func:
		str := funcVarToJson5(key, *ptr)
		sb.WriteString(fmt.Sprintf("\"func%s\"", str))
	case reflect.Map:
		sb.WriteString("{")
		keys := val.MapKeys()
		length := len(keys)
		i := 0
		for _, key := range keys {
			strct := val.MapIndex(key)
			objectDump(fmt.Sprintf("%v", key.Interface()), strct, sb)
			if i < length-1 {
				sb.WriteString(",")
			}
			i++
		}
		sb.WriteString("}")
	default:
		v1 := exposeInterface(val)
		if nil == v1 {
			sb.WriteString(fmt.Sprintf("%v", "null"))
		} else {
			sb.WriteString(fmt.Sprintf("%v", exposeInterface(val)))
		}
	}
}

// 判断对象是否为指定类型
//   - @temp_param T                类型
//   - @param      interface{}  obj 变量
//   - @return bool 是否为指定类型
func InstanceOf[T interface{}](obj interface{}) bool {
	switch obj.(type) {
	case T:
		return true
	}
	return false
}

func mustBeValue(obj any) reflect.Value {
	var val reflect.Value
	if InstanceOf[reflect.Value](obj) { // 反射对象
		val, _ = obj.(reflect.Value)
	} else { // 普通对象
		val = reflect.ValueOf(obj)
	}
	switch val.Kind() {
	case reflect.Ptr:
		return val.Elem()
	default:
		return val
	}
}

// // - 返回 obj 对象的 反射PtrValue 和 reflect.Value
// //   - 如果obj 是普通类型, 返回 反射PtrValue 和 reflect.Value
// //   - 如果obj是反射PtrValue , 返回PtrValue自身和其指向的 reflect.Value
// //   - 如果obj是reflect.Value, 返回新构建的PtrValue和其指向的 reflect.Value
// func mustBePtrValue(obj any) (reflect.Value, reflect.Value) {
// 	var ptr, val reflect.Value
// 	if InstanceOf[reflect.Value](obj) { // 反射对象
// 		val, _ = obj.(reflect.Value)
// 	} else { // 普通对象
// 		val = reflect.ValueOf(obj)
// 	}
//
// 	switch val.Kind() {
// 	case reflect.Struct, reflect.Func:
// 		// val = val.Elem()
// 		// val = reflect.Indirect(val) // 获取ptr类型的Value指向的值
// 		ptr = reflect.New(val.Type()) // 根据val的类型, 新建一个ptr类型的Value
// 		return ptr, ptr.Elem()
// 	case reflect.Ptr:
// 		return val, val.Elem()
// 	default:
// 		var tmp interface{}
// 		return reflect.ValueOf(&tmp), val
// 	}
// }

// 将函数类型变量转换为字符串格式
func funcVarToJson5(name string, val reflect.Value) string {
	var sb strings.Builder

	valv := val.Elem()
	tpv := valv.Type()         // 获取参数个数
	numArgs := tpv.NumIn()     // 获取参数个数
	numResults := tpv.NumOut() // 获取返回值个数

	sb.WriteString("(")
	for i := 0; i < numArgs; i++ {
		argType := tpv.In(i)
		sb.WriteString(argType.String())
		if i < numArgs-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString(")")

	if numResults > 0 {
		sb.WriteString("(")
		// 获取返回值类型
		for i := 0; i < numResults; i++ {
			resultType := tpv.Out(i)
			sb.WriteString(resultType.String())
			if i < numResults-1 {
				sb.WriteString(", ")
			}
		}
		sb.WriteString(")")
	}

	return sb.String()
}

// 将成员函数转换为字符串格式
func funcToJson5(name string, val reflect.Value) string {
	var sb strings.Builder

	method := val.MethodByName(name) // 获取 ProcessComplexObject 方法的反射值
	if method.IsValid() {
		methodType := method.Type() // 获取方法的类型
		sb.WriteString("func(")

		numArgs := methodType.NumIn() // 获取参数个数
		// 获取参数类型
		for i := 0; i < numArgs; i++ {
			argType := methodType.In(i)
			sb.WriteString(argType.String())
			if i < numArgs-1 {
				sb.WriteString(", ")
			}
		}

		sb.WriteString(")")

		// 获取返回值个数
		numResults := methodType.NumOut()
		if numResults > 0 {
			sb.WriteString("(")
			// 获取返回值类型
			for i := 0; i < numResults; i++ {
				resultType := methodType.Out(i)
				sb.WriteString(resultType.String())
				if i < numResults-1 {
					sb.WriteString(", ")
				}
			}
			sb.WriteString(")")
		}
	} else {
		fmt.Printf("error: %s\n", "method is invalid")
	}

	return sb.String()
}

// 获取对象的所有成员函数
//   - @param reflect.Value val 反射对象
//   - @return map[string]string map<函数名称, 函数签名>
func getMethodMap(val reflect.Value) map[string]string {
	result := make(map[string]string)
	argType := reflect.TypeOf(val.Interface())
	itemLen := argType.NumMethod()
	for i := 0; i < itemLen; i++ {
		m := argType.Method(i)
		result[m.Name] = funcToJson5(m.Name, val)
	}
	return result
}

// 获取对象的所有成员函数
func getMethod(val reflect.Value, sb *strings.Builder) {
	argType := reflect.TypeOf(val.Interface())

	itemLen := argType.NumMethod()
	for i := 0; i < itemLen; i++ {
		m := argType.Method(i)
		sb.WriteString(fmt.Sprintf("%s:\"%s\"", m.Name, funcToJson5(m.Name, val)))
		if i < itemLen-1 {
			sb.WriteString(",")
		}
	}
}

// 根据成员变量的名称判断是否为Public成员
//   - @param string fieldName 变量名
//   - @return bool
func fieldPublic(fieldName string) bool {
	if len(fieldName) > 0 && fieldName[0] >= 'A' && fieldName[0] <= 'Z' {
		return true
	}

	return false
}
