package go_reflect

import (
	"container/list"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/kenshin579/tutorials-go/go-reflect/model"
	"github.com/stretchr/testify/assert"
)

func Example_Type_Value_정보_확인() {
	type Foo struct {
		x int
		y float64
		z string
	}

	foo := Foo{
		x: 1,
		y: 1.0,
		z: "str",
	}

	fmt.Printf("foo: %v(%v)\n", reflect.ValueOf(foo), reflect.TypeOf(foo))
	fmt.Printf("x: %v(%v)\n", reflect.ValueOf(foo.x).Int(), reflect.TypeOf(foo.x))
	fmt.Printf("y: %v(%v)\n", reflect.ValueOf(foo.y).Float(), reflect.TypeOf(foo.y))
	fmt.Printf("z: %v(%v)\n", reflect.ValueOf(foo.z).String(), reflect.TypeOf(foo.z))

	// Output:
	// foo: {1 1 str}(go_reflect.Foo)
	// x: 1(int)
	// y: 1(float64)
	// z: str(string)
}

func Example_Struct_Type_Value_메타_정보_확인() {
	type ArticleRequest struct {
		Title string `json:"title" validate:"required"`
		Body  string `json:"body" validate:"required"`
	}

	a := ArticleRequest{
		Title: "title1",
		Body:  "this is a test",
	}

	uType := reflect.TypeOf(a)
	if fName, ok := uType.FieldByName("Title"); ok {
		fmt.Println(fName.Type, fName.Name, fName.Tag)
	}
	if fId, ok := uType.FieldByName("Body"); ok {
		fmt.Println(fId.Type, fId.Name, fId.Tag)
	}

	// Output:
	// string Title json:"title" validate:"required"
	// string Body json:"body" validate:"required"
}

func Example_Value_변경() {
	languages := []string{"golang", "java", "c++"}
	sliceValue := reflect.ValueOf(languages)
	value := sliceValue.Index(1)
	value.SetString("ruby") // 값 변경을 함
	fmt.Println(languages)

	x := 1
	if v := reflect.ValueOf(x); v.CanSet() { // CanSet으로 변경 가능한 값인지 확인함
		v.SetInt(2) // 호출되지 않음
	}

	fmt.Println(x) // 1

	v := reflect.ValueOf(&x) // pointer
	p := v.Elem()            // Elem() 메서드를 사용하여 값의 주소레 접근하여 다른 값으로 변경함
	p.SetInt(3)
	fmt.Println(x)

	// Output:
	// [golang ruby c++]
	// 1
	// 3

}

func Example_Method_동적_호출() {
	caption := "go is an open source programming language"
	// 1. TitleCase를 바로 호출
	title := TitleCase(caption)
	fmt.Println(title)

	// 2. TitleCase를 동적으로 호출
	titleFuncValue := reflect.ValueOf(TitleCase)
	values := titleFuncValue.Call([]reflect.Value{reflect.ValueOf(caption)})
	title = values[0].String()
	fmt.Println(title)

	// Output:
	// Go Is An Open Source Programming Language
	// Go Is An Open Source Programming Language
}

func TitleCase(s string) string {
	return strings.Title(s)
}

func Example_Len() {
	list1 := list.New() // list1.Len() == 0
	list2 := list.New()
	list2.PushFront(0.5) // list2.Len() == 1

	mapStringInt := map[string]int{"A": 1, "B": 2} // len(mapStringInt) == 2
	str := "one"                                   // len(str) == 3
	sliceInt := []int{5, 0, 4, 1}                  // len(sliceInt) == 4

	fmt.Println(Len(list1), Len(list2), Len(mapStringInt), Len(str), Len(sliceInt))

	// Output:
	// 0 1 2 3 4
}

func Len(x interface{}) int {
	value := reflect.ValueOf(x)
	switch reflect.TypeOf(x).Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		return value.Len()
	default:
		if method := value.MethodByName("Len"); method.IsValid() {
			values := method.Call(nil) // Len 메서드를 동적으로 호출함
			return int(values[0].Int())
		}
	}
	panic(fmt.Sprintf("'%v' does not have a length", x))
}

func Example_구조체_필드_순회하기() {
	cat := &model.Cat{
		Name:  "nabi",
		Age:   5,
		Child: []string{"nyang", "kong"},
	}
	IterateStructField(cat)

	// Output:
	// Name: Name / Type: string / Value: nabi / Tag: name
	// Name: Age / Type: int / Value: 5 / Tag: age
	// Name: Child / Type: []string / Value: [nyang kong] / Tag: child
}

func IterateStructField(object interface{}) {
	elem := reflect.ValueOf(object).Elem()
	fieldNum := elem.NumField()
	for i := 0; i < fieldNum; i++ {
		field := elem.Field(i)            // field
		fieldType := elem.Type().Field(i) // field type
		fieldValue := field.Interface()   // field value 값
		tag := fieldType.Tag.Get("custom")

		fmt.Printf("Name: %s / Type: %s / Value: %v / Tag: %s\n",
			fieldType.Name, fieldType.Type, fieldValue, tag)
	}
}

// todo: 여기서 부터 다시 하면 됨
func Test_타입을_통해_구조체_생성(t *testing.T) {
	cat := model.Cat{
		Name:  "nabi",
		Age:   5,
		Child: []string{"nyang", "kong"},
	}
	cat2 := createStructFromType(cat).(model.Cat)
	cat2 = model.Cat{
		Name: "nyang",
		Age:  1,
	}
	fmt.Println("Cat1: ", cat)
	fmt.Println("Cat2: ", cat2)

}

func createStructFromType(object interface{}) interface{} {
	e := reflect.TypeOf(object)
	return reflect.New(e).Elem().Interface()
}

// http://pyrasis.com/book/GoForTheReallyImpatient/Unit36
func Test_Reflect_Method에_대한_설명(t *testing.T) {
	var f float64 = 1.3
	typ := reflect.TypeOf(f)  // f의 타입 정보를 typ에 저장
	val := reflect.ValueOf(f) // f의 값 정보를 val에 저장

	fmt.Println(typ.Name())                    // float64: 자료형 이름 출력
	fmt.Println(typ.Size())                    // 8: 자료형 크기 출력
	fmt.Println(typ.Kind() == reflect.Float64) // true: 자료형 종류를 알아내어

	// reflect.Float64와 비교
	fmt.Println(typ.Kind() == reflect.Int64) // false: 자료형 종류를 알아내어 reflect.Int64와 비교

	fmt.Println(val.Type()) // float64: 값이 담긴 변수의 자료형 이름 출력

	fmt.Println(val.Kind() == reflect.Float64) // true: 값이 담긴 변수의 자료형 종류를
	// 알아내어 reflect.Float64와 비교
	fmt.Println(val.Kind() == reflect.Int64) // false: 값이 담긴 변수의 자료형 종류를

	// 알아내어 reflect.Int64와 비교
	fmt.Println(val.Float()) // 1.3: 값을 실수형으로 출력
}

// https://cjwoov.tistory.com/16
type Cat struct {
	Name  string   `custom:"name"`
	Age   int      `custom:"age"`
	Child []string `custom:"child"`
}

func TestCatFieldLoop(t *testing.T) {
	cat := &Cat{
		Name:  "nabi",
		Age:   5,
		Child: []string{"nyang", "kong"},
	}
	LoopObjectField(cat)
}

func LoopObjectField(object interface{}) {
	e := reflect.ValueOf(object).Elem()
	fieldNum := e.NumField()
	var childStr string
	for i := 0; i < fieldNum; i++ {
		childStr = ""
		v := e.Field(i)
		t := e.Type().Field(i)
		fmt.Printf("Name: %s / Type: %s / Value: %v / Tag: %s \n",
			t.Name, t.Type, v.Interface(), t.Tag.Get("custom"))
		fmt.Printf("%v\n", v.Kind())
		if v.Kind().String() == "slice" {
			for j := 0; j < v.Len()-1; j++ {
				childStr += v.Index(j).String() + ","
			}
			childStr += v.Index(v.Len() - 1).String()
			fmt.Printf("childStr:%v\n", childStr)
		}
	}
}

func Test_isAllFieldEmptyForStruct(t *testing.T) {
	type person struct {
		Name    string `json:"name"`
		Age     int    `json:"age"`
		Address struct {
			City string `json:"city"`
			Zip  int    `json:"zip"`
		} `json:"address"`
	}

	p1 := person{
		Address: struct {
			City string `json:"city"`
			Zip  int    `json:"zip"`
		}{
			City: "seoul",
		},
	}

	p2 := person{}
	p3 := person{
		Name: "",
		Age:  0,
		Address: struct {
			City string `json:"city"`
			Zip  int    `json:"zip"`
		}{
			City: "",
			Zip:  0,
		},
	}

	assert.False(t, isAllFieldEmpty(p1))
	assert.True(t, isAllFieldEmpty(p2))
	assert.True(t, isAllFieldEmpty(p3))
	assert.False(t, isAllFieldEmpty([]byte(`{"a":"","b":"","c":1}`)))
	assert.False(t, isAllFieldEmpty([]byte(`{"a":"","b":"b","c":0}`)))
	assert.True(t, isAllFieldEmpty([]byte(`{"a":"","b":"","c":0}`)))
	assert.True(t, isAllFieldEmpty([]byte(`{"a":"","b":"","c":0.0}`)))
}

func isAllFieldEmpty(inter any) bool {
	val := reflect.ValueOf(inter)
	if val.IsZero() {
		return true
	}

	switch val.Kind() {
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			zeroValue := reflect.Zero(field.Type())
			if reflect.DeepEqual(field.Interface(), zeroValue.Interface()) {
				continue
			}
			return false
		}
	case reflect.Slice:
		var dataMap map[string]any

		if err := json.Unmarshal(inter.([]byte), &dataMap); err != nil {
			fmt.Printf("fail to unmarshal json. err:%v\n", err)
			return false
		}
		return isMapEmpty(dataMap)
	}

	return true
}

func isMapEmpty(m map[string]any) bool {
	for _, v := range m {
		switch val := v.(type) {
		case string:
			if val != "" {
				return false
			}
		case int:
			if val != 0 {
				return false
			}
		case float64:
			if val != 0 {
				return false
			}
		case bool:
			if val == false {
				return false
			}
		}
	}
	return true
}
