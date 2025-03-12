package go_methods

import (
	"fmt"
	"math"
)

type Car struct {
	brand   string
	color   string
	mileage int
	speed   int
}

func (c Car) Color() string {
	return c.color
}

func (c *Car) SpeedUp(s int) {
	c.speed += s
}

func Example_Method_Value_Receiver() {
	hyundaiCar := Car{"현대", "빨강", 10000, 0}
	//fmt.Println("hyundaiCar", hyundaiCar)

	fmt.Println(hyundaiCar.Color())

	//Output:
	//빨강
}

func Example_Method_Pointer_Receiver() {
	hyundaiCar := Car{"현대", "빨강", 10000, 0}
	fmt.Println("hyundaiCar", hyundaiCar)

	hyundaiCar.SpeedUp(10)
	fmt.Println("hyundaiCar", hyundaiCar)

	//https://golangbot.com/methods/

	//Output:
	//hyundaiCar {현대 빨강 10000 0}
	//hyundaiCar {현대 빨강 10000 10}
}

//int 타입과 ceil 메서드는 같은 패키지 레벨이 존재하지 않기 떄문에 컴파일 오류가 발생한다
//func (f float64) ceil() float64 {
//	return math.Ceil(float64(i))
//}

type myFloat float64

func (m myFloat) ceil() float64 {
	return math.Ceil(float64(m))
}

func Example_Method_Non_Struct_Type() {
	v := myFloat(1.3)
	fmt.Println(v)

	//Output:
	//1.3
}
