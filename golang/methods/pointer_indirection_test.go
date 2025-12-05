package go_methods

import "fmt"

type Rectangle struct {
	height, width int
}

func (r *Rectangle) area() {
	fmt.Println(r.height * r.width)
}

func area(r *Rectangle) {
	fmt.Println(r.height * r.width)
}

func (r Rectangle) perimeter() {
	fmt.Println(2 * (r.height * r.width))
}

func perimeter(r Rectangle) {
	fmt.Println(2 * (r.height * r.width))
}

func Example_Indirection_Func_Pointer_Parameter() {
	r := Rectangle{
		height: 10,
		width:  3,
	}

	//area(r) //컴파일 오류 - 함수는 포인터 인자만 받을 수 있음
	area(&r)

	//Output:
	//30

}

func Example_Indirection_Method_Pointer_Receiver() {
	r := Rectangle{
		height: 10,
		width:  3,
	}

	r.area()
	(&r).area()

	//Output:
	//30
	//30

}

func Example_Indirection_Func_Value_Parameter() {
	r := Rectangle{
		height: 10,
		width:  3,
	}

	//perimeter(&r) //컴파일 오류 - 함수는 value 인자만 받을 수 있음
	perimeter(r)

	//Output:
	//60
}

func Example_Indirection_Method_Value_Receiver() {
	r := Rectangle{
		height: 10,
		width:  3,
	}

	r.perimeter()
	(&r).perimeter()

	//Output:
	//60
	//60
}
