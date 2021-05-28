package go_ternary

import "fmt"

//go에서는 ternary operator가 존재하지 않음
func Example_Ternary_Operator_GetMax() {
	//return val > val2 ? val2 : val1
	getMax := func(val1, val2 int) int {
		if val1 > val2 {
			return val1
		} else {
			return val2
		}
	}
	fmt.Println(getMax(1, 2))
	//Output: 2
}

func Example_Ternary_Operator_EqualCheck() {
	checkStatus := func(status string) bool {
		if status == "done" {
			return true
		}
		return false
	}("done")

	fmt.Println(checkStatus)
	//Output: true
}
