package main

import (
	"context"
	"fmt"
	"os"
	"runtime/trace"
	"sync"
)

// go tool trace 예제
//
// 실행 방법:
//   go run main.go
//   go tool trace trace.out

func main() {
	f, err := os.Create("trace.out")
	if err != nil {
		fmt.Fprintf(os.Stderr, "트레이스 파일 생성 실패: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	if err := trace.Start(f); err != nil {
		fmt.Fprintf(os.Stderr, "트레이스 시작 실패: %v\n", err)
		os.Exit(1)
	}
	defer trace.Stop()

	ctx := context.Background()

	// 여러 고루틴에서 동시 작업 수행
	var wg sync.WaitGroup
	for i := range 5 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(ctx, i)
		}()
	}
	wg.Wait()

	fmt.Println("트레이스 수집 완료: trace.out")
	fmt.Println("분석: go tool trace trace.out")
}

func worker(ctx context.Context, id int) {
	// 사용자 정의 태스크 영역 생성
	ctx, task := trace.NewTask(ctx, fmt.Sprintf("worker-%d", id))
	defer task.End()

	// CPU 작업
	trace.WithRegion(ctx, "compute", func() {
		sum := 0
		for range 10_000_000 {
			sum++
		}
	})

	// 채널 통신
	trace.WithRegion(ctx, "channel-work", func() {
		ch := make(chan int, 10)
		go func() {
			for i := range 10 {
				ch <- i
			}
			close(ch)
		}()
		for v := range ch {
			_ = v
		}
	})

	trace.Log(ctx, "status", fmt.Sprintf("worker-%d completed", id))
}
