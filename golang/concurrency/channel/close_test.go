package channel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestChannelClose - channel close 기본 동작
func TestChannelClose(t *testing.T) {
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3
	close(ch) // 더 이상 값을 보내지 않음을 알림

	// close된 channel에서도 버퍼에 남은 값을 receive 가능
	assert.Equal(t, 1, <-ch)
	assert.Equal(t, 2, <-ch)
	assert.Equal(t, 3, <-ch)
}

// TestReceiveFromClosedChannel - 닫힌 channel에서 receive하면 zero value 반환
func TestReceiveFromClosedChannel(t *testing.T) {
	ch := make(chan int, 1)
	ch <- 42
	close(ch)

	// 버퍼에 값이 있으면 정상 반환
	val, ok := <-ch
	assert.Equal(t, 42, val)
	assert.True(t, ok) // 값이 유효함

	// 버퍼가 비고 닫힌 channel → zero value + false
	val, ok = <-ch
	assert.Equal(t, 0, val)  // int의 zero value
	assert.False(t, ok)      // channel이 닫히고 비어있음
}

// TestReceiveFromClosedStringChannel - string channel의 zero value
func TestReceiveFromClosedStringChannel(t *testing.T) {
	ch := make(chan string)
	close(ch)

	val, ok := <-ch
	assert.Equal(t, "", val) // string의 zero value
	assert.False(t, ok)
}

// TestRangeOverChannel - range로 channel의 모든 값 수신
func TestRangeOverChannel(t *testing.T) {
	ch := make(chan int, 5)

	go func() {
		for i := 1; i <= 5; i++ {
			ch <- i
		}
		close(ch) // range가 종료되려면 channel이 닫혀야 한다
	}()

	var results []int
	for v := range ch { // channel이 닫힐 때까지 반복
		results = append(results, v)
	}

	assert.Equal(t, []int{1, 2, 3, 4, 5}, results)
}

// TestCloseResponsibility - close 책임은 sender에게 있다
func TestCloseResponsibility(t *testing.T) {
	// 패턴: sender가 channel을 생성하고, 데이터를 보내고, close한다
	generator := func() <-chan int {
		ch := make(chan int)
		go func() {
			defer close(ch) // sender가 close 책임
			for i := range 5 {
				ch <- i
			}
		}()
		return ch // receive-only channel 반환
	}

	var results []int
	for v := range generator() {
		results = append(results, v)
	}

	assert.Equal(t, []int{0, 1, 2, 3, 4}, results)
}

// TestChannelSignaling - 빈 struct channel을 신호용으로 사용
func TestChannelSignaling(t *testing.T) {
	done := make(chan struct{}) // 데이터 없이 신호만 전달

	go func() {
		// 작업 수행...
		close(done) // 완료 신호 (모든 receiver에게 전달)
	}()

	<-done // 완료 대기
	t.Log("작업 완료 신호 수신")
}
