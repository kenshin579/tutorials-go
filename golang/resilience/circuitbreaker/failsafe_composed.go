package circuitbreaker

import (
	"time"

	"github.com/failsafe-go/failsafe-go"
	"github.com/failsafe-go/failsafe-go/circuitbreaker"
	"github.com/failsafe-go/failsafe-go/fallback"
	"github.com/failsafe-go/failsafe-go/retrypolicy"
)

// ComposedPolicy는 Fallback + Retry + Circuit Breaker를 조합한 정책이다.
type ComposedPolicy struct {
	Fallback       fallback.Fallback[string]
	RetryPolicy    retrypolicy.RetryPolicy[string]
	CircuitBreaker circuitbreaker.CircuitBreaker[string]
}

// NewComposedPolicy는 기본 설정으로 조합 정책을 생성한다.
//
// 실행 흐름: Fallback → Retry → CircuitBreaker → fn
//   - CircuitBreaker가 먼저 요청 허용 여부를 결정
//   - 실패 시 Retry가 재시도
//   - Retry 소진 시 Fallback이 대체값 반환
func NewComposedPolicy(fallbackValue string) *ComposedPolicy {
	fb := fallback.NewBuilderWithResult(fallbackValue).
		HandleErrors(ErrExternal, circuitbreaker.ErrOpen).
		Build()

	rt := retrypolicy.NewBuilder[string]().
		HandleErrors(ErrExternal).
		WithMaxRetries(3).
		WithBackoff(100*time.Millisecond, 1*time.Second).
		Build()

	cb := circuitbreaker.NewBuilder[string]().
		HandleErrors(ErrExternal).
		WithFailureThreshold(5).
		WithSuccessThreshold(2).
		WithDelay(1 * time.Second).
		Build()

	return &ComposedPolicy{
		Fallback:       fb,
		RetryPolicy:    rt,
		CircuitBreaker: cb,
	}
}

// Execute는 조합된 정책으로 함수를 실행한다.
func (p *ComposedPolicy) Execute(fn func() (string, error)) (string, error) {
	return failsafe.With(p.Fallback, p.RetryPolicy, p.CircuitBreaker).Get(fn)
}
