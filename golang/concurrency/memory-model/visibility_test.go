package memory_model

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestVisibilityProblem - 공유 변수 visibility 문제
// 한 goroutine에서 변경한 값이 다른 goroutine에서 보이지 않을 수 있음
func TestVisibilityProblem(t *testing.T) {
	var data int
	var ready atomic.Bool

	go func() {
		data = 42
		ready.Store(true) // atomic으로 happens-before 보장
	}()

	// ready가 true가 될 때까지 대기
	for !ready.Load() {
		// busy wait
	}

	// atomic 연산이 happens-before를 보장하므로 data=42가 보임
	assert.Equal(t, 42, data)
}

// TestVisibilityWithChannel - channel은 happens-before를 보장
func TestVisibilityWithChannel(t *testing.T) {
	var data int
	done := make(chan struct{})

	go func() {
		data = 100
		close(done) // channel close는 happens-before 보장
	}()

	<-done
	assert.Equal(t, 100, data)
}

// TestVisibilityWithMutex - Mutex도 happens-before를 보장
func TestVisibilityWithMutex(t *testing.T) {
	var mu sync.Mutex
	var data int

	go func() {
		mu.Lock()
		data = 200
		mu.Unlock()
	}()

	// 약간의 시간을 두고 읽기
	mu.Lock()
	// Unlock → Lock 은 happens-before 관계
	// 하지만 어느 goroutine이 먼저 Lock을 잡을지는 비결정적
	val := data
	mu.Unlock()

	t.Logf("data = %d (0 또는 200)", val)
	assert.True(t, val == 0 || val == 200)
}

// TestHappensBeforeRules - Go의 happens-before 규칙 정리
func TestHappensBeforeRules(t *testing.T) {
	// 1. 같은 goroutine 내에서는 순서 보장
	x := 1
	x = 2
	assert.Equal(t, 2, x) // 항상 보장

	// 2. channel send는 receive보다 happens-before
	ch := make(chan int)
	var result int
	go func() {
		result = 42
		ch <- 1 // send happens-before receive
	}()
	<-ch
	assert.Equal(t, 42, result)

	// 3. channel close는 receive(zero value)보다 happens-before
	ch2 := make(chan int)
	var result2 int
	go func() {
		result2 = 99
		close(ch2)
	}()
	<-ch2
	assert.Equal(t, 99, result2)

	// 4. sync.Once - Do()는 한 번만 실행되고 모든 호출에 happens-before
	var once sync.Once
	var initialized int
	var wg sync.WaitGroup

	wg.Add(10)
	for range 10 {
		go func() {
			defer wg.Done()
			once.Do(func() {
				initialized = 1
			})
			assert.Equal(t, 1, initialized)
		}()
	}
	wg.Wait()
}
