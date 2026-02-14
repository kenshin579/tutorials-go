package channel

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestChannelSendReceive - channel 기본 send/receive
func TestChannelSendReceive(t *testing.T) {
	ch := make(chan int) // unbuffered channel

	go func() {
		ch <- 42 // send
	}()

	value := <-ch // receive (send될 때까지 blocking)
	assert.Equal(t, 42, value)
}

// TestChannelBlocking - unbuffered channel은 send와 receive가 동시에 준비되어야 한다
func TestChannelBlocking(t *testing.T) {
	ch := make(chan string)
	var result string
	done := make(chan struct{})

	// receiver goroutine
	go func() {
		result = <-ch // sender가 보낼 때까지 blocking
		close(done)
	}()

	// sender
	ch <- "hello" // receiver가 받을 때까지 blocking
	<-done

	assert.Equal(t, "hello", result)
}

// TestChannelMultipleValues - channel로 여러 값 전달
func TestChannelMultipleValues(t *testing.T) {
	ch := make(chan int)
	var results []int
	var mu sync.Mutex

	go func() {
		for i := range 5 {
			ch <- i
		}
		close(ch) // 더 이상 보낼 값이 없으면 close
	}()

	for v := range ch { // channel이 닫힐 때까지 receive
		mu.Lock()
		results = append(results, v)
		mu.Unlock()
	}

	assert.Equal(t, []int{0, 1, 2, 3, 4}, results)
}

// TestChannelStringType - 다양한 타입의 channel
func TestChannelStringType(t *testing.T) {
	ch := make(chan string)

	go func() {
		ch <- "Go"
		ch <- "Concurrency"
	}()

	first := <-ch
	second := <-ch

	assert.Equal(t, "Go", first)
	assert.Equal(t, "Concurrency", second)
}

// TestChannelStructType - struct를 channel로 전달
func TestChannelStructType(t *testing.T) {
	type Result struct {
		Value int
		Err   error
	}

	ch := make(chan Result)

	go func() {
		ch <- Result{Value: 100, Err: nil}
	}()

	result := <-ch
	assert.Equal(t, 100, result.Value)
	assert.NoError(t, result.Err)
}
