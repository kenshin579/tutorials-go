package status_code

import (
	"fmt"
)

func Example() {
	fmt.Println(Code.Message(Unprocessed))
	fmt.Println(Unprocessed.Message())
	fmt.Println(Unprocessed.String())
	json, _ := Unprocessed.MarshalJSON()
	fmt.Println(string(json))

	//Output:
	//UNPROCESSED
	//UNPROCESSED
	//string: 1001
	//"1001"

}
