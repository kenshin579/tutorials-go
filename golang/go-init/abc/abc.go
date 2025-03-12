package abc

import "fmt"

var version string

func init() {
	fmt.Println("abc init")
	version = "1.1"
}

func Version() string {
	fmt.Println("abc getVersion")
	return version
}
