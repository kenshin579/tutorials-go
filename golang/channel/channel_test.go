package go_channel

import (
	"fmt"
	"testing"
)

// https://golangbyexample.com/channel-function-argument-go/
func Test_Bidirectional_Channel(t *testing.T) {
	ch := make(chan int, 3)
	readAndWriteToChannel(ch)
}

func readAndWriteToChannel(ch chan int) {
	ch <- 2
	s := <-ch
	fmt.Println(s)
}

func Test_Only_Write_Channel(t *testing.T) {
	ch := make(chan int, 3)
	writeToChannel(ch)
	fmt.Println(<-ch)
}

func writeToChannel(ch chan<- int) { // write 채널만 가능하도록 채널 방향 설정 : chan<- int
	ch <- 2
	//s := <-ch //reading을 할 수 없음
}

func Test_Only_Receive_Channel(t *testing.T) {
	ch := make(chan int, 3)
	ch <- 2
	readChannel(ch)
}

func readChannel(ch <-chan int) {
	s := <-ch
	fmt.Println(s)
	//ch <- 2 //channel로 write가 안됨
}

/*
채널을 통과시키는 이러한 방법은 채널을 수정하려는 경우에만 적합합니다.
이것은 매우 드문 일이며 선호하는 방법은 아닙니다.
*/
func Test_Channel_Pointer(t *testing.T) {
	ch := make(chan int, 3)
	readAndWriteChannelAsPointer(&ch)
	fmt.Println(ch)
}

func readAndWriteChannelAsPointer(ch *chan int) {
	*ch <- 2
	s := <-*ch
	*ch = nil
	fmt.Println(s)
}
