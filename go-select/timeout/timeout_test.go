package main

import (
	"fmt"
	"testing"
	"time"
)

//https://gobyexample.com/timeouts
func TestTimeout_Sleep1(t *testing.T) {
	c1 := make(chan string, 1)
	go func() {
		fmt.Println("1111") //1
		time.Sleep(time.Second * 2)
		c1 <- "result 1"
		fmt.Println("2222") //2
	}()

	fmt.Println("3333", <-c1) //2
}

func TestTimeout_Sleep2(t *testing.T) {
	go func() {
		fmt.Println("1111") //1
		time.Sleep(time.Second * 2)
		fmt.Println("2222") //2
	}()
	fmt.Println("3333") //1
	time.Sleep(time.Second * 3)
}

//todo : 다시 확인하기
func TestTimeout_All(t *testing.T) {
	c1 := make(chan string, 1)
	go func() {
		time.Sleep(2 * time.Second)
		c1 <- "result 1"
	}()

	select {
	case res := <-c1:
		fmt.Println("result:c1", res)
	case <-time.After(1 * time.Second):
		fmt.Println("timeout 1")
	}

	c2 := make(chan string, 1)
	go func() {
		time.Sleep(2 * time.Second)
		c2 <- "result 2"
	}()
	select {
	case res := <-c2:
		fmt.Println("result:c2", res)
	case <-time.After(3 * time.Second):
		fmt.Println("timeout 2")
	}
}
