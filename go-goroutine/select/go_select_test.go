package _select

import (
	"fmt"
	"testing"
	"time"
)

/*
1.Select : goroutine wait on multiple communication operations
- A select blocks until one of its cases can run, then it executes that case. It chooses one at random if multiple are ready.
goroutine -> [c, quit] -> main:fabonacci
main:fabonacci:select
*/
func Test_Select(t *testing.T) {
	c := make(chan int)
	quit := make(chan int)

	go func() {
		for i := 0; i < 10; i++ {
			time.Sleep(100 * time.Millisecond)
			fmt.Println(<-c) //받을 떄까지 block
		}
		quit <- 0
	}()
	fibonacci(c, quit)
}

func Test_Default_Select(t *testing.T) {
	tick := time.Tick(100 * time.Millisecond)
	boom := time.After(500 * time.Millisecond)
	for {
		select {
		case <-tick:
			fmt.Println("tick.")
		case <-boom:
			fmt.Println("BOOM!")
			return
		default:
			fmt.Println("    .")
			time.Sleep(50 * time.Millisecond)
		}
	}
}

//https://hamait.tistory.com/1017
func Test_Select2(t *testing.T) {
	c1 := make(chan string)
	c2 := make(chan string)

	go func() {
		for {
			time.Sleep(5 * time.Second)
			c1 <- "one"
		}
	}()

	go func() {
		for {
			time.Sleep(10 * time.Second)
			c2 <- "two"
		}
	}()

	for {
		fmt.Println("start select-----------------")

		//처널에 값이 들어올 때까지 select문에서 블록된다
		select {
		case msg1 := <-c1:
			fmt.Println("received", msg1)
		case msg2 := <-c2:
			fmt.Println("received", msg2)
		}
		fmt.Println("end select-----------------")
	}
}

func fibonacci(c, quit chan int) {
	x, y := 0, 1
	for {
		select {
		case c <- x: //channel로 데이터를 보냄
			x, y = y, x+y
		case <-quit:
			fmt.Println("quit")
			return
		}
	}
}
