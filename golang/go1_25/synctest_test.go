package go1_25

import (
	"testing"
	"testing/synctest"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Synctest_가상시간_타이머(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		// 가상 시간에서 1초 타이머 생성
		result := make(chan string, 1)

		go func() {
			time.Sleep(1 * time.Second) // 가상 시간에서는 즉시 진행
			result <- "완료"
		}()

		// 가상 시간 1초 진행 후 모든 goroutine 대기
		time.Sleep(1 * time.Second)
		synctest.Wait()

		select {
		case v := <-result:
			assert.Equal(t, "완료", v)
		default:
			t.Fatal("타이머가 만료되지 않았다")
		}
	})
}

func Test_Synctest_Wait_goroutine_동기화(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ch := make(chan int, 3)

		// 여러 goroutine 시작
		go func() { ch <- 1 }()
		go func() { ch <- 2 }()
		go func() { ch <- 3 }()

		// 모든 goroutine이 블록될 때까지 대기
		synctest.Wait()

		// 채널에 3개의 값이 모두 들어있어야 함
		assert.Equal(t, 3, len(ch))
	})
}
