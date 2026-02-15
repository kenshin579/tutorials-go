package go_generics

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

// ============================================================
// 1. interface 기반 구현 - 런타임 다형성
// ============================================================

// Adder 인터페이스 - 메서드 기반 다형성
type Adder interface {
	Add(other Adder) Adder
	Value() string
}

// IntVal - Adder 구현
type IntVal struct {
	val int
}

func (i IntVal) Add(other Adder) Adder {
	o := other.(IntVal) // 타입 단언 필요 (런타임 체크)
	return IntVal{val: i.val + o.val}
}

func (i IntVal) Value() string {
	return fmt.Sprintf("%d", i.val)
}

// FloatVal - Adder 구현
type FloatVal struct {
	val float64
}

func (f FloatVal) Add(other Adder) Adder {
	o := other.(FloatVal) // 타입 단언 필요
	return FloatVal{val: f.val + o.val}
}

func (f FloatVal) Value() string {
	return fmt.Sprintf("%.1f", f.val)
}

// interface 기반 합산 함수
func sumInterface(items []Adder) Adder {
	result := items[0]
	for _, item := range items[1:] {
		result = result.Add(item)
	}
	return result
}

func Example_interfaceApproach() {
	ints := []Adder{IntVal{1}, IntVal{2}, IntVal{3}}
	result := sumInterface(ints)
	fmt.Println(result.Value())

	floats := []Adder{FloatVal{1.1}, FloatVal{2.2}, FloatVal{3.3}}
	result = sumInterface(floats)
	fmt.Println(result.Value())

	// Output:
	// 6
	// 6.6
}

// ============================================================
// 2. generics 기반 구현 - 컴파일 타임 다형성
// ============================================================

// generics 기반 합산 함수 - 타입 단언 불필요
func sumGeneric[T constraints.Integer | constraints.Float](items []T) T {
	var result T
	for _, item := range items {
		result += item
	}
	return result
}

func Example_genericApproach() {
	ints := []int{1, 2, 3}
	fmt.Println(sumGeneric(ints))

	floats := []float64{1.1, 2.2, 3.3}
	fmt.Printf("%.1f\n", sumGeneric(floats))

	// Output:
	// 6
	// 6.6
}

// ============================================================
// 3. interface가 적합한 경우 - 런타임에 타입이 결정
// ============================================================

// Shape 인터페이스 - 서로 다른 타입을 한 컬렉션에 담을 때
type Shape interface {
	Area() float64
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return 3.14159 * c.Radius * c.Radius
}

type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// 서로 다른 구체 타입을 한 슬라이스에 담을 수 있다
func totalArea(shapes []Shape) float64 {
	var total float64
	for _, s := range shapes {
		total += s.Area()
	}
	return total
}

func Example_interfaceMixedTypes() {
	shapes := []Shape{
		Circle{Radius: 5},
		Rectangle{Width: 3, Height: 4},
		Circle{Radius: 2},
	}
	fmt.Printf("%.2f\n", totalArea(shapes))

	// Output:
	// 103.11
}

// ============================================================
// 4. generics가 적합한 경우 - 타입 안전한 유틸리티
// ============================================================

// 동일 타입의 슬라이스에 대한 안전한 연산
func indexOf[T comparable](s []T, target T) int {
	for i, v := range s {
		if v == target {
			return i
		}
	}
	return -1
}

func Example_genericTypeSafe() {
	fmt.Println(indexOf([]int{10, 20, 30}, 20))
	fmt.Println(indexOf([]string{"a", "b", "c"}, "c"))
	fmt.Println(indexOf([]string{"a", "b", "c"}, "d"))

	// 컴파일 에러: 서로 다른 타입은 혼합 불가
	// indexOf([]int{1, 2, 3}, "hello") → 컴파일 에러

	// Output:
	// 1
	// 2
	// -1
}
