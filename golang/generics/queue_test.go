package go_generics

import (
	"fmt"
)

// Generic Queue 구현 (FIFO)
type Queue[T any] struct {
	items []T
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{items: make([]T, 0)}
}

func (q *Queue[T]) Enqueue(item T) {
	q.items = append(q.items, item)
}

func (q *Queue[T]) Dequeue() (T, bool) {
	var zero T
	if len(q.items) == 0 {
		return zero, false
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}

func (q *Queue[T]) Front() (T, bool) {
	var zero T
	if len(q.items) == 0 {
		return zero, false
	}
	return q.items[0], true
}

func (q *Queue[T]) Size() int {
	return len(q.items)
}

func (q *Queue[T]) IsEmpty() bool {
	return len(q.items) == 0
}

func Example_queueInt() {
	q := NewQueue[int]()
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)

	fmt.Println(q.Size())
	fmt.Println(q.Front())

	val, _ := q.Dequeue()
	fmt.Println(val)

	val, _ = q.Dequeue()
	fmt.Println(val)

	fmt.Println(q.Size())

	//Output:
	//3
	//1 true
	//1
	//2
	//1
}

func Example_queueString() {
	q := NewQueue[string]()
	q.Enqueue("first")
	q.Enqueue("second")
	q.Enqueue("third")

	for !q.IsEmpty() {
		val, _ := q.Dequeue()
		fmt.Println(val)
	}

	//Output:
	//first
	//second
	//third
}
