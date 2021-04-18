package main

import (
	"fmt"
	"time"
)

func Example_ParseDuration() {
	h, _ := time.ParseDuration("4h30m")
	fmt.Printf("%.1f\n", h.Hours())
	fmt.Println(h)

	//Output:
	//4.5
	//4h30m0s
}

func ExampleTimeDuration() {
	duration := 5 * time.Minute

	fmt.Println(duration.Minutes())
	fmt.Println(duration.Seconds())
	fmt.Println(duration.String())

	//Output:
	//5
	//300
	//5m0s
}

func ExampleTruncate() {
	d, err := time.ParseDuration("1h15m30.918273645s")
	if err != nil {
		panic(err)
	}

	trunc := []time.Duration{
		time.Nanosecond,
		time.Microsecond,
		time.Millisecond,
		time.Second,
		2 * time.Second,
		time.Minute,
		10 * time.Minute,
		time.Hour,
	}

	for _, t := range trunc {
		fmt.Printf("d.Truncate(%6s) = %s\n", t, d.Truncate(t).String())
	}

	//Output:
	//d.Truncate(   1ns) = 1h15m30.918273645s
	//d.Truncate(   1µs) = 1h15m30.918273s
	//d.Truncate(   1ms) = 1h15m30.918s
	//d.Truncate(    1s) = 1h15m30s
	//d.Truncate(    2s) = 1h15m30s
	//d.Truncate(  1m0s) = 1h15m0s
	//d.Truncate( 10m0s) = 1h10m0s
	//d.Truncate(1h0m0s) = 1h0m0s
}
