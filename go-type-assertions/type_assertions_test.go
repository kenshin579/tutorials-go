package go_type_assertions

import "fmt"

func Example_TypeAssertions() {
	var i interface{} = "foo"

	var s = i.(string)
	fmt.Println(s)

	s, ok := i.(string)
	fmt.Println(s, ok)

	n, ok := i.(int) //panic 대신 ok:false, t는 zero 값으로 반환한다
	fmt.Println(n, ok)

	//Output:
	//foo
	//foo true
	//0 false
}

func Example_TypeAssertions_Illegal() {
	var i interface{} = "foo"

	var n = i.(int) //panic: interface conversion: interface {} is string, not int
	fmt.Println(n)

	//Output:
}
