package main

import (
	"fmt"
	"testing"
)

// https://yourbasic.org/golang/detect-deadlock/
func TestDeadlock1(t *testing.T) {
	ch := make(chan int)
	ch <- 1 //deadline
	fmt.Println(<-ch)
}

func TestDeadlock1_Fix(t *testing.T) {
	ch := make(chan int)

	go func() {
		fmt.Println("START")
		fmt.Println(<-ch)
		fmt.Println("END")
	}()
	ch <- 1

}

// todo: deadlock이 안생김
func TestDeadlock2(t *testing.T) {
	c1 := make(chan int)

	go func() {
		g1 := <-c1 // wait for value
		fmt.Println("get g1: ", g1)
	}()

	fmt.Println("push c1: ")
	c1 <- 10 // send value and wait until it is received.
}
