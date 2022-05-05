package main

import (
	"fmt"
	"testing"
	"time"
)

/*
1.Goroutine:
*/
func Test_Hello_World(t *testing.T) {
	go say("world")
	say("hello")
}

func say(s string) {
	for i := 0; i < 5; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Printf("%s ", s)
	}
}

/*
2.Channels :
ch <- v    // Send v to channel ch.
v := <-ch  // Receive from ch, and
           // assign value to v.
*/
func Test_Channel_Usage(t *testing.T) {
	s := []int{7, 2, 8, -9, 4, 0}

	c := make(chan int)
	go sum(s[:len(s)/2], c)
	go sum(s[len(s)/2:], c)
	x, y := <-c, <-c // receive from c :

	fmt.Println(x, y, x+y)
}

func sum(s []int, c chan int) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	c <- sum // send sum to c :
}

/*
3.Buffered Channels
*/
func Test_Buffered_Channel(t *testing.T) {
	ch := make(chan int, 2)
	ch <- 1
	ch <- 2
	fmt.Println(<-ch)
	fmt.Println(<-ch)
}

/*
4.Close 체크
*/
func Test_Channel_Close_Check(t *testing.T) {
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

/*
5.Range Loop
*/
func Test_Range_Loop(t *testing.T) {
	c := make(chan int, 10)
	go fibonacci(cap(c), c)
	for i := range c {
		fmt.Println(i)
	}
}

func fibonacci(n int, c chan int) {
	x, y := 0, 1
	for i := 0; i < n; i++ {
		time.Sleep(100 * time.Millisecond)
		c <- x
		x, y = y, x+y
	}
	close(c) //sender에서 보낼 데이터가 없다고 알려주는 거임
}
