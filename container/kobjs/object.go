package kobjs

import (
	"fmt"
	"reflect"
	"strings"
)

func ObjToJson5(obj any) string {
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	var sb strings.Builder
	if val.Kind() == reflect.Struct {
		sb.WriteString("{")
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			fieldName := val.Type().Field(i).Name
			if !fieldPublic(fieldName) {
				continue
			}
			if field.Kind() == reflect.Ptr {
				field = field.Elem()
			}

			switch field.Kind() {
			case reflect.Struct:
				sb.WriteString(fmt.Sprintf("%s: %s", fieldName, ObjToJson5(field.Interface())))

			case reflect.Slice, reflect.Array:
				sb.WriteString(fmt.Sprintf("%s: [", fieldName))
				for j := 0; j < field.Len(); j++ {

					sb.WriteString(ObjToJson5(field.Index(j).Interface()))
					if j < field.Len()-1 {
						sb.WriteString(",")
					}
				}
				sb.WriteString("]")
			case reflect.String:
				str := fmt.Sprintf("%s", field.Interface())
				str = strings.ReplaceAll(str, "\\", "\\\\")
				str = strings.ReplaceAll(str, "\n", "\\n")
				sb.WriteString(fmt.Sprintf("%s: \"%v\"", fieldName, str))

			default:
				sb.WriteString(fmt.Sprintf("%s: %v", fieldName, field.Interface()))
			}

			if i < val.NumField()-1 {
				sb.WriteString(",")
			}
		}
		sb.WriteString("}")

	} else if val.Kind() == reflect.Array || val.Kind() == reflect.Slice {
		sb.WriteString("[")
		for j := 0; j < val.Len(); j++ {
			sb.WriteString(ObjToJson5(val.Index(j).Interface()))
		}
		sb.WriteString("]")
	} else if val.Kind() == reflect.String {
		str := fmt.Sprintf("%s", val.Interface())
		str = strings.ReplaceAll(str, "\\", "\\\\")
		str = strings.ReplaceAll(str, "\n", "\\n")
		sb.WriteString(fmt.Sprintf("\"%v\"", str))
	} else {
		sb.WriteString(fmt.Sprintf("%v", val.Interface()))
	}

	return sb.String()
}

func ObjToBeautifulJson5(obj any) string {
	return objToBeautifulJson5(obj, -1)
}

// func objectToBeautifulJson5(object any, level int) (string, error) {
// 	pValue, ok := object.(reflect.Value)
// 	if !ok {
// 		pValue = reflect.ValueOf(object)
// 	}
// 	if pValue.Kind() != reflect.Ptr {
// 		pValue = reflect.ValueOf(&object)
// 	}
// }

func objToBeautifulJson5(obj any, level int) string {
	sdent := "  "
	// argType := reflect.TypeOf(obj)
	// if argType.Kind() != reflect.Ptr {
	// 	tmpVal, ok := obj.(reflect.Value)
	// 	if ok {
	// 		tmp := tmpVal.Interface()
	// 		argType = reflect.TypeOf(&tmp)
	// 	} else {
	// 		argType = reflect.TypeOf(&obj)
	// 	}
	// }

	val := reflect.ValueOf(obj)
	// log.Printf("  %#v\n", val)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if level < -1 {
		level = -1
	}

	level++

	var sb strings.Builder
	if val.Kind() == reflect.Struct {
		sb.WriteString(fmt.Sprintf("%s{\n", strings.Repeat(sdent, level*2)))
		publicFieldNums := 0 // public成员遍历数量
		fieldNums := val.NumField()
		for i := 0; i < fieldNums; i++ {
			field := val.Field(i)
			fieldName := val.Type().Field(i).Name
			// log.Printf("%s\n", fieldName)
			// 判断字段的首字母是否是大写
			if !fieldPublic(fieldName) {
				continue
			}

			publicFieldNums++

			if field.Kind() == reflect.Ptr {
				field = field.Elem()
			}

			switch field.Kind() {
			case reflect.Struct:
				sb.WriteString(fmt.Sprintf("%s%s: \n%s", strings.Repeat(sdent, (level+1)*2), fieldName, objToBeautifulJson5(field.Interface(), level)))

			case reflect.Slice, reflect.Array:
				sb.WriteString(fmt.Sprintf("%s%s:[\n", strings.Repeat(sdent, (level+1)*2), fieldName))
				for j := 0; j < field.Len(); j++ {

					sb.WriteString(objToBeautifulJson5(field.Index(j).Interface(), level+1))
					if j < field.Len()-1 {
						sb.WriteString(",\n")
					}
				}
				sb.WriteString(fmt.Sprintf("\n%s]", strings.Repeat(sdent, (level+1)*2)))
			case reflect.String:

				str := fmt.Sprintf("%s", field.Interface())
				str = strings.ReplaceAll(str, "\\", "\\\\")
				str = strings.ReplaceAll(str, "\n", "\\n")
				sb.WriteString(fmt.Sprintf("%s%s: \"%v\"", strings.Repeat(sdent, (level+1)*2), fieldName, str))

			default:
				sb.WriteString(fmt.Sprintf("%s%s: %v", strings.Repeat(sdent, (level+1)*2), fieldName, field.Interface()))
			}

			if i < fieldNums-1 {
				sb.WriteString(",\n")
			}
		}

		// methodNums := argType.NumMethod()
		// for i := 0; i < methodNums; i++ {
		// 	m := argType.Method(i)
		// 	if publicFieldNums > 0 {
		// 		sb.WriteString(",\n")
		// 	}
		// 	sb.WriteString(fmt.Sprintf("%s%s: \"%s\"", strings.Repeat(sdent, (level+1)*2), m.Name, funcToJson5(m.Name, val)))
		// 	if i < methodNums-1 {
		// 		sb.WriteString(",\n")
		// 	}
		// }

		sb.WriteString(fmt.Sprintf("\n%s}", strings.Repeat(sdent, level*2)))

	} else if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		sb.WriteString("[\n")
		for j := 0; j < val.Len(); j++ {
			sb.WriteString(objToBeautifulJson5(val.Index(j).Interface(), level+1))
		}
		sb.WriteString("\n]\n")
	} else if val.Kind() == reflect.String {
		str := fmt.Sprintf("%s", val.Interface())
		str = strings.ReplaceAll(str, "\\", "\\\\")
		str = strings.ReplaceAll(str, "\n", "\\n")
		sb.WriteString(fmt.Sprintf("%s\"%v\"", strings.Repeat(sdent, (level+1)*2), str))
	} else {
		sb.WriteString(fmt.Sprintf("%s%v", strings.Repeat(sdent, (level+1)*2), val.Interface()))
	}

	return sb.String()
}

// func funcToJson5(name string, val reflect.Value) string {
// 	pVal := reflect.ValueOf(val)
// 	if pVal.Kind() != reflect.Ptr {
// 		pVal = val.Addr()
// 		// 	// // log.Printf("%#v\n", tmp)
// 		// pVal = reflect.ValueOf(&tmp)
// 		// 	// pVal = reflect.ValueOf(&val)
// 	}

// 	log.Printf("func %s, %s pVal: %#v", name, pVal.Kind(), pVal)

// 	var sb strings.Builder

// 	// 获取 ProcessComplexObject 方法的反射值
// 	method := pVal.MethodByName(name)
// 	if method.IsValid() {
// 		// 获取方法的类型
// 		methodType := method.Type()

// 		sb.WriteString("func(")

// 		// 获取参数个数
// 		numArgs := methodType.NumIn()
// 		// 获取参数类型
// 		for i := 0; i < numArgs; i++ {
// 			argType := methodType.In(i)
// 			sb.WriteString(argType.String())
// 			if i < numArgs-1 {
// 				sb.WriteString(",")
// 			}
// 		}

// 		sb.WriteString(")")

// 		// 获取返回值个数
// 		numResults := methodType.NumOut()
// 		if numResults > 0 {
// 			sb.WriteString("(")
// 			// 获取返回值类型
// 			for i := 0; i < numResults; i++ {
// 				resultType := methodType.Out(i)
// 				sb.WriteString(resultType.String())
// 				if i < numResults-1 {
// 					sb.WriteString(",")
// 				}
// 			}
// 			sb.WriteString(")")
// 		}
// 	}

// 	return sb.String()

// 	// // 准备参数
// 	// args := []reflect.Value{reflect.ValueOf(10), reflect.ValueOf(1.234)}

// 	// // 调用方法
// 	// results := method.Call(args)

// 	// // 解析返回值
// 	// resultString := results[0].Interface().(string)
// 	// resultCount := results[1].Interface().(int)

// 	// return ""
// }

// func getMethod(obj any) {
// 	argType := reflect.TypeOf(obj)
// 	if argType.Kind() != reflect.Ptr {
// 		argType = reflect.TypeOf(&obj)
// 	}
// 	fmt.Println(argType.NumMethod())
// 	for i := 0; i < argType.NumMethod(); i++ {
// 		m := argType.Method(i)
// 		fmt.Printf("%s: %v\n", m.Name, m.Type)
// 	}
// }

func fieldPublic(fieldName string) bool {
	if len(fieldName) > 0 && fieldName[0] >= 'A' && fieldName[0] <= 'Z' {
		return true
	}

	return false
}
