package debugging

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestDeadlockFixed - deadlock 상황을 수정한 버전
// 원래 코드: 두 goroutine이 서로의 lock을 기다림 (circular wait)
// 수정: lock 순서를 통일 (항상 muA → muB 순서)
func TestDeadlockFixed(t *testing.T) {
	var muA, muB sync.Mutex
	var wg sync.WaitGroup

	wg.Add(2)

	// goroutine 1: muA → muB 순서
	go func() {
		defer wg.Done()
		muA.Lock()
		time.Sleep(1 * time.Millisecond)
		muB.Lock()
		// 작업 수행
		muB.Unlock()
		muA.Unlock()
	}()

	// goroutine 2: 같은 순서 muA → muB (deadlock 방지)
	go func() {
		defer wg.Done()
		muA.Lock()
		time.Sleep(1 * time.Millisecond)
		muB.Lock()
		// 작업 수행
		muB.Unlock()
		muA.Unlock()
	}()

	wg.Wait()
	t.Log("deadlock 없이 완료 (lock 순서 통일)")
}

// TestChannelDeadlockFixed - channel deadlock 수정
// 원래 코드: unbuffered channel에 같은 goroutine에서 send → 영원히 block
// 수정: buffered channel 사용 또는 별도 goroutine에서 send
func TestChannelDeadlockFixed(t *testing.T) {
	// 방법 1: buffered channel 사용
	ch := make(chan int, 1)
	ch <- 42 // buffer가 있으므로 block되지 않음
	val := <-ch
	assert.Equal(t, 42, val)

	// 방법 2: 별도 goroutine에서 send
	ch2 := make(chan int) // unbuffered
	go func() {
		ch2 <- 100
	}()
	val2 := <-ch2
	assert.Equal(t, 100, val2)
}

// TestTimeoutPreventDeadlock - timeout으로 잠재적 deadlock 방지
func TestTimeoutPreventDeadlock(t *testing.T) {
	ch := make(chan int)

	// 아무도 보내지 않는 channel에서 대기 → timeout으로 방지
	select {
	case val := <-ch:
		t.Fatalf("예상하지 않은 값 수신: %d", val)
	case <-time.After(50 * time.Millisecond):
		t.Log("timeout으로 deadlock 방지")
	}
}

// TestMutexTimeoutPattern - mutex에 timeout 적용하는 패턴
func TestMutexTimeoutPattern(t *testing.T) {
	mu := make(chan struct{}, 1) // channel을 mutex처럼 사용

	// lock 획득
	mu <- struct{}{}

	// 다른 goroutine에서 lock 시도 (timeout 포함)
	acquired := make(chan bool, 1)
	go func() {
		select {
		case mu <- struct{}{}:
			acquired <- true
		case <-time.After(50 * time.Millisecond):
			acquired <- false
		}
	}()

	result := <-acquired
	assert.False(t, result, "lock을 획득하지 못해야 함 (timeout)")

	// lock 해제
	<-mu
}
