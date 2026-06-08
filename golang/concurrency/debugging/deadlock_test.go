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

	// 핵심: 모든 goroutine이 동일한 순서(muA → muB)로 lock을 잡는다
	// circular wait가 형성되지 않으므로 deadlock 자체가 불가능
	go func() {
		defer wg.Done()
		muA.Lock()
		time.Sleep(1 * time.Millisecond)
		muB.Lock()
		// 작업 수행
		muB.Unlock()
		muA.Unlock()
	}()

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
	// 방법 1: buffered channel - 버퍼에 즉시 저장되므로 receive 없이도 진행 가능
	ch := make(chan int, 1)
	ch <- 42
	val := <-ch
	assert.Equal(t, 42, val)

	// 방법 2: send와 receive를 서로 다른 goroutine으로 분리 → rendezvous 성립
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

	// select + time.After: 두 channel 중 먼저 준비된 쪽이 실행됨
	// sender가 없어도 timeout이 발동하여 무한 대기를 방지
	select {
	case val := <-ch:
		t.Fatalf("예상하지 않은 값 수신: %d", val)
	case <-time.After(50 * time.Millisecond):
		t.Log("timeout으로 deadlock 방지")
	}
}

// TestMutexTimeoutPattern - mutex에 timeout 적용하는 패턴
func TestMutexTimeoutPattern(t *testing.T) {
	// 버퍼 크기 1 channel을 mutex로 활용 → sync.Mutex와 달리 timeout 적용 가능
	mu := make(chan struct{}, 1)

	// lock 획득 (버퍼에 값 넣기)
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

	// lock 해제 (버퍼에서 값 빼기)
	<-mu
}
