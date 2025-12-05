package go_shadowed

import (
	"fmt"
)

func Example_Shadowed() {
	if value := getValue("key1"); value == "value1" {
		fmt.Println("value", value)
		if value := getValue("key2"); value == "value2" {
			fmt.Println("value", value)
		}
		fmt.Println("value", value)
	}

	//Output:
	//value value1
	//value value2
	//value value1
}

func getValue(key string) string {
	switch key {
	case "key1":
		return "value1"
	case "key2":
		return "value2"
	default:
		return ""
	}
}
