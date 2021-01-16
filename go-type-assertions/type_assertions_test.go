package go_type_assertions

import "fmt"

func Example_TypeAssertion() {
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

func Example_TypeAssertion_Illegal() {
	var i interface{} = "foo"

	var n = i.(int) //panic: interface conversion: interface {} is string, not int
	fmt.Println(n)

	//Output:
}

type geometry interface {
	area() float64
	perimeter() float64
}

type rect struct {
	width, height float64
}

func (r rect) area() float64 {
	return r.width * r.height
}

func Example() {
	r := rect{width: 3, height: 4}
	fmt.Println(r.area())

	//type assertion
	//r2 := r.()

	//Output:
}
