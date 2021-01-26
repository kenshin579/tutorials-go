package go_methods

import "fmt"

type Point struct {
	x, y float64
}

func (v *Point) ScalePointerReceiver(f float64) {
	v.x = v.x * f
	v.y = v.y * f
}

func (v Point) ScaleValueReceiver(f float64) {
	fmt.Printf("v.x:%.2f, v.y:%.2f\n", v.x*f, v.y*f)
}

func ScaleFuncPointerParameter(v *Point, f float64) {
	v.x = v.x * f
	v.y = v.y * f
}

func ScaleFuncValueParameter(v Point, f float64) {
	fmt.Printf("v.x:%.2f, v.y:%.2f\n", v.x*f, v.y*f)
}

//indirection 의미 : dereferencing이라는 의미로 자체 값 대신에 다른 방법(ex. 메모리 주소)으로 참조할 수 있는 방법
func Example_Func_Pointer_Parameter_Indirection() {
	p := Point{3, 4}
	fmt.Printf("value:%+v type:%T\n", p, p)
	//ScaleFuncPointerParameter(p, 10) //컴파일 오류 발생함
	ScaleFuncPointerParameter(&p, 10) //포인터 인자는 주소 값만 받음
	fmt.Println(p)

	//Output:
	//value:{x:3 y:4} type:go_methods.Point
	//{30 40}
}

func Example_Func_Value_Parameter_Indirection() {
	p := Point{3, 4}
	fmt.Printf("value:%+v type:%T\n", p, p)
	//ScaleFuncPointerParameter(&p, 10) //컴파일 오류 - value 인자는 주소 값으로 넘겨줄 수 없음
	ScaleFuncValueParameter(p, 10)
	fmt.Println(p)

	//Output:
	//value:{x:3 y:4} type:go_methods.Point
	//v.x:30.00, v.y:40.00
	//{3 4}
}

func Example_Method_Pointer_Receiver_Indirection() {
	p1 := Point{1, 2}
	fmt.Printf("value:%+v type:%T\n", p1, p1)
	p1.ScalePointerReceiver(3) //p는 value이지만, Go는 자동으로 (&p1).ScalePointerReceiver(3)으로 해석해서 실행함
	fmt.Println(p1)

	p2 := &Point{4, 3}
	fmt.Printf("value:%+v type:%T\n", p2, p2)
	p2.ScalePointerReceiver(3) //p2 포인터이여도 잘 실행됨
	fmt.Println(*p2)

	//Output:
	//value:{x:1 y:2} type:go_methods.Point
	//{3 6}
	//{12 9}
}

func Example_Method_Value_Receiver_Indirection() {
	p := &Point{4, 3}
	fmt.Printf("value:%+v type:%T\n", p, p)
	p.ScaleValueReceiver(3) //p는 pointer이지만, Go는 자동으로 (*p).ScaleValueReceiver(3)으로 해석해서 실행함
	fmt.Println(*p)

	//Output:
	//value:&{x:4 y:3} type:*go_methods.Point
	//v.x:12.00, v.y:9.00
	//&{4 3}
}
