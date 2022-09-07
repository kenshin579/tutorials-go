package util

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

func Timer() func(timerName string) {
	name := callerName(1)
	start := time.Now()
	return func(timerName string) {
		fmt.Printf("timerName:%s, %s took %v ms\n", timerName, name, time.Since(start).Milliseconds())
	}
}
