package config

import "fmt"

var Version = "Unknown"
var BuildTime = "Unknown"

func PrintBuildInfo() {
	fmt.Println("build.Version:\t", Version)
	fmt.Println("build.Time:\t", BuildTime)
}
