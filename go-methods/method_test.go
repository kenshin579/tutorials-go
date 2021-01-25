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
	hyundaiCar1 := Car{"현대", "빨강", 10000, 0}
	fmt.Println("hyundaiCar1", hyundaiCar1)

	hyundaiCar1.SpeedUp(10)
	fmt.Println("hyundaiCar1", hyundaiCar1)

	//https://golangbot.com/methods/

	//Output:
	//hyundaiCar1 {현대 빨강 10000 0}
	//hyundaiCar1 {현대 빨강 10000 10}
}

type myFloat float64

func (v myFloat) ceil() float64 {
	return math.Ceil(float64(v))
}

func Example_Method_Non_Struct_Type() {
	v := myFloat(1.3)
	fmt.Println(v)

	//Output:
	//1.3
}
