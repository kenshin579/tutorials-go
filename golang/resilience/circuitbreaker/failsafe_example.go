package circuitbreaker

import (
	"errors"
	"log"
	"time"

	"github.com/failsafe-go/failsafe-go/circuitbreaker"
)

var ErrExternal = errors.New("external service error")

// NewCountBasedBreaker는 실패 횟수 기반 Circuit Breaker를 생성한다.
// - 5회 연속 실패 시 Open 전환
// - 30초 후 Half-Open 전환
// - Half-Open에서 2회 연속 성공 시 Closed 복귀
func NewCountBasedBreaker() circuitbreaker.CircuitBreaker[any] {
	return circuitbreaker.NewBuilder[any]().
		HandleErrors(ErrExternal).
		WithFailureThreshold(5).
		WithSuccessThreshold(2).
		WithDelay(30 * time.Second).
		OnStateChanged(func(e circuitbreaker.StateChangedEvent) {
			log.Printf("Count-based CB: %s → %s", e.OldState, e.NewState)
		}).
		Build()
}

// NewTimeBasedBreaker는 시간 윈도우 기반 Circuit Breaker를 생성한다.
// - 1분 동안 3회 이상 실패 시 Open 전환
// - 30초 후 Half-Open 전환
// - Half-Open에서 2회 연속 성공 시 Closed 복귀
func NewTimeBasedBreaker() circuitbreaker.CircuitBreaker[any] {
	return circuitbreaker.NewBuilder[any]().
		HandleErrors(ErrExternal).
		WithFailureThresholdPeriod(3, 1*time.Minute).
		WithSuccessThreshold(2).
		WithDelay(30 * time.Second).
		OnStateChanged(func(e circuitbreaker.StateChangedEvent) {
			log.Printf("Time-based CB: %s → %s", e.OldState, e.NewState)
		}).
		Build()
}
