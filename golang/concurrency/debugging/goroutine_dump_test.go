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
	// 블로킹된 goroutine을 일부러 만들어 덤프에 노출시킴
	done := make(chan struct{})
	for range 3 {
		go func() {
			<-done // block 상태로 유지
		}()
	}

	time.Sleep(10 * time.Millisecond)

	buf := make([]byte, 1<<16)
	// 두 번째 인자 true → 프로세스 내 모든 goroutine의 스택을 덤프
	// hang 상태에서 어떤 goroutine이 어디서 멈췄는지 진단하는 핵심 도구
	n := runtime.Stack(buf, true)
	stackDump := string(buf[:n])

	t.Logf("=== Goroutine Dump ===\n%s", stackDump[:min(len(stackDump), 2000)])

	// 덤프에 goroutine 정보가 포함되어 있는지 확인
	assert.Contains(t, stackDump, "goroutine")

	close(done) // goroutine 정리
	time.Sleep(10 * time.Millisecond)
}

// TestNumGoroutine - goroutine 수 모니터링
func TestNumGoroutine(t *testing.T) {
	// baseline → during → after 3시점 비교가 leak 탐지의 기본 패턴
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

	// after가 baseline 수준으로 돌아오지 않으면 leak 의심
	after := runtime.NumGoroutine()
	t.Logf("after goroutines: %d (정리 완료)", after)
	assert.Less(t, after, during)
}

// TestGoroutineLeakDetection - goroutine leak 탐지 패턴
func TestGoroutineLeakDetection(t *testing.T) {
	baseline := runtime.NumGoroutine()

	runWithCleanup := func() {
		done := make(chan struct{})
		go func() {
			<-done
		}()
		// 종료 신호를 반드시 보낸다 → 함수가 끝나도 goroutine이 남지 않음
		close(done)
	}

	runWithCleanup()
	time.Sleep(50 * time.Millisecond)

	current := runtime.NumGoroutine()
	// baseline과 비교하여 leak이 없는지 확인 (스케줄러 노이즈 허용으로 +1)
	assert.LessOrEqual(t, current, baseline+1,
		"goroutine leak 발생: baseline=%d, current=%d", baseline, current)
}

// TestGoroutineStackInfo - 현재 goroutine의 스택 정보 확인
func TestGoroutineStackInfo(t *testing.T) {
	buf := make([]byte, 4096)
	// 두 번째 인자 false → 현재 goroutine 스택만 덤프 (성능 부담 적음)
	n := runtime.Stack(buf, false)
	stack := string(buf[:n])

	t.Logf("현재 goroutine 스택:\n%s", stack)

	// 현재 테스트 함수가 스택에 포함되어 있는지 확인
	assert.True(t, strings.Contains(stack, "TestGoroutineStackInfo"))
}

// TestRuntimeMemStats - 메모리 통계 확인
func TestRuntimeMemStats(t *testing.T) {
	var m runtime.MemStats
	// ReadMemStats는 GC를 강제 호출하므로 hot path에서는 주의
	runtime.ReadMemStats(&m)

	t.Logf("Alloc = %d KB", m.Alloc/1024)
	t.Logf("TotalAlloc = %d KB", m.TotalAlloc/1024)
	t.Logf("Sys = %d KB", m.Sys/1024)
	t.Logf("NumGC = %d", m.NumGC)
	t.Logf("NumGoroutine = %d", runtime.NumGoroutine())

	assert.Greater(t, m.Alloc, uint64(0))
}
