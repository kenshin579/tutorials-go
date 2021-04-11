package go_reflect

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/kenshin579/tutorials-go/go-reflect/model"
)

//http://pyrasis.com/book/GoForTheReallyImpatient/Unit36
func Test_Reflect_Method에_대한_설명(t *testing.T) {
	var f float64 = 1.3
	typ := reflect.TypeOf(f)  // f의 타입 정보를 typ에 저장
	val := reflect.ValueOf(f) // f의 값 정보를 val에 저장

	fmt.Println(typ.Name())                    // float64: 자료형 이름 출력
	fmt.Println(typ.Size())                    // 8: 자료형 크기 출력
	fmt.Println(typ.Kind() == reflect.Float64) // true: 자료형 종류를 알아내어

	// reflect.Float64와 비교
	fmt.Println(typ.Kind() == reflect.Int64) // false: 자료형 종류를 알아내어 reflect.Int64와 비교

	fmt.Println(val.Type())                    // float64: 값이 담긴 변수의 자료형 이름 출력
	fmt.Println(val.Kind() == reflect.Float64) // true: 값이 담긴 변수의 자료형 종류를

	// 알아내어 reflect.Float64와 비교
	fmt.Println(val.Kind() == reflect.Int64) // false: 값이 담긴 변수의 자료형 종류를

	// 알아내어 reflect.Int64와 비교
	fmt.Println(val.Float()) // 1.3: 값을 실수형으로 출력
}

func Example_구조체_필드_순회하기() {
	cat := &model.Cat{
		Name:  "nabi",
		Age:   5,
		Child: []string{"nyang", "kong"},
	}
	IterateStructField(cat)

	//Output:
	//Name: Name / Type: string / Value: nabi / Tag: name
	//Name: Age / Type: int / Value: 5 / Tag: age
	//Name: Child / Type: []string / Value: [nyang kong] / Tag: child
}

func IterateStructField(object interface{}) {
	elem := reflect.ValueOf(object).Elem()
	fieldNum := elem.NumField()
	for i := 0; i < fieldNum; i++ {
		field := elem.Field(i)            //field
		fieldType := elem.Type().Field(i) //field type
		fieldValue := field.Interface()   //field value 값
		tag := fieldType.Tag.Get("custom")

		fmt.Printf("Name: %s / Type: %s / Value: %v / Tag: %s\n",
			fieldType.Name, fieldType.Type, fieldValue, tag)
	}
}

//todo: 여기서 부터 다시 하면 됨
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
