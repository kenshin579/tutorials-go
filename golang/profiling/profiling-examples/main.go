package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/google/gops/agent"
	"github.com/kenshin579/tutorials-go/go-profiling/profiling-examples/pkg/block"
	"github.com/kenshin579/tutorials-go/go-profiling/profiling-examples/pkg/cpu"
	"github.com/kenshin579/tutorials-go/go-profiling/profiling-examples/pkg/memory"
	"github.com/kenshin579/tutorials-go/go-profiling/profiling-examples/pkg/mutex"
	"github.com/kenshin579/tutorials-go/go-profiling/profiling-examples/pkg/threadcreate"
)

// 종합 프로파일링 예제 - CPU, 메모리, 블로킹, 뮤텍스, 스레드 생성 프로파일을 동시에 수집할 수 있다.
//
// 실행 후 아래 명령어로 각 프로파일을 수집할 수 있다:
//
//	go tool pprof http://localhost:6060/debug/pprof/profile?seconds=10  (CPU)
//	go tool pprof http://localhost:6060/debug/pprof/heap                (힙 메모리)
//	go tool pprof http://localhost:6060/debug/pprof/goroutine           (고루틴)
//	go tool pprof http://localhost:6060/debug/pprof/block               (블로킹)
//	go tool pprof http://localhost:6060/debug/pprof/mutex               (뮤텍스 경합)
//	go tool pprof http://localhost:6060/debug/pprof/threadcreate        (스레드 생성)
//
// 참고: https://github.com/ssup2/golang-profiling-example
func main() {
	// pprof HTTP 엔드포인트를 제공하는 서버 시작 (localhost:6060)
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// gops 에이전트 시작 - 실행 중인 Go 프로세스 정보를 조회할 수 있다
	go func() {
		agent.Listen(agent.Options{})
	}()

	// 블로킹/뮤텍스 프로파일은 기본 비활성화이므로 명시적으로 활성화해야 한다
	runtime.SetBlockProfileRate(1)      // 모든 블로킹 이벤트 기록 (1 = 나노초 임계값)
	runtime.SetMutexProfileFraction(1)  // 모든 뮤텍스 경합 이벤트 기록 (1 = 1/1 확률)

	// 각 프로파일 유형별 부하를 생성하는 고루틴 시작
	go cpu.IncreaseInt()                    // CPU 부하 (무한 루프 연산)
	go cpu.IncreaseIntGoroutine()           // CPU 부하 (중첩 고루틴)
	go memory.AllocMemory()                 // 힙 메모리 할당
	go block.PrintHello()                   // stdout 블로킹 (I/O Lock 경합)
	go block.PrintWorld()                   // stdout 블로킹 (I/O Lock 경합)
	go threadcreate.CreateGoroutine1000()   // 대량 고루틴 생성 → OS 스레드 생성 유발
	go mutex.Mutex01()                      // 뮤텍스 경합
	go mutex.Mutex02()                      // 뮤텍스 경합
	go mutex.Mutex03()                      // 뮤텍스 경합

	// SIGTERM 또는 SIGINT 시그널을 받을 때까지 대기
	log.Println("프로파일링 서버 시작: http://localhost:6060/debug/pprof/")
	termSignal := make(chan os.Signal, 1)
	signal.Notify(termSignal, syscall.SIGTERM, syscall.SIGINT)
	<-termSignal
}
