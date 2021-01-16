package go_type_switch

import "fmt"

func typeSwitchTest(i interface{}) {
	switch v := i.(type) {
	case nil:
		fmt.Println("x is nil")
	case int:
		fmt.Println("x is", v)
	case bool, string:
		fmt.Println("x is bool or string")
	default:
		fmt.Printf("type unknown %T\n", v)
	}
}

func Example_TypeSwitch() {
	typeSwitchTest("value")
	typeSwitchTest(23)
	typeSwitchTest(true)
	typeSwitchTest(nil)
	typeSwitchTest([]int{})

	//Output:
	//x is bool or string
	//x is 23
	//x is bool or string
	//x is nil
	//type unknown []int
}
