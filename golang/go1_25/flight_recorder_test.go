package go1_25

import (
	"bytes"
	"runtime/trace"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_FlightRecorder_스냅샷_저장(t *testing.T) {
	// FlightRecorder 생성 (최소 1초, 최대 1MB)
	fr := trace.NewFlightRecorder(trace.FlightRecorderConfig{
		MinAge:   1 * time.Second,
		MaxBytes: 1 << 20,
	})

	// 기록 시작
	err := fr.Start()
	assert.NoError(t, err)
	assert.True(t, fr.Enabled(), "FlightRecorder가 활성화되어야 한다")

	// goroutine 작업 수행 (트레이스 데이터 생성)
	done := make(chan struct{})
	go func() {
		time.Sleep(10 * time.Millisecond)
		close(done)
	}()
	<-done

	// 스냅샷 저장
	var buf bytes.Buffer
	n, err := fr.WriteTo(&buf)
	assert.NoError(t, err)
	assert.Positive(t, n, "스냅샷 데이터가 비어있지 않아야 한다")
	t.Logf("스냅샷 크기: %d bytes", n)

	// 기록 중지
	fr.Stop()
	assert.False(t, fr.Enabled(), "Stop 후 비활성화되어야 한다")
}
