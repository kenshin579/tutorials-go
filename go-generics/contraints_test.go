package go_generics

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

//type Ordered interface {
//	Integer | Float | ~string
//}

// constraints.Ordered 은 크기 비교가 가능한 타입 제한자입니다.
// 새로 추가된 comparable 키워드는 == 또는 != 연산이 가능한 타입 제한자입니다.
func min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func Example_Ordered() {
	var (
		a int     = 10
		b int     = 20
		c int16   = 10
		d int16   = 20
		e float32 = 3.14
		f float32 = 1.14
	)
	fmt.Println(min(a, b))
	fmt.Println(min(c, d))
	fmt.Println(min(e, f))
	var (
		h = "Hello"
		i = "World"
	)
	fmt.Println(min(h, i))

	//Output:
	//10
	//10
	//1.14
	//Hello
}

// int 앞에 ~(틸트)가 붙어 있습니다. 이 토큰도 Go 1.18에서 새로 추가된 문법으로 "확장된" 이란 의미를 가집니다.
type Integer interface {
	//int | int8 | int16 | int32 | int64
	~int | int8 | int16 | int32 | int64
}

// 변수 c, d는 MyInt를 사용하고 Integer 인터페이스에는 ~int 가 있어서 int를 확장한 MyInt를 사용할 수 있다.
type MyInt int

func min2[T Integer](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func Example_Extended() {
	var (
		a int = 10
		b int = 20
	)
	fmt.Println(min2(a, b))
	var (
		c MyInt = 10
		d MyInt = 20
	)
	//possibly missing ~ for int in constraint Integer
	//type Integer 안에 'int' 앞에 틸트(~)를 붙여워야 에러가 안난다. '~int'
	fmt.Println(min2(c, d))

	//Output:
	//10
	//10
}
