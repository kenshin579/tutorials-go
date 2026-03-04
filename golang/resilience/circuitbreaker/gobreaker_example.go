package circuitbreaker

import (
	"log"
	"time"

	"github.com/sony/gobreaker/v2"
)

// NewBreaker는 기본 설정으로 Circuit Breaker를 생성한다.
// - 연속 5회 실패 시 Open 전환
// - Open 상태에서 30초 후 Half-Open 전환
// - Half-Open에서 최대 3개 요청 허용
func NewBreaker(name string) *gobreaker.CircuitBreaker[[]byte] {
	settings := gobreaker.Settings{
		Name:        name,
		MaxRequests: 3,                // Half-Open 상태에서 허용할 요청 수
		Interval:    10 * time.Second, // Closed 상태 카운터 초기화 주기
		Timeout:     30 * time.Second, // Open → Half-Open 전환 대기 시간
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > 5
		},
		OnStateChange: func(name string, from, to gobreaker.State) {
			log.Printf("Circuit Breaker %s: %s → %s", name, from, to)
		},
	}
	return gobreaker.NewCircuitBreaker[[]byte](settings)
}

// Execute는 Circuit Breaker로 보호된 함수를 실행한다.
func Execute(cb *gobreaker.CircuitBreaker[[]byte], fn func() ([]byte, error)) ([]byte, error) {
	return cb.Execute(fn)
}
