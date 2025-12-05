package go_generics

import "fmt"

func foo1(a interface{}) interface{} {
	return a
}

func foo2[T any](a T) T {
	return a
}

func Example_Print() {
	var (
		a int = 10
		b int = 20
		c int
	)
	c = foo1(a).(int) // 리턴 타입이 interface{} 이다.
	fmt.Println(c)
	c = foo2(b) // 리턴 타입이 int이다.
	fmt.Println(c)

	//Output:
	//10
	//20
}
