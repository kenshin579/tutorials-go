package goroutine

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestGoroutineCreation - goroutine 생성 기본 예제
func TestGoroutineCreation(t *testing.T) {
	// goroutine 생성 전 goroutine 수 확인
	initialCount := runtime.NumGoroutine()
	t.Logf("초기 goroutine 수: %d", initialCount)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		t.Log("goroutine 실행됨")
	}()

	wg.Wait()
}

// TestGoroutineNonDeterministicOrder - goroutine 실행 순서는 비결정적이다
func TestGoroutineNonDeterministicOrder(t *testing.T) {
	var mu sync.Mutex
	var order []int
	var wg sync.WaitGroup

	const numGoroutines = 10
	wg.Add(numGoroutines)

	for i := range numGoroutines {
		go func() {
			defer wg.Done()
			mu.Lock()
			order = append(order, i)
			mu.Unlock()
		}()
	}

	wg.Wait()

	t.Logf("실행 순서: %v", order)

	// 실행 순서가 0,1,2,...,9 순서와 다를 수 있다
	// 결정적이지 않으므로 단순히 모든 값이 포함되어 있는지 확인
	assert.Len(t, order, numGoroutines)
}

// TestGoroutineCount - runtime.NumGoroutine()으로 goroutine 수 확인
func TestGoroutineCount(t *testing.T) {
	initialCount := runtime.NumGoroutine()

	const numGoroutines = 5
	blocker := make(chan struct{})

	for range numGoroutines {
		go func() {
			<-blocker // goroutine이 살아있도록 blocking
		}()
	}

	// goroutine이 스케줄링될 시간을 줌
	time.Sleep(50 * time.Millisecond)

	currentCount := runtime.NumGoroutine()
	t.Logf("초기: %d, 현재: %d, 생성한 수: %d", initialCount, currentCount, numGoroutines)

	assert.GreaterOrEqual(t, currentCount, initialCount+numGoroutines)

	// 정리: goroutine 종료
	close(blocker)
	time.Sleep(50 * time.Millisecond)
}

// TestGoroutineLightweight - goroutine은 매우 가볍다 (수천 개 생성 가능)
func TestGoroutineLightweight(t *testing.T) {
	const numGoroutines = 10000
	var counter atomic.Int64
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for range numGoroutines {
		go func() {
			defer wg.Done()
			counter.Add(1)
		}()
	}

	wg.Wait()

	assert.Equal(t, int64(numGoroutines), counter.Load())
	t.Logf("%d개의 goroutine이 모두 완료됨", numGoroutines)
}

// TestAnonymousVsNamedGoroutine - 익명 함수 vs 이름 있는 함수로 goroutine 실행
func TestAnonymousVsNamedGoroutine(t *testing.T) {
	var wg sync.WaitGroup

	// 이름 있는 함수
	sayHello := func(name string) {
		defer wg.Done()
		fmt.Printf("Hello, %s!\n", name)
	}

	wg.Add(2)

	// 이름 있는 함수로 goroutine 실행
	go sayHello("World")

	// 익명 함수로 goroutine 실행
	go func() {
		defer wg.Done()
		fmt.Println("Hello from anonymous goroutine!")
	}()

	wg.Wait()
}

// TestGOMAXPROCS - runtime.GOMAXPROCS로 사용할 CPU 수 설정
func TestGOMAXPROCS(t *testing.T) {
	// 현재 설정된 GOMAXPROCS 값 확인
	currentProcs := runtime.GOMAXPROCS(0) // 0을 전달하면 현재 값 반환
	numCPU := runtime.NumCPU()

	t.Logf("CPU 수: %d", numCPU)
	t.Logf("현재 GOMAXPROCS: %d", currentProcs)

	// GOMAXPROCS를 1로 설정 (단일 OS 쓰레드에서 실행)
	prev := runtime.GOMAXPROCS(1)
	t.Logf("GOMAXPROCS=1 설정 (이전 값: %d)", prev)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		sum := 0
		for i := range 1000 {
			sum += i
		}
	}()

	go func() {
		defer wg.Done()
		sum := 0
		for i := range 1000 {
			sum += i
		}
	}()

	wg.Wait()

	// 원래 값 복원
	runtime.GOMAXPROCS(prev)
}
