package go_generics

import (
	"fmt"
)

// 숫자 타입을 묶는 커스텀 constraint
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type FloatType interface {
	~float32 | ~float64
}

// constraint 합성 - 여러 constraint를 하나로 조합
type Number interface {
	Signed | Unsigned | FloatType
}

func sum[T Number](nums []T) T {
	var total T
	for _, n := range nums {
		total += n
	}
	return total
}

func Example_customConstraintNumber() {
	fmt.Println(sum([]int{1, 2, 3, 4, 5}))
	fmt.Println(sum([]float64{1.5, 2.5, 3.0}))
	fmt.Println(sum([]uint{10, 20, 30}))

	//Output:
	//15
	//7
	//60
}

// 메서드를 요구하는 constraint + 타입 제한을 결합하는 패턴
type Printable interface {
	~int | ~string
	Format() string
}

type Label string

func (l Label) Format() string {
	return fmt.Sprintf("[%s]", string(l))
}

type Code int

func (c Code) Format() string {
	return fmt.Sprintf("#%d", int(c))
}

func printFormatted[T Printable](items []T) {
	for _, item := range items {
		fmt.Println(item.Format())
	}
}

func Example_customConstraintWithMethod() {
	printFormatted([]Label{"hello", "world"})
	printFormatted([]Code{100, 200})

	//Output:
	//[hello]
	//[world]
	//#100
	//#200
}

// Ordered constraint 직접 구현 (표준 라이브러리의 cmp.Ordered와 유사)
type MyOrdered interface {
	Signed | Unsigned | FloatType | ~string
}

func maxVal[T MyOrdered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func Example_customOrdered() {
	fmt.Println(maxVal(10, 20))
	fmt.Println(maxVal(3.14, 2.71))
	fmt.Println(maxVal("apple", "banana"))

	//Output:
	//20
	//3.14
	//banana
}
