package go_type_conversion

import "fmt"

func Example_TypeConversion() {
	var i = 52
	var j float64 = float64(i)
	var k = uint32(j)

	fmt.Println(i)
	fmt.Println(j)
	fmt.Println(k)

	//Output:
	//52
	//52
	//52
}
