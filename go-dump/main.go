package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	// Capture starting number of goroutines.
	startingGs := runtime.NumGoroutine()

	leak()

	// Hold the program from terminating for 1 second to see
	// if any goroutines created by leak terminate.
	time.Sleep(time.Second)

	// Capture ending number of goroutines.
	endingGs := runtime.NumGoroutine()

	// Report the results.
	fmt.Println("========================================")
	fmt.Println("Number of goroutines before:", startingGs)
	fmt.Println("Number of goroutines after :", endingGs)
	fmt.Println("Number of goroutines leaked:", endingGs-startingGs)
}

// leak is a buggy function. It launches a goroutine that
// blocks receiving from a channel. Nothing will ever be
// sent on that channel and the channel is never closed so
// that goroutine will be blocked forever.
func leak() {
	ch := make(chan int)

	go func() {
		val := <-ch
		fmt.Println("We received a value:", val)
	}()
}
