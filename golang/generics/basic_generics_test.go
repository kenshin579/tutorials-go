package go_generics

import (
	"fmt"
)

func minNoGeneric(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func minInt16(a, b int16) int16 {
	if a < b {
		return a
	}
	return b
}

func Example_noGenerics() {
	var (
		a int = 10
		b int = 20
	)
	fmt.Println(minNoGeneric(a, b))
	var (
		c int16 = 10
		d int16 = 20
	)
	fmt.Println(minInt16(c, d))

	//Output:
	//10
	//10
}

// func name['식별자' '타입 제한자']
func printAny[T any](a T) {
	fmt.Println(a)
}

func Example_generics() {
	var (
		a int     = 10
		b float32 = 3.14
		c string  = "hello"
	)
	printAny(a)
	printAny(b)
	printAny(c)
}

//func min[T any](a, b T) T {
//	if a < b { // 문법 오류가 발생. any는 < 연산을 지원하지 않는다.
//		return a
//	}
//	return b
//}

// 타입 제안자에 언떤 타입이 들어갈지 범위를 정함
func minType[T int | int16 | int32 | int64 | float32 | float64](a, b T) T {
	if a < b { // 위 타입들이 < 연산자를 지원하기 때문에 문법 오류가 없다.
		return a
	}
	return b
}

// 타입 제한자는 파이프 연산자로 여러 개를 쉽게 추가가 가능하다.
func Example_typeConstraint() {
	var (
		a int     = 10
		b int     = 20
		c int16   = 10
		d int16   = 20
		e float32 = 3.14
		f float32 = 1.14
	)
	fmt.Println(minType(a, b))
	fmt.Println(minType(c, d))
	fmt.Println(minType(e, f))

	//Output:
	//10
	//10
	//1.14
}

// 타입 제한자 선언
// 매번 타입 제한자를 만드는 것은 비효율적이므로 타입 제한자를 interface 키워드로 선언하여 사용이 가능하다.
type ComparableNumbers interface {
	int | int16 | int32 | int64 | float32 | float64
}

type IntegerType interface {
	int | int16 | int32 | int64
}

type Float interface {
	float32 | float64
}

type ComparableNumbers2 interface {
	IntegerType | Float
}

func minComparableNumbers[T ComparableNumbers](a, b T) T {
	if a < b { // 위 타입들이 < 연산자를 지원하기 때문에 문법 오류가 없다.
		return a
	}
	return b
}

func minComparableNumbers2[T ComparableNumbers](a, b T) T {
	if a < b { // 위 타입들이 < 연산자를 지원하기 때문에 문법 오류가 없다.
		return a
	}
	return b
}

func Example_comparableNumbers() {
	var (
		a int     = 10
		b int     = 20
		c int16   = 10
		d int16   = 20
		e float32 = 3.14
		f float32 = 1.14
	)
	fmt.Println(minComparableNumbers(a, b))
	fmt.Println(minComparableNumbers(c, d))
	fmt.Println(minComparableNumbers(e, f))

	fmt.Println(minComparableNumbers2(a, b))
	fmt.Println(minComparableNumbers2(c, d))
	fmt.Println(minComparableNumbers2(e, f))

	//Output:
	//10
	//10
	//1.14
	//10
	//10
	//1.14
}
