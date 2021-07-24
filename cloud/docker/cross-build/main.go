package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Printf("Hello from %s architecture\n", runtime.GOARCH)
}
