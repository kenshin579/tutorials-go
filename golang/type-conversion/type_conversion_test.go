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

type myInt int

func Example_TypeConversion_custmonType() {
	var i myInt = 4
	originalInt := int(i)
	i = myInt(originalInt)

	fmt.Println(originalInt)
	fmt.Println(i)

	//Output:
	//4
	//4
}

func Example_TypeConversion_Byte_And_String() {
	greeting := []byte("hello")
	greetingStr := string(greeting)

	fmt.Println(greetingStr)

	//Output: hello
}

type Parent struct {
	name string
	age  int
}

type Child struct {
	name string
	age  int
}

type Pet struct {
	name string
}

func Example_TypeConversion_Between_Structs() {
	parent := Parent{
		name: "parent",
		age:  15,
	}
	child := Child(parent)
	fmt.Println(child)
	parent = Parent(child)
	fmt.Println(parent)

	//Output:
	//{parent 15}
	//{parent 15}
}

func Example_TypeConversion_Between_Structs_변환_안되는_경우() {
	parent := Parent{
		name: "parent",
		age:  15,
	}

	//babyBod := Pet(parent) //cannot convert parent (type Parent) to type Pet
	fmt.Println(parent)

	//Output:
	//{parent 15}
}
