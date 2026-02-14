package _select

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestSelectBasic - select로 여러 channel 동시 대기
func TestSelectBasic(t *testing.T) {
	ch1 := make(chan string, 1)
	ch2 := make(chan string, 1)

	ch1 <- "from ch1"

	var result string
	select {
	case msg := <-ch1:
		result = msg
	case msg := <-ch2:
		result = msg
	}

	assert.Equal(t, "from ch1", result)
}

// TestSelectMultipleReady - 여러 case가 동시에 준비되면 랜덤 선택
func TestSelectMultipleReady(t *testing.T) {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)

	ch1Count := 0
	ch2Count := 0
	iterations := 1000

	for range iterations {
		ch1 <- 1
		ch2 <- 2

		select {
		case <-ch1:
			ch1Count++
		case <-ch2:
			ch2Count++
		}

		// 선택되지 않은 channel 비우기
		select {
		case <-ch1:
		case <-ch2:
		default:
		}
	}

	t.Logf("ch1 선택: %d, ch2 선택: %d", ch1Count, ch2Count)
	// 완전히 균등하지는 않지만, 한쪽으로 치우치지 않아야 함
	assert.Greater(t, ch1Count, 0)
	assert.Greater(t, ch2Count, 0)
}

// TestSelectDefault - default case로 non-blocking 동작
func TestSelectDefault(t *testing.T) {
	ch := make(chan int)

	var result string
	select {
	case val := <-ch:
		result = "received: " + string(rune(val))
	default:
		result = "no data available" // channel에 데이터가 없으면 즉시 실행
	}

	assert.Equal(t, "no data available", result)
}

// TestSelectDefaultNonBlockingSend - default로 non-blocking send
func TestSelectDefaultNonBlockingSend(t *testing.T) {
	ch := make(chan int, 1)
	ch <- 1 // 버퍼 가득 참

	sent := false
	select {
	case ch <- 2:
		sent = true
	default:
		sent = false // 버퍼가 가득 차서 send 불가 → default 실행
	}

	assert.False(t, sent)
}

// TestSelectWithTimer - select로 주기적 작업
func TestSelectWithTimer(t *testing.T) {
	done := make(chan struct{})
	var ticks int

	go func() {
		ticker := time.NewTicker(10 * time.Millisecond)
		defer ticker.Stop()
		timeout := time.After(55 * time.Millisecond)

		for {
			select {
			case <-ticker.C:
				ticks++
			case <-timeout:
				close(done)
				return
			}
		}
	}()

	<-done
	t.Logf("ticks: %d", ticks)
	assert.GreaterOrEqual(t, ticks, 3)
}
