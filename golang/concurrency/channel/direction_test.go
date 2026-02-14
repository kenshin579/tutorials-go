package channel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// send-only channel: chan<- T (보내기만 가능)
// receive-only channel: <-chan T (받기만 가능)

// producer는 send-only channel을 받음
func produce(ch chan<- int, values []int) {
	for _, v := range values {
		ch <- v
	}
	close(ch)
}

// consumer는 receive-only channel을 받음
func consume(ch <-chan int) []int {
	var results []int
	for v := range ch {
		results = append(results, v)
	}
	return results
}

// TestChannelDirection - channel 방향 제한
func TestChannelDirection(t *testing.T) {
	ch := make(chan int, 5)

	// 양방향 channel을 send-only 또는 receive-only로 변환 가능
	go produce(ch, []int{1, 2, 3, 4, 5})
	results := consume(ch)

	assert.Equal(t, []int{1, 2, 3, 4, 5}, results)
}

// TestSendOnlyChannel - send-only channel에서는 receive 불가 (컴파일 에러)
func TestSendOnlyChannel(t *testing.T) {
	ch := make(chan int, 1)

	var sendOnly chan<- int = ch // 양방향 → send-only 변환
	sendOnly <- 42

	// val := <-sendOnly // 컴파일 에러: receive from send-only channel

	val := <-ch // 원래 양방향 채널에서는 receive 가능
	assert.Equal(t, 42, val)
}

// TestReceiveOnlyChannel - receive-only channel에서는 send 불가 (컴파일 에러)
func TestReceiveOnlyChannel(t *testing.T) {
	ch := make(chan int, 1)
	ch <- 42

	var recvOnly <-chan int = ch // 양방향 → receive-only 변환
	val := <-recvOnly

	// recvOnly <- 100 // 컴파일 에러: send to receive-only channel

	assert.Equal(t, 42, val)
}

// TestDirectionInFunctionSignature - 함수 시그니처에서 channel 방향 제한
func TestDirectionInFunctionSignature(t *testing.T) {
	// 변환 함수: 입력 channel에서 값을 읽어 2배로 변환 후 출력 channel에 쓰기
	doubler := func(in <-chan int, out chan<- int) {
		for v := range in {
			out <- v * 2
		}
		close(out)
	}

	input := make(chan int, 3)
	output := make(chan int, 3)

	input <- 1
	input <- 2
	input <- 3
	close(input)

	go doubler(input, output)

	var results []int
	for v := range output {
		results = append(results, v)
	}

	assert.Equal(t, []int{2, 4, 6}, results)
}
