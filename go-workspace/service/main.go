package main

import (
	"fmt"

	"github.com/kenshin579/tutorials-go/go-workspace/adder"
)

func main() {
	sum := adder.Add(1, 2)
	fmt.Println(sum)
}
