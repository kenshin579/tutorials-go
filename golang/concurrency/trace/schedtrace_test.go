package _trace

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSchedtrace_GODEBUG_출력_파싱 - GODEBUG=schedtrace로 스케줄러 상태 확인
// subprocess로 실행하여 GODEBUG 환경변수의 효과를 검증한다.
//
// schedtrace 출력 형식:
//
//	SCHED 0ms: gomaxprocs=8 idleprocs=6 threads=4 spinningthreads=1
//	  idlethreads=0 runqueue=0 [0 0 0 0 0 0 0 0]
//
// 각 필드:
//   - gomaxprocs: GOMAXPROCS 값 (P의 수)
//   - idleprocs: 유휴 P의 수
//   - threads: OS 스레드(M) 수
//   - spinningthreads: 실행할 goroutine을 찾고 있는 M 수
//   - runqueue: 글로벌 실행 큐의 goroutine 수
//   - []: 각 P의 로컬 실행 큐 goroutine 수
func TestSchedtrace_GODEBUG_출력_파싱(t *testing.T) {
	if os.Getenv("SCHEDTRACE_HELPER") == "1" {
		// helper 프로세스: goroutine을 생성하여 스케줄러 활동 유발
		var wg sync.WaitGroup
		for range 50 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				sum := 0
				for range 100_000 {
					sum++
				}
			}()
		}
		wg.Wait()
		return
	}

	// 메인 테스트: subprocess로 helper 실행
	cmd := exec.Command(os.Args[0],
		"-test.run=TestSchedtrace_GODEBUG_출력_파싱",
		"-test.v",
	)
	cmd.Env = append(os.Environ(),
		"SCHEDTRACE_HELPER=1",
		"GODEBUG=schedtrace=100",
	)

	output, err := cmd.CombinedOutput()
	outputStr := string(output)
	t.Logf("schedtrace 출력:\n%s", outputStr)

	// 프로세스가 정상 종료되어야 한다
	assert.NoError(t, err, "subprocess 실행 실패: %s", outputStr)

	// schedtrace 출력에 핵심 필드가 포함되어야 한다
	assert.Contains(t, outputStr, "SCHED", "schedtrace 출력이 있어야 한다")
	assert.Contains(t, outputStr, "gomaxprocs=", "gomaxprocs 필드가 있어야 한다")
	assert.Contains(t, outputStr, "runqueue=", "runqueue 필드가 있어야 한다")
}

// TestSchedtrace_Scheddetail_상세출력 - scheddetail=1로 P/M/G 상세 정보 확인
//
// scheddetail=1 추가 시 P, M, G 각각의 상태가 출력된다:
//   - P: status(idle/running/syscall), schedtick, syscalltick, 로컬 큐 정보
//   - M: P 바인딩 상태, spinning 여부, blocked 여부
//   - G: status(idle/runnable/running/syscall/waiting), 대기 이유
//
// Goroutine 상태 코드:
//
//	_Gidle(0): 생성 직후, 아직 초기화 안 됨
//	_Grunnable(1): 실행 가능, 큐에서 대기 중
//	_Grunning(2): M에서 실행 중
//	_Gsyscall(3): 시스템 콜 실행 중
//	_Gwaiting(4): 채널, 뮤텍스 등에서 블로킹
func TestSchedtrace_Scheddetail_상세출력(t *testing.T) {
	if os.Getenv("SCHEDDETAIL_HELPER") == "1" {
		var wg sync.WaitGroup
		ch := make(chan struct{})

		// 대기 중인 goroutine 생성 (Gwaiting 상태)
		for range 10 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				<-ch
			}()
		}

		// 실행 중인 goroutine 생성 (Grunning/Grunnable 상태)
		for range 10 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				sum := 0
				for range 500_000 {
					sum++
				}
			}()
		}

		runtime.Gosched() // 스케줄러에게 양보하여 출력 기회 제공
		close(ch)
		wg.Wait()
		return
	}

	cmd := exec.Command(os.Args[0],
		"-test.run=TestSchedtrace_Scheddetail_상세출력",
		"-test.v",
	)
	cmd.Env = append(os.Environ(),
		"SCHEDDETAIL_HELPER=1",
		"GODEBUG=schedtrace=100,scheddetail=1",
	)

	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	// 최대 2000자까지만 로깅 (scheddetail은 출력이 매우 길다)
	logOutput := outputStr
	if len(logOutput) > 2000 {
		logOutput = logOutput[:2000] + "\n... (truncated)"
	}
	t.Logf("scheddetail 출력:\n%s", logOutput)

	assert.NoError(t, err, "subprocess 실행 실패")
	assert.Contains(t, outputStr, "SCHED")

	// scheddetail 출력에는 P, M, G 상세 정보가 포함된다
	hasDetail := strings.Contains(outputStr, "P") &&
		(strings.Contains(outputStr, "M") || strings.Contains(outputStr, "G"))
	assert.True(t, hasDetail, "scheddetail 출력에 P/M/G 정보가 있어야 한다")
}
