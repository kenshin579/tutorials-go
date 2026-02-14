package _select

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNilChannelBlock - nil channel은 send/receive 모두 영원히 blocking
func TestNilChannelBlock(t *testing.T) {
	var ch chan int // nil channel

	// nil channel에 대한 select는 해당 case를 무시한다
	select {
	case <-ch:
		t.Fatal("nil channel에서 receive할 수 없다")
	default:
		t.Log("nil channel은 select에서 무시됨")
	}
}

// TestNilChannelDisable - nil channel로 select case 동적 비활성화
func TestNilChannelDisable(t *testing.T) {
	ch1 := make(chan int, 3)
	ch2 := make(chan int, 3)

	ch1 <- 1
	ch1 <- 2
	ch1 <- 3
	close(ch1)

	ch2 <- 10
	ch2 <- 20
	close(ch2)

	var results []int
	var active1, active2 = (<-chan int)(ch1), (<-chan int)(ch2)

	for active1 != nil || active2 != nil {
		select {
		case v, ok := <-active1:
			if !ok {
				active1 = nil // ch1이 닫히면 nil로 설정 → 이 case 비활성화
				continue
			}
			results = append(results, v)
		case v, ok := <-active2:
			if !ok {
				active2 = nil // ch2가 닫히면 nil로 설정 → 이 case 비활성화
				continue
			}
			results = append(results, v)
		}
	}

	// 두 channel의 모든 값을 수신
	assert.Len(t, results, 5) // ch1: 3개 + ch2: 2개
	t.Logf("results: %v", results)
}

// TestNilChannelMerge - nil channel을 활용한 두 정렬된 channel merge
func TestNilChannelMerge(t *testing.T) {
	// 두 개의 정렬된 channel
	ch1 := make(chan int, 4)
	ch2 := make(chan int, 4)

	for _, v := range []int{1, 3, 5, 7} {
		ch1 <- v
	}
	close(ch1)

	for _, v := range []int{2, 4, 6, 8} {
		ch2 <- v
	}
	close(ch2)

	// 두 channel에서 번갈아 값을 받아 하나의 슬라이스로 합침
	var merged []int
	var a, b = (<-chan int)(ch1), (<-chan int)(ch2)

	for a != nil || b != nil {
		select {
		case v, ok := <-a:
			if !ok {
				a = nil
				continue
			}
			merged = append(merged, v)
		case v, ok := <-b:
			if !ok {
				b = nil
				continue
			}
			merged = append(merged, v)
		}
	}

	assert.Len(t, merged, 8)
	t.Logf("merged: %v", merged)
}
