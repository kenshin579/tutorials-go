package channel

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUnbufferedChannel - unbuffered channel은 send 시 receiver가 있어야 진행
func TestUnbufferedChannel(t *testing.T) {
	ch := make(chan int) // unbuffered: 버퍼 크기 0

	go func() {
		ch <- 1 // receiver가 없으면 여기서 blocking
	}()

	val := <-ch
	assert.Equal(t, 1, val)
}

// TestBufferedChannel - buffered channel은 버퍼가 가득 차기 전까지 blocking 없이 send 가능
func TestBufferedChannel(t *testing.T) {
	ch := make(chan int, 3) // 버퍼 크기 3

	// receiver 없이도 3개까지 send 가능 (blocking 없음)
	ch <- 1
	ch <- 2
	ch <- 3
	// ch <- 4 → 여기서 blocking (버퍼 가득 참)

	assert.Equal(t, 1, <-ch)
	assert.Equal(t, 2, <-ch)
	assert.Equal(t, 3, <-ch)
}

// TestBufferedChannelCapLen - cap과 len으로 채널 상태 확인
func TestBufferedChannelCapLen(t *testing.T) {
	ch := make(chan string, 5)

	assert.Equal(t, 5, cap(ch)) // 버퍼 용량
	assert.Equal(t, 0, len(ch)) // 현재 버퍼에 대기 중인 값 수

	ch <- "a"
	ch <- "b"

	assert.Equal(t, 5, cap(ch))
	assert.Equal(t, 2, len(ch))
}

// TestBufferedChannelAsQueue - buffered channel은 FIFO 큐처럼 동작
func TestBufferedChannelAsQueue(t *testing.T) {
	ch := make(chan int, 5)

	// 순서대로 넣기
	for i := range 5 {
		ch <- i
	}

	// 순서대로 꺼내기 (FIFO)
	for i := range 5 {
		val := <-ch
		assert.Equal(t, i, val)
	}
}

// TestProducerConsumer - producer/consumer 패턴
func TestProducerConsumer(t *testing.T) {
	ch := make(chan int, 10)
	var results []int
	var mu sync.Mutex

	// Producer
	go func() {
		for i := 1; i <= 10; i++ {
			ch <- i * i // 제곱수 전달
		}
		close(ch)
	}()

	// Consumer
	for val := range ch {
		mu.Lock()
		results = append(results, val)
		mu.Unlock()
	}

	expected := []int{1, 4, 9, 16, 25, 36, 49, 64, 81, 100}
	assert.Equal(t, expected, results)
}

// TestMultipleProducers - 여러 producer가 하나의 channel에 데이터 전송
func TestMultipleProducers(t *testing.T) {
	ch := make(chan int, 20)
	var wg sync.WaitGroup

	// 3개의 producer
	for p := range 3 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range 5 {
				ch <- p*100 + i
			}
		}()
	}

	// producer가 모두 끝나면 channel close
	go func() {
		wg.Wait()
		close(ch)
	}()

	var results []int
	for val := range ch {
		results = append(results, val)
	}

	// 총 15개의 값 (3 producers x 5 values)
	assert.Len(t, results, 15)
}
