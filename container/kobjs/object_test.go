package kobjs

import (
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/khan-lau/kutils/container/kstrings"
	"github.com/khan-lau/kutils/katomic"
)

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

type MyStruct struct {
	PublicField int
}

func (s MyStruct) PublicMethod0() {
	fmt.Println("Public method called")
}

func Test_reflect(t *testing.T) {
	myInstance := MyStruct{
		PublicField: 42,
	}

	// // 获取对象的反射类型
	// objType := reflect.TypeOf(myInstance)
	// fmt.Printf("objType %s\n", objType.Kind().String())
	// pObjType := reflect.TypeOf(&myInstance)
	// fmt.Printf("pObjType %s\n", pObjType.Kind().String())

	// val := reflect.ValueOf(myInstance)
	// log.Printf("%#v  %#v %#v\n", val, val.Elem(), val.Addr())

	// 获取对象的反射值
	val := reflect.ValueOf(myInstance)
	fmt.Printf("Value of myInstance: %#v\n", val)
	val00 := reflect.ValueOf(val)
	fmt.Printf("Value of myInstance: %#v\n", val00)

	// // 如果 val 是指向对象的指针类型，获取对象本身的反射值
	// if val.Kind() == reflect.Ptr {
	// 	val = val.Elem()
	// }
	// fmt.Printf("Value of myInstance: %#v\n", val)

	// inst, ok := val.Interface().(MyStruct)
	// if ok {
	// 	fmt.Printf("Value of myInstance: %#v\n", inst)
	// }

	// fmt.Printf("Type of myInstance: %s\n", val.Type())

}

// // type User struct {
// // 	Id   int
// // 	Name string
// // 	Age  int
// // }

// // func (u *User) Call() {
// // 	u.Id = 100
// // 	fmt.Printf("%v\n", u)
// // }

// // func GetField(arg interface{}) {
// // 	argType := reflect.TypeOf(arg)
// // 	fmt.Println(argType.NumMethod())
// // 	for i := 0; i < argType.NumMethod(); i++ {
// // 		m := argType.Method(i)
// // 		fmt.Printf("%s: %v\n", m.Name, m.Type)
// // 	}
// // }

func Test_ObjToJson5(t *testing.T) {
	cat := katomic.NewUint32(4)
	// log.Printf("%s", ObjToBeautifulJson5(cat))

	pVal := reflect.ValueOf(cat)
	val := pVal.Elem()

	argType := reflect.TypeOf(val)
	log.Printf("%s \n", argType.Kind())
	if argType.Kind() != reflect.Ptr {
		argType = reflect.TypeOf(val.Interface())
	}

	log.Printf("%s \n", argType.Kind())

	// pVal := reflect.ValueOf(cat)
	// log.Printf("%s\n", pVal.Kind())
	// log.Printf("%#v\n", pVal)
	// if pVal.Kind() != reflect.Ptr {
	// 	log.Printf("%#v\n", pVal)
	// 	pVal = reflect.ValueOf(&cat)

	// 	log.Printf("%#v\n", pVal)
	// }

	// val := pVal.Elem()
	// log.Printf("%#v\n", val)

	// ptrVal := val.Addr()
	// log.Printf("%#v\n", ptrVal)

	methodNums := argType.NumMethod()
	for i := 0; i < methodNums; i++ {
		m := argType.Method(i)

		method := pVal.MethodByName(m.Name)
		log.Printf("name:%s - %s %v %v\n", m.Name, method.Kind(), method.IsValid(), method)

		// 	// 获取方法的类型
		// 	methodType := method.Type()

		// 	log.Printf("%s: \"%v\"    %#v\n", m.Name, methodType, pVal)

		// 	//log.Printf("%s", funcToJson5(m.Name, pVal))
		// 	// sb.WriteString(fmt.Sprintf("%s%s: \"%s\"", strings.Repeat(sdent, (level+1)*2), m.Name, funcToJson5(m.Name, val)))
		// 	// if i < methodNums-1 {
		// 	// 	sb.WriteString(",\n")
		// 	// }
	}

	// myInstance := MyStruct{
	// 	PublicField:  42,
	// 	privateField: "secret",
	// }

	// // 获取对象的反射类型
	// objType := reflect.TypeOf(myInstance)
	// // fmt.Printf("objType %s\n", objType.Kind().String())
	// // pObjType := reflect.TypeOf(&myInstance)
	// // fmt.Printf("pObjType %s\n", pObjType.Kind().String())

	// // 遍历成员字段
	// for i := 0; i < objType.NumField(); i++ {
	// 	field := objType.Field(i)
	// 	fmt.Printf("Field Name: %s, Type: %v\n", field.Name, field.Type)
	// }
	// fmt.Println()

	// // 遍历方法
	// pObjType := objType
	// if objType.Kind() != reflect.Ptr {
	// 	pObjType = reflect.TypeOf(&myInstance)
	// }
	// for i := 0; i < pObjType.NumMethod(); i++ {
	// 	method := pObjType.Method(i)
	// 	fmt.Printf("* Method Name: %s\n", method.Name)
	// }

	// fmt.Println()

	// // value := reflect.ValueOf(&myInstance)
	// // // 调用公有方法
	// // methodValue0 := value.MethodByName("PublicMethod0")
	// // methodValue0.Call(nil)

	// // // 获取私有方法的反射值
	// // methodValue1 := value.MethodByName("privateMethod")
	// // if methodValue1.IsValid() {
	// // 	methodValue1.Call(nil)
	// // } else {
	// // 	t.Errorf("%s", "private Method not found")
	// // }

	// // user := User{1, "ding", 18}
	// // GetField(&user) // 传递指针类型的实例

}
