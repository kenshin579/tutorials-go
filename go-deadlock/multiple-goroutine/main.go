package main

import "time"

// https://stackoverflow.com/questions/48548928/detect-deadlock-between-a-group-of-goroutines
func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	// goroutine 1
	go func() {
		ch1 <- 12
		ch2 <- 13 // oh oh, wrong channel. deadlock between goroutine 1 and 2
	}()

	// goroutine 2
	go func() {
		println(<-ch1)
		println(<-ch1)
	}()

	for {
		// i'm busy
		time.Sleep(time.Second)
	}

}
