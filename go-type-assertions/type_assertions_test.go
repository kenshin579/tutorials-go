package go_type_assertions

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Example_TypeAssertion_Empty_Interface() {
	var i interface{} = "foo"

	var s = i.(string)
	fmt.Println(s)

	s, ok := i.(string)
	fmt.Println(s, ok)

	n, ok := i.(int) // panic 대신 ok:false, n는 zero 값으로 반환한다
	fmt.Println(n, ok)

	// Output:
	// foo
	// foo true
	// 0 false
}

func Example_TypeAssertion_Empty_Interface_Illegal() {
	var i interface{} = "foo"

	var n = i.(int) // panic: interface conversion: interface {} is string, not int
	fmt.Println(n)

	// Output:
}

type Person interface {
	getName() string
}

type Phone interface {
	getPhone() string
}

type Student struct {
	name  string
	age   int
	phone string
}

func (c Student) getName() string {
	return c.name
}

func (c Student) getPhone() string {
	return c.phone
}

type Animal interface {
	walk()
}

func Example_TypeAssertion_인터페이스_데이터_타입_Student_값을_가져온다() {
	var p Person = Student{"Frank", 13, "1111"}
	fmt.Println(p.getName())

	s := p.(Student) // Person -> Student - student의 실제 값을 가져온다.
	fmt.Println(s.getName())
	fmt.Println(s.getPhone())

	// Output:
	// Frank
	// Frank
	// 1111
}

func Example_TypeAssertion_다른_인터페이스로_값을_가져온다() {
	var p Person = Student{"Frank", 13, "1111"}

	ph := p.(Phone) // Person -> Phone
	fmt.Println(ph.getPhone())

	// Output:
	// 1111
}

// 타입 T가 인터페이스를 구현하고 있지 않기 때문에 컴파일 에러가 발생한다
func Example_TypeAssertion_인터페이스가_타입_T의_동적_값을_소유하지_않을_경우_컴파일_에러가_발생한다() {
	// var p Person = Student{"Frank", 13, "1111"}
	// value := p.(string) //impossible type assertion: string does not implement person (missing getName method)
	// fmt.Printf("%v, %T\n", value, value)

	// Output:
}

func Example_TypeAssertion_인터페이스가_타입_T의_실제_값을_가지고_있지_않는_경우_panic이_발생한다() {
	var p Person = nil
	// value := p.(Student) //panic: interface conversion: go_type_assertions.Person is nil, not go_type_assertions.Student
	value, ok := p.(Student)
	fmt.Printf("(%v, %T), ok: %v\n", value, value, ok)

	// Output:
	// ({ 0 }, go_type_assertions.Student), ok: false
}

func Example_TypeAssertion_다른_인터페이스가_타입_T를_구현하지_않고_있으면_panic이_발생한다() {
	var p Person = Student{"Frank", 13, "1111"}
	// value := p.(Animal) //panic: interface conversion: go_type_assertions.Student is not go_type_assertions.Animal: missing method walk
	value, ok := p.(Animal)
	fmt.Printf("(%v, %T) %v\n", value, value, ok)

	// Output:
	// (<nil>, <nil>) false
}

type ListStudent []Student

func TestTypeConversion이_안되는_케이스(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	var studentList interface{} = []Student{
		{
			name: "name-1",
			age:  10,
		},
	}

	result := studentList.(ListStudent) // panic: interface conversion: interface {} is []go_type_assertions.Student, not go_type_assertions.ListStudent
	fmt.Printf("result:%+v", result)
}

type Event string

var (
	UpdateEvent = Event("UPDATE")
)

/*
type assertion이 안되는 케이스

The issue in your code is that you're attempting to perform a type assertion from str to model.Event, but the type of str is interface{} and not model.Event.
Type assertions work when the underlying value has the specified type or implements the specified interface.
*/
func Test_Type_Assertion_잘_못사용하는_케이스(t *testing.T) {
	var str interface{}
	str = "UPDATE"

	event, ok := str.(Event)
	assert.False(t, ok)
	assert.NotEqualf(t, UpdateEvent, event, "type assertion을 할수가 없다")

	// interface{}을 Event로 변경하려면
	s := str.(string)
	e := Event(s)
	assert.Equal(t, UpdateEvent, e)

}
