package main

import (
	"fmt"
)

func Example_iota_기본() {
	const (
		c0 = iota
		c1 = iota
		c2 = iota
	)

	fmt.Println(c0, c1, c2)
	//Output:
	//0 1 2
}

func Example_iota_기본2() {
	const (
		c0 = iota
		c1
		c2
	)

	fmt.Println(c0, c1, c2)
	//Output:
	//0 1 2
}

func Example_iota_시작값_변경() {
	const (
		c0 = iota + 5
		c1
		c2
	)

	fmt.Println(c0, c1, c2)
	//Output:
	//5 6 7
}

func Example_iota_중간값_skipped() {
	const (
		c0 = iota + 1
		_
		c1
		c2
	)

	fmt.Println(c0, c1, c2)
	//Output:
	//1 3 4
}

/*
https://stackoverflow.com/questions/57053373/how-to-skip-a-lot-of-values-when-define-const-variable-with-iota/57053431#57053431
*/
func Example_iota_중간값_다르게_지정() {
	//수동으로 계산하기
	const (
		c0 = iota + 1   //1
		c1              //2
		c2 = iota + 98  //100
		c3              //101
		c4 = iota + 496 //500
		c5              //501
	)

	fmt.Println(c0, c1, c2, c3, c4, c5)
	//Output:
	//1 2 100 101 500 501
}

func Example_iota_중간값_다르게_지정2() {
	//수동으로 계산하기
	const (
		c0 = iota + 1 //1
		c1            //2
	)
	const (
		c2 = iota + 100 //100
		c3              //101
	)
	const (
		c4 = iota + 500 //500
		c5              //501
	)

	fmt.Println(c0, c1, c2, c3, c4, c5)
	//Output:
	//1 2 100 101 500 501
}

type Direction int

func Example_iota_Direction_예제() {
	const (
		North Direction = iota
		East
		South
		West
	)

	var direction Direction = North
	fmt.Print(direction)

	switch direction {
	case North:
		fmt.Println(" goes up.")
	case South:
		fmt.Println(" goes down.")
	default:
		fmt.Println(" stays put.")
	}

	//Output:
	//North goes up.
}

func (d Direction) String() string {
	return [...]string{"North", "East", "South", "West"}[d]
}

func Example_iota_BYTE_예제() {
	type ByteSize float64

	const (
		_           = iota             // ignore first value by assigning to blank identifier
		KB ByteSize = 1 << (10 * iota) //2^(10*1) = 1024
		MB                             //2^(10*2) = 1,048,576
		GB
		TB
		PB
		EB
		ZB
		YB
	)
	var fileSize ByteSize = 4000000000 //4 GB
	fmt.Printf("%.2f GB", fileSize/GB)
	//Output:
	//3.73 GB
}

type WeekDay int

func Example_iterate_weekdays() {
	const (
		Sunday WeekDay = iota
		Monday
		Tuesday
		Wednesday
		Thursday
		Friday
		Saturday
		numberOfDays
	)

	for day := WeekDay(0); day < numberOfDays; day++ {
		fmt.Print(" ", day)
	}
	fmt.Println("")
	//Output:
	//Sunday Monday Tuesday Wednesday Thursday Friday Saturday
}

func (d WeekDay) String() string {
	return [...]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}[d]
}
