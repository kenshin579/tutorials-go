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

func Example_PrintCat() {
	var a MyInt = 10
	var b MyInt = 100
	PrintCat(a, b)

	//Output:

}

// 타입 제한자. 인터페이스로 사용 불가능
type Integer interface {
	~int8 | ~int16 | ~int32 | ~int64 | ~int
}

// 문법 오류 발생
func PrintMin1(a, b Integer) {
	if a < b {
		fmt.Println(a.String())
	} else {
		fmt.Println(b.String())
	}
}

// 타입 제한자+인터페이스 ==> 타입 제한자. 인터페이스로 사용 불가능
type Stringer interface {
	Integer
	ToString
}

func PrintMin2[T Stringer](a, b T) {
	if a < b {
		fmt.Println(a.String())
	} else {
		fmt.Println(b.String())
	}
}

// 문법 오류 발생
func PrintMin3(a, b Stringer) {
	if a < b {
		fmt.Println(a.String())
	} else {
		fmt.Println(b.String())
	}
}

type MyInt int

func (m MyInt) String() string {
	return fmt.Sprintf("%d", m)
}

func Example_Print() {
	var a MyInt = 10
	var b MyInt = 100
	PrintMin3(a, b)

	//Output:
}
