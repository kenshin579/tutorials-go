package main

import (
	"fmt"
	"testing"
	"time"
)

/*
- 문자열 파싱 <-- todo : 이거 부터 하면 됨
- 날짜 계산
- days in month
- time zones
- measure of piece of code
- time.Now(), time.Date(), time.Parse()

*/

func Test(t *testing.T) {
	// The date we're trying to parse, work with and format
	myDateString := "2018-01-20 04:35"
	fmt.Println("My Starting Date:\t", myDateString)

	// Parse the date string into Go's time object
	// The 1st param specifies the format, 2nd is our date string
	myDate, err := time.Parse("2006-01-02 15:04", myDateString)
	if err != nil {
		panic(err)
	}

	// Format uses the same formatting style as parse, or we can use a pre-made constant
	fmt.Println("My Date Reformatted:\t", myDate.Format(time.RFC822))

	// In Y-m-d
}

func TestParse_Date_String(t *testing.T) {
	dateStr := "2019-03-20 02:22:30"
	parseDate, _ := time.Parse("1999-01-02 00:00:00", dateStr)
	fmt.Println(parseDate.Format("2000-01-02"))
	//fmt.Println(parseDate.Year())
}

func Example_diff() {
	currentTime := time.Date(2020, 1, 1, 1, 48, 32, 0, time.UTC)
	totalTime := time.Date(2020, 1, 1, 6, 39, 57, 0, time.UTC)
	diffTime := totalTime.Sub(currentTime)
	totalTimeInSeconds := float64(totalTime.Second() + totalTime.Hour()*24*60 + totalTime.Minute()*60)
	fmt.Println("currentTime", currentTime)
	fmt.Println("totalTime", totalTime)
	fmt.Println("totalTimeInSeconds", totalTimeInSeconds)
	fmt.Println("diffTime", diffTime.Seconds())
	fmt.Println("T", diffTime.Seconds()/totalTimeInSeconds)

	//Output:

}
