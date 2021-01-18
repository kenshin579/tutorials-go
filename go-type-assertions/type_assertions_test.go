package go_type_assertions

import (
	"fmt"
)

func Example_TypeAssertion_Empty_Interface() {
	var i interface{} = "foo"

	var s = i.(string)
	fmt.Println(s)

	s, ok := i.(string)
	fmt.Println(s, ok)

	n, ok := i.(int) //panic 대신 ok:false, n는 zero 값으로 반환한다
	fmt.Println(n, ok)

	//Output:
	//foo
	//foo true
	//0 false
}

func Example_TypeAssertion_Empty_Interface_Illegal() {
	var i interface{} = "foo"

	var n = i.(int) //panic: interface conversion: interface {} is string, not int
	fmt.Println(n)

	//Output:
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

	s := p.(Student) //Person -> Student - student의 실제 값을 가져온다.
	fmt.Println(s.getName())
	fmt.Println(s.getPhone())

	//Output:
	//Frank
	//Frank
	//1111
}

func Example_TypeAssertion_다른_인터페이스로_값을_가져온다() {
	var p Person = Student{"Frank", 13, "1111"}

	ph := p.(Phone) //Person -> Phone
	fmt.Println(ph.getPhone())

	//Output:
	//1111
}

//타입 T가 인터페이스를 구현하고 있지 않기 때문에 컴파일 에러가 발생한다
func Example_TypeAssertion_인터페이스가_타입_T의_동적_값을_소유하지_않을_경우_컴파일_에러가_발생한다() {
	//var p Person = Student{"Frank", 13, "1111"}
	//value := p.(string) //impossible type assertion: string does not implement person (missing getName method)
	//fmt.Printf("%v, %T\n", value, value)

	//Output:
}

func Example_TypeAssertion_인터페이스가_타입_T의_실제_값을_가지고_있지_않는_경우_panic이_발생한다() {
	var p Person = nil
	//value := p.(Student) //panic: interface conversion: go_type_assertions.Person is nil, not go_type_assertions.Student
	value, ok := p.(Student)
	fmt.Printf("(%v, %T), ok: %v\n", value, value, ok)

	//Output:
	//({ 0 }, go_type_assertions.Student), ok: false
}

func Example_TypeAssertion_다른_인터페이스가_타입_T를_구현하지_않고_있으면_panic이_발생한다() {
	var p Person = Student{"Frank", 13, "1111"}
	//value := p.(Animal) //panic: interface conversion: go_type_assertions.Student is not go_type_assertions.Animal: missing method walk
	value, ok := p.(Animal)
	fmt.Printf("(%v, %T) %v\n", value, value, ok)

	//Output:
	//(<nil>, <nil>) false
}
