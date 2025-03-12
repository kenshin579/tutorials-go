package lib1

import "fmt"

var version string

func init() {
	fmt.Println("lib1 init")
	version = "1.1"
}

func Version() string {
	fmt.Println("lib1 getVersion")
	return version
}
