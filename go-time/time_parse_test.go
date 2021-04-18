package main

import (
	"fmt"
	"time"
)

const (
	LayoutISO    = "2006-01-02"
	LayoutCustom = "2006-01-02 15:04:05"
)

func ExampleTimeParse_Date_String() {
	dateStr := "2020-04-20"
	parseDate, _ := time.Parse(LayoutISO, dateStr)
	fmt.Println(parseDate)
	fmt.Println(parseDate.Format(time.RFC822))
	fmt.Println(parseDate.Format(time.RFC3339))

	//Output:
	//2020-04-20 00:00:00 +0000 UTC
	//20 Apr 20 00:00 UTC
	//2020-04-20T00:00:00Z
}

func ExampleTimeParse_DateTime_String() {
	dateStr := "2020-04-20 12:33:30"
	parseDate, _ := time.Parse(LayoutCustom, dateStr)
	fmt.Println(parseDate)
	fmt.Println(parseDate.Format(time.RFC822))
	fmt.Println(parseDate.Format(time.RFC3339))

	//Output:
	//2020-04-20 12:33:30 +0000 UTC
	//20 Apr 20 12:33 UTC
	//2020-04-20T12:33:30Z
}

//https://mingrammer.com/gobyexample/time-formatting-parsing/
func Example_Predefined_Layout_TimeParse_DateTime_RFC3339() {
	dateStr := "2021-04-18T15:04:05+09:00"
	parseDate, _ := time.Parse(time.RFC3339, dateStr)
	fmt.Println(parseDate)
	fmt.Println(parseDate.Format(time.RFC822))
	fmt.Println(parseDate.Format(time.RFC3339))

	//Output:
	//2021-04-18 15:04:05 +0900 KST
	//18 Apr 21 15:04 KST
	//2021-04-18T15:04:05+09:00
}
