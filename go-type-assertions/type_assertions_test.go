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

type Shape interface {
	Area() float64
}

type Object interface {
	Volume() float64
}

type Cube struct {
	side float64
}

func (c Cube) Area() float64 {
	return 6 * (c.side * c.side)
}

func (c Cube) Volume() float64 {
	return c.side * c.side * c.side
}

type Skin interface {
	Color() float64
}

func Example_TypeAssertion_Interface_With_Method() {
	var s Shape = Cube{3}
	c := s.(Cube) //Shape -> Cube
	fmt.Println(c.Area())
	fmt.Println(c.Volume())

	//Output:
	//54
	//27
}

func Example() {
	var s Shape = Cube{3}
	value1, ok1 := s.(Object) //Shape -> Object
	fmt.Printf("shape 값(%v) Object 인터페이스를 구현되어 있나? %v\n", value1, ok1)

	value2, ok2 := s.(Skin)
	fmt.Printf("shape 값(%v) Skin 인테퍼이스를 구현되어 있나?? %v\n", value2, ok2)

	//Output:
	//shape 값({3}) Object 인터페이스를 구현되어 있나? true
	//shape 값(<nil>) Skin 인테퍼이스를 구현되어 있나?? false
}
