// Grafana Pyroscope SDK 기본 연동 예제
//
// Pyroscope 서버에 CPU, 메모리, 고루틴, 뮤텍스, 블로킹 프로파일을 Push 모드로 전송한다.
// Labels(TagWrapper)를 사용하여 워크로드 유형별로 프로파일 데이터를 구분할 수 있다.
//
// 실행:
//
//	PYROSCOPE_SERVER=http://localhost:4040 go run .
//
// Docker Compose:
//
//	docker compose up -d  (루트 디렉토리에서)
package main

import (
	"context"
	"log"
	"math"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/grafana/pyroscope-go"
)

func main() {
	// 뮤텍스/블로킹 프로파일은 기본 비활성화이므로 명시적으로 활성화
	runtime.SetMutexProfileFraction(5)
	runtime.SetBlockProfileRate(5)

	serverAddr := getEnv("PYROSCOPE_SERVER", "http://localhost:4040")

	profiler, err := pyroscope.Start(pyroscope.Config{
		ApplicationName: "simple.golang.app",
		ServerAddress:   serverAddr,
		Logger:          pyroscope.StandardLogger,
		Tags:            map[string]string{"hostname": hostname()},
		ProfileTypes: []pyroscope.ProfileType{
			// 기본 활성화 프로파일
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,
			// 추가 프로파일 (선택)
			pyroscope.ProfileGoroutines,
			pyroscope.ProfileMutexCount,
			pyroscope.ProfileMutexDuration,
			pyroscope.ProfileBlockCount,
			pyroscope.ProfileBlockDuration,
		},
	})
	if err != nil {
		log.Fatalf("pyroscope 시작 실패: %v", err)
	}
	defer profiler.Stop()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	// 워크로드별 Labels로 프로파일 데이터 구분
	go runWorkload(ctx, "cpu", cpuWork)
	go runWorkload(ctx, "memory", memoryWork)
	go runWorkload(ctx, "mutex", mutexWork)

	log.Printf("Pyroscope 프로파일링 시작 (서버: %s) - Ctrl+C로 종료", serverAddr)
	<-ctx.Done()
	log.Println("종료합니다")
}

// runWorkload는 지정된 label과 함께 워크로드를 반복 실행한다.
// Pyroscope TagWrapper로 감싸서 Flame Graph에서 워크로드별로 필터링할 수 있다.
func runWorkload(ctx context.Context, label string, work func()) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			pyroscope.TagWrapper(ctx, pyroscope.Labels("workload", label), func(c context.Context) {
				work()
			})
		}
	}
}

// cpuWork는 피보나치 연산으로 CPU 부하를 생성한다.
func cpuWork() {
	n := 35
	fibonacci(n)
}

func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

// memoryWork는 슬라이스 할당으로 힙 메모리 부하를 생성한다.
func memoryWork() {
	data := make([]byte, 1024*1024) // 1MB 할당
	for i := range data {
		data[i] = byte(i % 256)
	}
	// math.Sqrt를 호출하여 data가 최적화로 제거되지 않도록 함
	_ = math.Sqrt(float64(len(data)))
	time.Sleep(100 * time.Millisecond)
}

// mutexWork는 여러 고루틴이 동일 뮤텍스를 경합하여 뮤텍스 프로파일 데이터를 생성한다.
func mutexWork() {
	var mu sync.Mutex
	var wg sync.WaitGroup
	counter := 0

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				mu.Lock()
				counter++
				time.Sleep(time.Microsecond)
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
}

func hostname() string {
	h, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return h
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
