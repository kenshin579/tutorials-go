package go1_26_test

import (
	"fmt"
	"testing"
)

// Ordered - 제네릭 타입이 자기 자신을 타입 파라미터로 참조
type Ordered[T Ordered[T]] interface {
	Less(T) bool
}

// MyInt - Ordered 인터페이스 구현
type MyInt int

func (a MyInt) Less(b MyInt) bool {
	return a < b
}

// MyString - Ordered 인터페이스 구현
type MyString string

func (a MyString) Less(b MyString) bool {
	return a < b
}

// Min - 제네릭 자기참조를 활용한 최솟값 함수
func Min[T Ordered[T]](a, b T) T {
	if a.Less(b) {
		return a
	}
	return b
}

// Max - 제네릭 자기참조를 활용한 최댓값 함수
func Max[T Ordered[T]](a, b T) T {
	if b.Less(a) {
		return a
	}
	return b
}

func TestRecursiveGenericMin(t *testing.T) {
	// MyInt 타입으로 Min 사용
	result := Min(MyInt(3), MyInt(7))
	fmt.Printf("Min(3, 7) = %d\n", result)
	if result != 3 {
		t.Errorf("expected 3, got %d", result)
	}
}

func TestRecursiveGenericMax(t *testing.T) {
	result := Max(MyInt(3), MyInt(7))
	fmt.Printf("Max(3, 7) = %d\n", result)
	if result != 7 {
		t.Errorf("expected 7, got %d", result)
	}
}

func TestRecursiveGenericString(t *testing.T) {
	result := Min(MyString("apple"), MyString("banana"))
	fmt.Printf("Min(apple, banana) = %s\n", result)
	if result != "apple" {
		t.Errorf("expected 'apple', got '%s'", result)
	}
}

// Adder - F-bounded 다형성 패턴
type Adder[A Adder[A]] interface {
	Add(A) A
}

// Vector2D - Adder 인터페이스 구현
type Vector2D struct {
	X, Y float64
}

func (v Vector2D) Add(other Vector2D) Vector2D {
	return Vector2D{X: v.X + other.X, Y: v.Y + other.Y}
}

// Sum - Adder 제약을 활용한 합산 함수
func Sum[A Adder[A]](values []A) A {
	var zero A
	result := zero
	for _, v := range values {
		result = result.Add(v)
	}
	return result
}

func TestAdderInterface(t *testing.T) {
	vectors := []Vector2D{
		{1, 2},
		{3, 4},
		{5, 6},
	}
	result := Sum(vectors)
	fmt.Printf("Sum = %+v\n", result)

	if result.X != 9 || result.Y != 12 {
		t.Errorf("expected {9, 12}, got %+v", result)
	}
}
