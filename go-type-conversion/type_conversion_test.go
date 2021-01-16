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

func Example_TypeConversion2() {
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
	bob := Parent{
		name: "bob",
		age:  15,
	}
	babyBob := Child(bob)
	fmt.Println(bob)
	fmt.Println(babyBob)

	//Output:
	//{bob 15}
	//{bob 15}
}

func Example_TypeConversion_Between_Structs_변환_안되는_경우() {
	bob := Parent{
		name: "bob",
		age:  15,
	}

	//babyBod := Pet(bob) //cannot convert bob (type Parent) to type Pet
	fmt.Println(bob)

	//Output:
	//{bob 15}
}
