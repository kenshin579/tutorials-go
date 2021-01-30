package go_select

import (
	"fmt"
	"testing"
	"time"
)

//https://hamait.tistory.com/1017
func Test_Select(t *testing.T) {
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
