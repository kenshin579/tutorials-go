package main

import (
	"fmt"
	"testing"
)

func Test_Check_Close_Channel(t *testing.T) {
	ch := make(chan int, 1)
	go func() {
		ch <- 1
	}()

	if a, ok := <-ch; ok {
		fmt.Println(a, ok)
	}
	close(ch)

	if _, ok := <-ch; !ok {
		fmt.Println("Channel is closed")
	}
}
