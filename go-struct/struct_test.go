package go_struct

import (
	"fmt"
	"unsafe"
)

//todo: https://dave.cheney.net/2014/03/25/the-empty-struct
func Example_Struct_Width() {
	var s string
	var c complex128
	fmt.Println(unsafe.Sizeof(s))
	fmt.Println(unsafe.Sizeof(c))

	//Output:
	//16
	//16
}

func Example_Empty_Struct() {
	var emptyStruct struct{} //zero bytes of storage
	fmt.Println(unsafe.Sizeof(emptyStruct))

	type emptyNestedStruct struct { //이것도 동일하게 zero bytes이다
		A struct{}
		B struct{}
	}
	var ns emptyNestedStruct
	fmt.Println(unsafe.Sizeof(ns))

	//Output:
	//0
	//0
}

func Example_Empty_Struct_Instance() {
	// not the zero value, a real new struct{} instance
	a := struct{}{} //zero 값이 아니라 샐제 struct{} 인스턴스이다
	b := struct{}{}
	fmt.Println(a == b)
	//Output: true
}
