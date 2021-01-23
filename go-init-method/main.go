package main

//import된 순서대로 init()가 호출된다
import (
	"fmt"
	"github.com/kenshin579/tutorials-go/go-init-method/abc"
	"github.com/kenshin579/tutorials-go/go-init-method/lib1"
)

var version string

func init() {
	fmt.Println("main init")
	version = "1"
}

func main() {
	fmt.Println("main called")
	fmt.Println(version)
	fmt.Println(lib1.Version())
	fmt.Println(abc.Version())
}
