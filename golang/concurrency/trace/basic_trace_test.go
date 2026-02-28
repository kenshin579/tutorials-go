package _trace

import (
	"os"
	rttrace "runtime/trace"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestBasicTrace_Start_Stop - trace.Start/Stop으로 기본 trace 수집
// 수집된 trace 파일은 `go tool trace trace_basic.out` 으로 분석 가능
func TestBasicTrace_Start_Stop(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "trace_basic_*.out")
	assert.NoError(t, err)
	defer f.Close()

	// trace 수집 시작
	err = rttrace.Start(f)
	assert.NoError(t, err)

	// 여러 goroutine에서 동시 작업 수행
	var wg sync.WaitGroup
	for i := range 5 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sum := 0
			for range 1_000_000 {
				sum++
			}
			t.Logf("goroutine %d: sum=%d", i, sum)
		}()
	}
	wg.Wait()

	// trace 수집 중지
	rttrace.Stop()

	// trace 파일이 생성되었는지 검증
	info, err := f.Stat()
	assert.NoError(t, err)
	assert.Positive(t, info.Size(), "trace 파일이 비어있지 않아야 한다")
	t.Logf("trace 파일: %s (크기: %d bytes)", f.Name(), info.Size())
}

// TestBasicTrace_GoTest_플래그 - go test -trace 플래그 사용 설명
// 실행: go test -v -trace=trace_gotest.out -run TestBasicTrace_GoTest_플래그
func TestBasicTrace_GoTest_플래그(t *testing.T) {
	// go test -trace=trace.out 플래그를 사용하면
	// 코드 수정 없이 테스트 실행 중 trace를 수집할 수 있다.
	//
	// 사용법:
	//   go test -v -trace=trace_gotest.out ./golang/concurrency/trace/...
	//   go tool trace trace_gotest.out
	//
	// 이 방법은 기존 테스트에 trace 코드를 추가하지 않고도
	// goroutine 스케줄링, GC 활동 등을 분석할 수 있어 편리하다.

	ch := make(chan int, 5)
	go func() {
		for i := range 5 {
			ch <- i * i
		}
		close(ch)
	}()

	var results []int
	for v := range ch {
		results = append(results, v)
	}
	assert.Len(t, results, 5)
}
