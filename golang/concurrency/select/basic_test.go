package _select

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestSelectBasic - select로 여러 channel 동시 대기
func TestSelectBasic(t *testing.T) {
	ch1 := make(chan string, 1) // 버퍼 1짜리 channel
	ch2 := make(chan string, 1)

	ch1 <- "from ch1" // ch1에만 값이 있으므로 ch1 case가 선택됨

	var result string
	select {
	case msg := <-ch1: // ch1에서 먼저 데이터가 오면 실행
		result = msg
	case msg := <-ch2: // ch2에서 먼저 데이터가 오면 실행
		result = msg
	}

	assert.Equal(t, "from ch1", result)
}

// TestSelectMultipleReady - 여러 case가 동시에 준비되면 랜덤 선택
func TestSelectMultipleReady(t *testing.T) {
	ch1 := make(chan int, 1) // 버퍼 1짜리 channel
	ch2 := make(chan int, 1)

	ch1Count := 0
	ch2Count := 0
	iterations := 1000

	for range iterations { // 1000번 반복하여 선택 비율 확인
		ch1 <- 1 // 두 channel에 동시에 값을 넣어 둘 다 준비 상태로 만듦
		ch2 <- 2

		select {
		case <-ch1: // 두 case 모두 준비됐으므로 runtime이 무작위 선택
			ch1Count++
		case <-ch2:
			ch2Count++
		}

		// 선택되지 않은 channel의 남은 값 비우기
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
	ch := make(chan int) // unbuffered channel (데이터 없음)

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
	case ch <- 2: // 버퍼가 가득 차서 send 불가
		sent = true
	default:
		sent = false // send 실패 시 즉시 default 실행
	}

	assert.False(t, sent)
}

// TestSelectWithTimer - select로 주기적 작업 + timeout 조합
func TestSelectWithTimer(t *testing.T) {
	done := make(chan struct{})
	var ticks int

	go func() {
		ticker := time.NewTicker(10 * time.Millisecond) // 10ms마다 tick
		defer ticker.Stop()
		timeout := time.After(55 * time.Millisecond) // 55ms 후 timeout

		for {
			select {
			case <-ticker.C: // tick마다 카운트 증가
				ticks++
			case <-timeout: // timeout 발생 시 종료
				close(done)
				return
			}
		}
	}()

	<-done
	t.Logf("ticks: %d", ticks)
	assert.GreaterOrEqual(t, ticks, 3)
}
