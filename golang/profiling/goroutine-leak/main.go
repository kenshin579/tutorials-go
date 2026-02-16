package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	_ "net/http/pprof"
)

func main() {
	// pprof HTTP 서버 시작
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	fmt.Printf("초기 고루틴 수: %d\n", runtime.NumGoroutine())

	// 고루틴 누수 시연: 닫히지 않는 채널 대기
	for i := 0; i < 100; i++ {
		go leakyGoroutine(i)
	}

	time.Sleep(1 * time.Second)
	fmt.Printf("누수 후 고루틴 수: %d\n", runtime.NumGoroutine())

	// 안전한 고루틴 사용: context로 취소 가능
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	for i := 0; i < 50; i++ {
		go safeGoroutine(ctx, i)
	}

	time.Sleep(1 * time.Second)
	fmt.Printf("안전한 고루틴 추가 후 고루틴 수: %d\n", runtime.NumGoroutine())

	// context 타임아웃 대기
	<-ctx.Done()
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("context 취소 후 고루틴 수: %d\n", runtime.NumGoroutine())

	fmt.Println("\npprof 확인: http://localhost:6060/debug/pprof/goroutine?debug=1")
	fmt.Println("누수된 100개의 고루틴이 여전히 남아있는 것을 확인할 수 있습니다.")

	// 프로파일링 확인을 위해 대기
	select {}
}

// leakyGoroutine은 아무도 닫지 않는 채널을 대기하여 고루틴 누수를 발생시킨다.
func leakyGoroutine(id int) {
	ch := make(chan struct{})
	<-ch // 영원히 대기 -> 고루틴 누수!
	fmt.Println("never reached", id)
}

// safeGoroutine은 context를 사용하여 취소 가능한 고루틴을 구현한다.
func safeGoroutine(ctx context.Context, id int) {
	ch := make(chan struct{})
	select {
	case <-ch:
		fmt.Println("received", id)
	case <-ctx.Done():
		return // context 취소 시 정상 종료
	}
}
