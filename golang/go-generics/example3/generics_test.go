package example3

import (
	"fmt"
	"strings"
)

// Go 제네릭은 함수와 구조체에서 사용이 가능하다. 아직 메소드에서 사용은 불가능하다.
type Node[T any] struct {
	val  T // struct 의 value 타입을 T로 사용한다.
	next *Node[T]
}

func NewNode[T any](v T) *Node[T] { // 새로운 Node를 만들 때도 제네릭이 필요하다.
	return &Node[T]{val: v}
}

/*
Node의 메소드인 Push에 제네릭 T가 포함된다. 하지만 이곳에서 새로운 다른 제네릭을 선언하거나 사용하는 것은 문법 오류이다.
문법오류: func (n *Node[T]) Push[F any](f F) * Node[T]
*/
func (n *Node[T]) Push(v T) *Node[T] {
	node := NewNode(v)
	n.next = node
	return node
}

func Example_Generics_Func() {
	node1 := NewNode(1) // *Node[int]
	node1.Push(2).Push(3).Push(4)

	for node1 != nil {
		fmt.Println(node1.val)
		node1 = node1.next
	}

	node2 := NewNode("hello") // *Node[string]
	node2.Push("how").Push("are").Push("you").Push("?")

	for node2 != nil {
		fmt.Println(node2.val)
		node2 = node2.next
	}

	//Output:
	//1
	//2
	//3
	//4
	//hello
	//how
	//are
	//you
	//?
}

// Map 함수는 값을 변환시켜 배열로 반환하는 함수이다.
func Map[F, T any](s []F, f func(F) T) []T {
	rst := make([]T, len(s))
	for i, v := range s {
		rst[i] = f(v)
	}
	return rst
}

func Example_Generic_Map() {
	doubled := Map([]int{1, 2, 3}, func(i int) int {
		return i * 2
	})
	fmt.Println(doubled)

	uppered := Map([]string{"Hello", "world", "abc"}, func(s string) string {
		return strings.ToUpper(s)
	})
	fmt.Println(uppered)

	toString := Map([]int{1, 2, 3}, func(i int) string {
		return fmt.Sprintf("str%d", i)
	})
	fmt.Println(toString)

	//Output:
	//[2 4 6]
	//[HELLO WORLD ABC]
	//[str1 str2 str3]

}
