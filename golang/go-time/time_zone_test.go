package main

import (
	"fmt"
	"testing"
	"time"
)

func ExampleTimeZone() {
	timeStr := "2019-01-12 00:00:00.000"

	customLayout := "2006-01-02 15:04:05.000"
	utc, _ := time.Parse(customLayout, timeStr)

	location, _ := time.LoadLocation("Asia/Seoul")
	kst := utc.In(location)

	fmt.Println(utc)
	fmt.Println(kst)

	//Output:
	//2019-01-12 00:00:00 +0000 UTC
	//2019-01-12 09:00:00 +0900 KST
}

func TestLocation(t *testing.T) {
	now := time.Now()
	fmt.Println(now)
	fmt.Println(now.Location()) //Local
}
