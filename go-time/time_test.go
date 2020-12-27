package main

import (
	"fmt"
	"time"
)

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
