package go_strings

import (
	"fmt"
	"os"
)

/*
https://gobyexample.com/string-formatting
*/

type point struct {
	x, y int
}

func Example_printf_struct() {
	p := point{1, 2}

	fmt.Printf("%v\n", p)
	fmt.Printf("%+v\n", p)
	fmt.Printf("%#v\n", p)
	fmt.Printf("%T\n", p)
	fmt.Printf("%p\n", &p)

	//Output:
	//{1 2}
	//{x:1 y:2}
	//go_string_formatting.point{x:1, y:2}
	//go_string_formatting.point
	//0xc00011c220
}

func Example_printf1() {
	fmt.Printf("%t\n", true)
	fmt.Printf("%d\n", 123)
	fmt.Printf("%b\n", 14)
	fmt.Printf("%c\n", 33)
	fmt.Printf("%x\n", 456)

	//Output:
	//true
	//123
	//1110
	//!
	//1c8
}

func Example_printf_float() {
	fmt.Printf("%f\n", 78.9)
	fmt.Printf("%e\n", 123400000.0)
	fmt.Printf("%E\n", 123400000.0)
	fmt.Printf("|%6d|%6d|\n", 12, 345)
	fmt.Printf("|%6.2f|%6.2f|\n", 1.2, 3.45)
	fmt.Printf("|%-6.2f|%-6.2f|\n", 1.2, 3.45)

	//Output:
	//78.900000
	//1.234000e+08
	//1.234000E+08
	//|    12|   345|
	//|  1.20|  3.45|
	//|1.20  |3.45  |
}

func Example_printf_string() {
	fmt.Printf("%s\n", "\"string\"")
	fmt.Printf("%q\n", "\"string\"")
	fmt.Printf("%x\n", "hex this")
	fmt.Printf("|%6s|%6s|\n", "foo", "b")
	fmt.Printf("|%-6s|%-6s|\n", "foo", "b")

	s := fmt.Sprintf("a %s", "string")
	fmt.Println(s)
	fmt.Fprintf(os.Stderr, "an %s\n", "error")

	//Output:
	//"string"
	//"\"string\""
	//6865782074686973
	//|   foo|     b|
	//|foo   |b     |
	//a string
}
