package _trace

import (
	"context"
	"fmt"
	"os"
	rttrace "runtime/trace"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestTask_NewTask_End - trace.NewTask로 논리적 작업 단위 정의
// Task는 여러 goroutine에 걸친 논리적 작업을 하나로 묶어준다.
// go tool trace에서 Task별 latency 히스토그램을 확인할 수 있다.
func TestTask_NewTask_End(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "trace_task_*.out")
	assert.NoError(t, err)
	defer f.Close()

	err = rttrace.Start(f)
	assert.NoError(t, err)
	defer rttrace.Stop()

	ctx := context.Background()

	var wg sync.WaitGroup
	for i := range 3 {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Task 생성 - trace 뷰에서 "order-process" 라벨로 표시
			ctx, task := rttrace.NewTask(ctx, fmt.Sprintf("order-process-%d", i))
			defer task.End()

			// Task 내에서 로그 기록
			rttrace.Log(ctx, "orderID", fmt.Sprintf("ORD-%04d", i))

			// 작업 시뮬레이션
			time.Sleep(time.Duration(i+1) * time.Millisecond)
			rttrace.Log(ctx, "status", "completed")
		}()
	}
	wg.Wait()

	info, _ := f.Stat()
	assert.Positive(t, info.Size())
	t.Logf("Task trace 수집 완료: %s (%d bytes)", f.Name(), info.Size())
}

// TestRegion_WithRegion - trace.WithRegion으로 구간 측정
// Region은 단일 goroutine 내에서 특정 구간의 시간을 측정한다.
func TestRegion_WithRegion(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "trace_region_*.out")
	assert.NoError(t, err)
	defer f.Close()

	err = rttrace.Start(f)
	assert.NoError(t, err)
	defer rttrace.Stop()

	ctx := context.Background()
	ctx, task := rttrace.NewTask(ctx, "data-pipeline")
	defer task.End()

	var result int

	// WithRegion으로 구간 측정 - trace 뷰에서 구간별 시간 확인 가능
	rttrace.WithRegion(ctx, "fetch-data", func() {
		time.Sleep(2 * time.Millisecond)
		result = 42
	})

	rttrace.WithRegion(ctx, "transform-data", func() {
		result = result * 2
	})

	rttrace.WithRegion(ctx, "save-data", func() {
		time.Sleep(1 * time.Millisecond)
	})

	assert.Equal(t, 84, result)
}

// TestRegion_StartRegion_End - StartRegion/End로 유연한 구간 측정
// WithRegion과 달리 함수 클로저 바깥에서도 Region을 제어할 수 있다.
func TestRegion_StartRegion_End(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "trace_startregion_*.out")
	assert.NoError(t, err)
	defer f.Close()

	err = rttrace.Start(f)
	assert.NoError(t, err)
	defer rttrace.Stop()

	ctx := context.Background()
	ctx, task := rttrace.NewTask(ctx, "manual-region")
	defer task.End()

	// StartRegion으로 Region 시작
	region := rttrace.StartRegion(ctx, "validation")
	time.Sleep(1 * time.Millisecond)
	region.End() // 명시적으로 종료

	// 중첩(nested) Region
	rttrace.WithRegion(ctx, "outer-region", func() {
		rttrace.Log(ctx, "phase", "outer-start")

		rttrace.WithRegion(ctx, "inner-region", func() {
			rttrace.Log(ctx, "phase", "inner")
			time.Sleep(1 * time.Millisecond)
		})

		rttrace.Log(ctx, "phase", "outer-end")
	})
}

// TestLog_이벤트_마킹 - trace.Log로 trace에 커스텀 이벤트 기록
// Log는 trace 타임라인에 특정 시점의 상태를 기록하여 디버깅에 활용한다.
func TestLog_이벤트_마킹(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "trace_log_*.out")
	assert.NoError(t, err)
	defer f.Close()

	err = rttrace.Start(f)
	assert.NoError(t, err)

	ctx := context.Background()
	ctx, task := rttrace.NewTask(ctx, "http-handler")

	// 다양한 카테고리의 로그 기록
	rttrace.Log(ctx, "request", "GET /api/users")
	rttrace.Log(ctx, "db-query", "SELECT * FROM users LIMIT 10")
	rttrace.Log(ctx, "cache", "HIT user-list")
	rttrace.Log(ctx, "response", "200 OK (15 items)")

	task.End()
	rttrace.Stop()

	info, _ := f.Stat()
	assert.Positive(t, info.Size())
}
