package _trace

import (
	"net/http"
	"net/http/httptest"
	_ "net/http/pprof"
	"os"
	"runtime/pprof"
	rttrace "runtime/trace"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestPprof_Trace_동시수집 - pprof CPU 프로파일과 trace를 동시에 수집
//
// pprof와 trace는 서로 다른 관점을 제공한다:
//   - pprof: "무엇이 CPU/메모리를 소비하는가?" (샘플링, 통계적)
//   - trace: "무슨 일이 어떤 순서로 일어났는가?" (정밀, 인과관계)
//
// 두 도구를 함께 사용하면:
//  1. pprof로 CPU/메모리 핫스팟을 식별하고
//  2. trace로 해당 구간의 동시성/타이밍을 분석하여
//  3. 최적화 타겟을 정확히 결정할 수 있다
func TestPprof_Trace_동시수집(t *testing.T) {
	// 1. trace 파일 생성
	traceFile, err := os.CreateTemp(t.TempDir(), "trace_combo_*.out")
	assert.NoError(t, err)
	defer traceFile.Close()

	// 2. CPU 프로파일 파일 생성
	cpuFile, err := os.CreateTemp(t.TempDir(), "cpu_combo_*.prof")
	assert.NoError(t, err)
	defer cpuFile.Close()

	// 3. trace 수집 시작
	err = rttrace.Start(traceFile)
	assert.NoError(t, err)

	// 4. CPU 프로파일 수집 시작
	err = pprof.StartCPUProfile(cpuFile)
	assert.NoError(t, err)

	// 5. 워크로드 실행
	var wg sync.WaitGroup
	const numWorkers = 4
	const numItems = 100

	ch := make(chan int, numItems)
	for i := range numItems {
		ch <- i
	}
	close(ch)

	for range numWorkers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for n := range ch {
				// CPU 집중 작업
				sum := 0
				for i := range (n + 1) * 10_000 {
					sum += i
				}
				_ = sum
			}
		}()
	}
	wg.Wait()

	// 6. 수집 중지 (순서: trace → pprof)
	rttrace.Stop()
	pprof.StopCPUProfile()

	// 7. 두 파일 모두 데이터가 있는지 검증
	traceInfo, _ := traceFile.Stat()
	cpuInfo, _ := cpuFile.Stat()

	assert.Positive(t, traceInfo.Size(), "trace 파일이 비어있지 않아야 한다")
	assert.Positive(t, cpuInfo.Size(), "CPU 프로파일이 비어있지 않아야 한다")

	t.Logf("trace 파일: %s (%d bytes)", traceFile.Name(), traceInfo.Size())
	t.Logf("CPU 프로파일: %s (%d bytes)", cpuFile.Name(), cpuInfo.Size())
	t.Log("분석: go tool trace <trace파일> / go tool pprof <cpu프로파일>")
}

// TestPprof_HTTP_엔드포인트 - net/http/pprof 엔드포인트로 런타임 프로파일 확인
// import _ "net/http/pprof"로 기본 엔드포인트가 등록된다.
//
// 주요 엔드포인트:
//   - /debug/pprof/             : 인덱스 (사용 가능한 프로파일 목록)
//   - /debug/pprof/profile      : CPU 프로파일 (기본 30초)
//   - /debug/pprof/trace        : execution trace (기본 1초)
//   - /debug/pprof/heap         : 메모리 할당 프로파일
//   - /debug/pprof/goroutine    : goroutine 스택 덤프
//   - /debug/pprof/block        : blocking 프로파일
//   - /debug/pprof/mutex        : mutex 경합 프로파일
func TestPprof_HTTP_엔드포인트(t *testing.T) {
	// pprof 핸들러가 등록된 서버 생성
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", http.DefaultServeMux.ServeHTTP)

	server := httptest.NewServer(http.DefaultServeMux)
	defer server.Close()

	// /debug/pprof/ 인덱스 페이지 접근
	resp, err := http.Get(server.URL + "/debug/pprof/")
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// /debug/pprof/goroutine 접근
	resp2, err := http.Get(server.URL + "/debug/pprof/goroutine?debug=1")
	assert.NoError(t, err)
	defer resp2.Body.Close()
	assert.Equal(t, http.StatusOK, resp2.StatusCode)

	t.Logf("pprof 서버: %s/debug/pprof/", server.URL)
	t.Logf("trace 수집: curl -o trace.out '%s/debug/pprof/trace?seconds=5'", server.URL)
}
