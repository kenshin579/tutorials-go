package main

import (
	"fmt"
	"runtime"
	"time"
)

func callerName(skip int) string {
	const unknown = "unknown"
	pcs := make([]uintptr, 1)
	n := runtime.Callers(skip+2, pcs)
	if n < 1 {
		return unknown
	}
	frame, _ := runtime.CallersFrames(pcs).Next()
	if frame.Function == "" {
		return unknown
	}
	return frame.Function
}

// timer returns a function that prints the name of the calling
// function and the elapsed time between the call to timer and
// the call to the returned function. The returned function is
// intended to be used in a defer statement:
//
//   defer timer()()
func timer() func() {
	name := callerName(1)
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}

func main() {
	defer timer()()
	time.Sleep(time.Second * 2)
} // prints: main.main took 2s
