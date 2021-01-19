package lib

import "fmt"

var version string

func init() {
	fmt.Println("lib init")
	version = "1.1"
}

func Version() string {
	fmt.Println("getVersion")
	return version
}
