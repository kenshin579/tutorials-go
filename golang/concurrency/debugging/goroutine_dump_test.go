package debugging

import (
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestGoroutineDump - runtime.Stack()으로 goroutine 상태 덤프
func TestGoroutineDump(t *testing.T) {
	// goroutine 몇 개 생성
	done := make(chan struct{})
	for range 3 {
		go func() {
			<-done // block 상태로 유지
		}()
	}

	time.Sleep(10 * time.Millisecond)

	// 현재 goroutine 스택 덤프
	buf := make([]byte, 1<<16)
	n := runtime.Stack(buf, true) // true = 모든 goroutine
	stackDump := string(buf[:n])

	t.Logf("=== Goroutine Dump ===\n%s", stackDump[:min(len(stackDump), 2000)])

	// 덤프에 goroutine 정보가 포함되어 있는지 확인
	assert.Contains(t, stackDump, "goroutine")

	close(done) // goroutine 정리
	time.Sleep(10 * time.Millisecond)
}

// TestNumGoroutine - goroutine 수 모니터링
func TestNumGoroutine(t *testing.T) {
	baseline := runtime.NumGoroutine()
	t.Logf("baseline goroutines: %d", baseline)

	done := make(chan struct{})
	for range 10 {
		go func() {
			<-done
		}()
	}

	time.Sleep(10 * time.Millisecond)
	during := runtime.NumGoroutine()
	t.Logf("during goroutines: %d (10개 추가)", during)
	assert.GreaterOrEqual(t, during, baseline+10)

	close(done)
	time.Sleep(50 * time.Millisecond)

	after := runtime.NumGoroutine()
	t.Logf("after goroutines: %d (정리 완료)", after)
	assert.Less(t, after, during)
}

// TestGoroutineLeakDetection - goroutine leak 탐지 패턴
func TestGoroutineLeakDetection(t *testing.T) {
	baseline := runtime.NumGoroutine()

	// 의도적으로 goroutine을 제대로 정리하는 함수
	runWithCleanup := func() {
		done := make(chan struct{})
		go func() {
			<-done
		}()
		close(done) // 반드시 정리
	}

	runWithCleanup()
	time.Sleep(50 * time.Millisecond)

	current := runtime.NumGoroutine()
	// baseline과 비교하여 leak이 없는지 확인
	assert.LessOrEqual(t, current, baseline+1,
		"goroutine leak 발생: baseline=%d, current=%d", baseline, current)
}

// TestGoroutineStackInfo - 현재 goroutine의 스택 정보 확인
func TestGoroutineStackInfo(t *testing.T) {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false) // false = 현재 goroutine만
	stack := string(buf[:n])

	t.Logf("현재 goroutine 스택:\n%s", stack)

	// 현재 테스트 함수가 스택에 포함되어 있는지 확인
	assert.True(t, strings.Contains(stack, "TestGoroutineStackInfo"))
}

// TestRuntimeMemStats - 메모리 통계 확인
func TestRuntimeMemStats(t *testing.T) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	t.Logf("Alloc = %d KB", m.Alloc/1024)
	t.Logf("TotalAlloc = %d KB", m.TotalAlloc/1024)
	t.Logf("Sys = %d KB", m.Sys/1024)
	t.Logf("NumGC = %d", m.NumGC)
	t.Logf("NumGoroutine = %d", runtime.NumGoroutine())

	assert.Greater(t, m.Alloc, uint64(0))
}
