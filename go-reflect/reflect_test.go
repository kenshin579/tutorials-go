package go_reflect

import (
	"container/list"
	"fmt"
	"reflect"
	"strings"
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

	//Output:
	//foo: {1 1 str}(go_reflect.Foo)
	//x: 1(int)
	//y: 1(float64)
	//z: str(string)
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

	//Output:
	//string Title json:"title" validate:"required"
	//string Body json:"body" validate:"required"
}

func Example_Value_변경() {
	languages := []string{"golang", "java", "c++"}
	sliceValue := reflect.ValueOf(languages)
	value := sliceValue.Index(1)
	value.SetString("ruby") //값 변경을 함
	fmt.Println(languages)

	x := 1
	if v := reflect.ValueOf(x); v.CanSet() { //CanSet으로 변경 가능한 값인지 확인함
		v.SetInt(2) // 호출되지 않음
	}

	fmt.Println(x) // 1

	v := reflect.ValueOf(&x) //pointer
	p := v.Elem()            //Elem() 메서드를 사용하여 값의 주소레 접근하여 다른 값으로 변경함
	p.SetInt(3)
	fmt.Println(x)

	//Output:
	//[golang ruby c++]
	//1
	//3

}

func Example_Method_동적_호출() {
	caption := "go is an open source programming language"
	//1. TitleCase를 바로 호출
	title := TitleCase(caption)
	fmt.Println(title)

	//2. TitleCase를 동적으로 호출
	titleFuncValue := reflect.ValueOf(TitleCase)
	values := titleFuncValue.Call([]reflect.Value{reflect.ValueOf(caption)})
	title = values[0].String()
	fmt.Println(title)

	//Output:
	//Go Is An Open Source Programming Language
	//Go Is An Open Source Programming Language
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

	//Output:
	//0 1 2 3 4
}

func Len(x interface{}) int {
	value := reflect.ValueOf(x)
	switch reflect.TypeOf(x).Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		return value.Len()
	default:
		if method := value.MethodByName("Len"); method.IsValid() {
			values := method.Call(nil) //Len 메서드를 동적으로 호출함
			return int(values[0].Int())
		}
	}
	panic(fmt.Sprintf("'%v' does not have a length", x))
}