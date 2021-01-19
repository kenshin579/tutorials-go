package main

import (
	"fmt"
	"github.com/kenshin579/tutorials-go/go-init-method/lib"
)

var version string

func init() {
	fmt.Println("main init")
	version = "1"
}

func main() {
	fmt.Println("main called")
	fmt.Println(version)
	fmt.Println(lib.Version())
}
