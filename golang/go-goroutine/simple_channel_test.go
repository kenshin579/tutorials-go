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

// http://golang.site/go/article/22-Go-%EC%B1%84%EB%84%90
func Test_Int_Channel(t *testing.T) {
	// 정수형 채널을 생성한다
	ch := make(chan int)

	go func() {
		time.Sleep(time.Second * 2)
		fmt.Println("1.START")
		ch <- 123 //채널에 123을 보낸다
		fmt.Println("1.END")
	}()

	var i int
	fmt.Println("2.START")
	i = <-ch // 채널로부터 123을 받는다 (대기함)
	fmt.Println("i", i)
	fmt.Println("2.END")
}

func Test_Done(t *testing.T) {
	done := make(chan bool)

	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(i)
		}
		done <- true
	}()

	// 위의 Go루틴이 끝날 때까지 대기
	<-done
}
