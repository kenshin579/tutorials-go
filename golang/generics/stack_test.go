package go_generics

import (
	"fmt"
)

// Generic Stack 구현
type Stack[T any] struct {
	items []T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{items: make([]T, 0)}
}

func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, bool) {
	var zero T
	if len(s.items) == 0 {
		return zero, false
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item, true
}

func (s *Stack[T]) Peek() (T, bool) {
	var zero T
	if len(s.items) == 0 {
		return zero, false
	}
	return s.items[len(s.items)-1], true
}

func (s *Stack[T]) Size() int {
	return len(s.items)
}

func (s *Stack[T]) IsEmpty() bool {
	return len(s.items) == 0
}

func Example_stackInt() {
	s := NewStack[int]()
	s.Push(1)
	s.Push(2)
	s.Push(3)

	fmt.Println(s.Size())
	fmt.Println(s.Peek())

	val, _ := s.Pop()
	fmt.Println(val)
	fmt.Println(s.Size())

	//Output:
	//3
	//3 true
	//3
	//2
}

func Example_stackString() {
	s := NewStack[string]()
	s.Push("hello")
	s.Push("world")

	val, ok := s.Pop()
	fmt.Println(val, ok)

	val, ok = s.Pop()
	fmt.Println(val, ok)

	val, ok = s.Pop()
	fmt.Printf("%q %v\n", val, ok)

	//Output:
	//world true
	//hello true
	//"" false
}
