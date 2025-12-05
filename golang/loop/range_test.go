package go_loop

import (
	"fmt"
)

func Example_range() {
	nodeEntities := []string{"a", "b", "c", "d", "e"}

	for i, fromNode := range nodeEntities[:len(nodeEntities)-1] {
		for j, toNode := range nodeEntities[i+1:] { //j의 값은 0부터 시작된다
			fmt.Println(i, j, fromNode, toNode)
		}
		fmt.Println()
	}

	//Output:
	//0 0 a b
	//0 1 a c
	//0 2 a d
	//0 3 a e
	//
	//1 0 b c
	//1 1 b d
	//1 2 b e
	//
	//2 0 c d
	//2 1 c e
	//
	//3 0 d e
}
