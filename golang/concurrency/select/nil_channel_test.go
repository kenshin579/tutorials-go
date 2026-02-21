package _select

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNilChannelBlock - nil channel은 send/receive 모두 영원히 blocking
func TestNilChannelBlock(t *testing.T) {
	var ch chan int // nil channel (초기화 안 함)

	// nil channel에 대한 select는 해당 case를 무시한다
	select {
	case <-ch: // nil channel이므로 이 case는 절대 선택되지 않음
		t.Fatal("nil channel에서 receive할 수 없다")
	default:
		t.Log("nil channel은 select에서 무시됨")
	}
}

// TestNilChannelDisable - nil channel로 select case 동적 비활성화
func TestNilChannelDisable(t *testing.T) {
	ch1 := make(chan int, 3)
	ch2 := make(chan int, 3)

	ch1 <- 1; ch1 <- 2; ch1 <- 3; close(ch1) // ch1에 3개 값 전송 후 닫기
	ch2 <- 10; ch2 <- 20; close(ch2)           // ch2에 2개 값 전송 후 닫기

	var results []int
	// receive 전용 channel 변수로 선언 — nil 할당으로 비활성화 가능
	var active1, active2 = (<-chan int)(ch1), (<-chan int)(ch2)

	for active1 != nil || active2 != nil { // 둘 다 nil이 되면 모든 데이터 소진
		select {
		case v, ok := <-active1:
			if !ok {
				active1 = nil // 닫힌 channel → nil로 설정하면 select에서 무시됨
				continue
			}
			results = append(results, v)
		case v, ok := <-active2:
			if !ok {
				active2 = nil // 같은 방식으로 ch2도 비활성화
				continue
			}
			results = append(results, v)
		}
	}

	// 두 channel의 모든 값을 수신
	assert.Len(t, results, 5) // ch1: 3개 + ch2: 2개 = 총 5개
	t.Logf("results: %v", results)
}

// TestNilChannelMerge - nil channel을 활용한 두 정렬된 channel merge
func TestNilChannelMerge(t *testing.T) {
	// 두 개의 정렬된 channel
	ch1 := make(chan int, 4)
	ch2 := make(chan int, 4)

	for _, v := range []int{1, 3, 5, 7} { // ch1: 홀수 값
		ch1 <- v
	}
	close(ch1)

	for _, v := range []int{2, 4, 6, 8} { // ch2: 짝수 값
		ch2 <- v
	}
	close(ch2)

	// 두 channel에서 번갈아 값을 받아 하나의 슬라이스로 합침
	var merged []int
	var a, b = (<-chan int)(ch1), (<-chan int)(ch2)

	for a != nil || b != nil { // 두 source 모두 소진될 때까지 반복
		select {
		case v, ok := <-a:
			if !ok {
				a = nil // ch1 소진 → 비활성화
				continue
			}
			merged = append(merged, v)
		case v, ok := <-b:
			if !ok {
				b = nil // ch2 소진 → 비활성화
				continue
			}
			merged = append(merged, v)
		}
	}

	assert.Len(t, merged, 8) // 4개 + 4개 = 총 8개
	t.Logf("merged: %v", merged)
}
