package go_generics

import (
	"fmt"
)

// 인터페이스. 타입 제한자로 사용 가능
type ToString interface {
	String() string
}

func PrintCat[T ToString](a, b T) {
	fmt.Printf("%s-%s", a.String(), b.String())
}

func Example_printCat() {
	var a MyIntType = 10
	var b MyIntType = 100
	PrintCat(a, b)

	//Output:
	//10-100
}

// 타입 제한자. 인터페이스로 사용 불가능
type IntegerTilde interface {
	~int8 | ~int16 | ~int32 | ~int64 | ~int
}

// 아래 코드는 타입 제한자를 일반 인터페이스처럼 사용하면 문법 오류가 발생하는 예시이다.
// func PrintMin1(a, b IntegerTilde) { ... }
// → cannot use type IntegerTilde outside a type constraint

// 타입 제한자+인터페이스 ==> 타입 제한자. 인터페이스로 사용 불가능
type Stringer interface {
	IntegerTilde
	ToString
}

func PrintMin2[T Stringer](a, b T) {
	if a < b {
		fmt.Println(a.String())
	} else {
		fmt.Println(b.String())
	}
}

// 아래 코드는 타입 제한자가 포함된 인터페이스를 일반 인터페이스처럼 사용하면 문법 오류가 발생하는 예시이다.
// func PrintMin3(a, b Stringer) { ... }
// → cannot use type Stringer outside a type constraint

type MyIntType int

func (m MyIntType) String() string {
	return fmt.Sprintf("%d", m)
}

func Example_printMin() {
	var a MyIntType = 10
	var b MyIntType = 100
	PrintMin2(a, b)

	//Output:
	//10
}
