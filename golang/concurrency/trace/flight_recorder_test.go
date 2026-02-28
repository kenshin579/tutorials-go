package _trace

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	rttrace "runtime/trace"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestFlightRecorder_HTTP서버_Latency_감지 - HTTP 서버에서 느린 요청 발생 시 자동 스냅샷 캡처
// FlightRecorder는 링 버퍼 방식으로 항상 trace를 수집하다가,
// 이상 징후(예: latency 초과) 발생 시 WriteTo로 스냅샷을 저장한다.
//
// 프로덕션 환경에서 간헐적으로 발생하는 tail latency 문제를
// 사후에 분석할 수 있는 "블랙박스 레코더" 역할을 한다.
func TestFlightRecorder_HTTP서버_Latency_감지(t *testing.T) {
	const latencyThreshold = 50 * time.Millisecond

	// FlightRecorder 설정
	fr := rttrace.NewFlightRecorder(rttrace.FlightRecorderConfig{
		MinAge:   200 * time.Millisecond,
		MaxBytes: 1 << 20, // 1 MiB
	})

	err := fr.Start()
	assert.NoError(t, err)
	defer fr.Stop()
	assert.True(t, fr.Enabled(), "FlightRecorder가 활성화되어야 한다")

	// 가변 응답 시간 서버: 요청 번호에 따라 지연 시간이 달라짐
	requestCount := 0
	var mu sync.Mutex
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		count := requestCount
		requestCount++
		mu.Unlock()

		// 5번째 요청마다 느린 응답 시뮬레이션
		if count%5 == 4 {
			time.Sleep(80 * time.Millisecond)
		} else {
			time.Sleep(5 * time.Millisecond)
		}
		fmt.Fprintf(w, "request-%d", count)
	}))
	defer server.Close()

	// 클라이언트: 요청 보내면서 latency 초과 시 스냅샷 캡처
	var snapshots []int64
	client := server.Client()

	for i := range 10 {
		start := time.Now()

		resp, err := client.Get(server.URL)
		assert.NoError(t, err)
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()

		elapsed := time.Since(start)
		t.Logf("요청 %d: %v", i, elapsed)

		// latency 초과 시 FlightRecorder 스냅샷 저장
		if elapsed > latencyThreshold {
			var buf bytes.Buffer
			n, err := fr.WriteTo(&buf)
			assert.NoError(t, err)
			snapshots = append(snapshots, n)
			t.Logf("  → 스냅샷 캡처! (크기: %d bytes, latency: %v)", n, elapsed)
		}
	}

	// 느린 요청이 있었으므로 최소 1개 이상의 스냅샷이 캡처되어야 한다
	assert.NotEmpty(t, snapshots, "latency 초과 시 스냅샷이 캡처되어야 한다")
	for _, size := range snapshots {
		assert.Positive(t, size, "스냅샷 데이터가 비어있지 않아야 한다")
	}
}

// TestFlightRecorder_Start_Stop_라이프사이클 - FlightRecorder의 기본 라이프사이클
func TestFlightRecorder_Start_Stop_라이프사이클(t *testing.T) {
	fr := rttrace.NewFlightRecorder(rttrace.FlightRecorderConfig{
		MinAge:   100 * time.Millisecond,
		MaxBytes: 512 << 10, // 512 KiB
	})

	// 시작 전에는 비활성화 상태
	assert.False(t, fr.Enabled())

	// 시작
	err := fr.Start()
	assert.NoError(t, err)
	assert.True(t, fr.Enabled())

	// goroutine 활동으로 trace 데이터 생성
	var wg sync.WaitGroup
	for range 5 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(10 * time.Millisecond)
		}()
	}
	wg.Wait()

	// 스냅샷 저장
	var buf bytes.Buffer
	n, err := fr.WriteTo(&buf)
	assert.NoError(t, err)
	assert.Positive(t, n, "스냅샷 데이터가 있어야 한다")
	t.Logf("스냅샷 크기: %d bytes", n)

	// 중지
	fr.Stop()
	assert.False(t, fr.Enabled(), "Stop 후 비활성화되어야 한다")
}
