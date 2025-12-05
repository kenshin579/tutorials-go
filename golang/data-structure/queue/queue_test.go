package queue

import (
	"container/list"
	"fmt"
)

// https://tutorialedge.net/golang/go-linked-lists-tutorial/
// https://golangbyexample.com/queue-in-golang/
func ExampleContainer_List() {
	list := list.New()

	//push
	list.PushBack("A")
	list.PushBack(100)
	list.PushBack(true)
	list.PushFront("B")

	//Pop
	front := list.Front()
	fmt.Println(front.Value)
	fmt.Println(list.Len())

	list.Remove(front)
	fmt.Println(list.Len())

	for element := list.Front(); element != nil; element = element.Next() {
		// do something with element.Value
		fmt.Println(element.Value)
	}

	//Output:
	//B
	//4
	//3
	//A
	//100
	//true
}

// https://gist.github.com/matishsiao/24050743170157072a4d219950237aab
func ExampleListUsingSlices() {
	var queue []string

	//push
	queue = append(queue, "A")
	queue = append(queue, "C")
	queue = append(queue, "B")
	queue = append(queue, "D")

	fmt.Println(len(queue))

	fmt.Println(queue[0]) //peek
	queue = queue[1:]     //dequeue
	fmt.Println(len(queue))

	for _, item := range queue {
		fmt.Println(item)
	}

	//Output:
	//4
	//A
	//3
	//C
	//B
	//D

}
