package main

import (
	"fmt"
	"testing"
	"time"
)

// https://gobyexample.com/timeouts
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

func TestTimeout_TimeAfter(t *testing.T) {
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
}

// http://golang.site/go/article/211-%EC%B1%84%EB%84%90-%ED%83%80%EC%9E%84%EC%95%84%EC%9B%83-%EA%B8%B0%EB%8A%A5
func TestAfterTime(t *testing.T) {
	//time.After()는 입력파라미터에 지정된 시간이 지나면 Ready되는 채널을 리턴한다
	timeoutChan := time.After(1 * time.Second)
	fmt.Println("timeoutChan", <-timeoutChan)
}

// https://golangbyexample.com/select-statement-with-timeout-go/
func Test(t *testing.T) {
	newsChannel := make(chan string)
	go newsFeed(newsChannel)

	printAllNews(newsChannel)
}

func printAllNews(news chan string) {
	for {
		select {
		case incomingMsg := <-news:
			fmt.Println(incomingMsg)
		case <-time.After(time.Second * 1):
			fmt.Println("Timeout: News feed finished")
			return
		}
	}
}

func newsFeed(ch chan string) {
	for i := 0; i < 2; i++ {
		time.Sleep(time.Millisecond * 400)
		ch <- fmt.Sprintf("News: %d", i+1)
	}
}
