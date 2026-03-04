package circuitbreaker

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/failsafe-go/failsafe-go"
	"github.com/failsafe-go/failsafe-go/circuitbreaker"
	"github.com/failsafe-go/failsafe-go/fallback"
	"github.com/failsafe-go/failsafe-go/retrypolicy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComposed_CBOpen_ReturnsFallback(t *testing.T) {
	policy := NewComposedPolicy("fallback-data")

	// CB를 Open 상태로 만들기
	for i := 0; i < 5; i++ {
		_, _ = policy.Execute(func() (string, error) {
			return "", ErrExternal
		})
	}
	assert.True(t, policy.CircuitBreaker.IsOpen())

	// Open 상태에서 Fallback 반환
	result, err := policy.Execute(func() (string, error) {
		return "should-not-reach", nil
	})
	require.NoError(t, err)
	assert.Equal(t, "fallback-data", result)
}

func TestComposed_TransientFailure_RetrySuccess(t *testing.T) {
	var attempts atomic.Int32

	fb := fallback.NewBuilderWithResult("fallback").
		HandleErrors(ErrExternal, circuitbreaker.ErrOpen).
		Build()
	rt := retrypolicy.NewBuilder[string]().
		HandleErrors(ErrExternal).
		WithMaxRetries(3).
		WithDelay(10 * time.Millisecond).
		Build()
	cb := circuitbreaker.NewBuilder[string]().
		HandleErrors(ErrExternal).
		WithFailureThreshold(10). // 높은 임계값으로 CB가 열리지 않도록
		WithDelay(1 * time.Second).
		Build()

	result, err := failsafe.With(fb, rt, cb).Get(func() (string, error) {
		if attempts.Add(1) <= 2 {
			return "", ErrExternal
		}
		return "success-after-retry", nil
	})

	require.NoError(t, err)
	assert.Equal(t, "success-after-retry", result)
	assert.Equal(t, int32(3), attempts.Load()) // 2회 실패 + 1회 성공
}

func TestComposed_RetryExhausted_CBOpen_Fallback(t *testing.T) {
	fb := fallback.NewBuilderWithResult("fallback").
		HandleErrors(ErrExternal, circuitbreaker.ErrOpen).
		Build()
	rt := retrypolicy.NewBuilder[string]().
		HandleErrors(ErrExternal).
		WithMaxRetries(2).
		WithDelay(10 * time.Millisecond).
		ReturnLastFailure().
		Build()
	cb := circuitbreaker.NewBuilder[string]().
		HandleErrors(ErrExternal).
		WithFailureThreshold(3).
		WithDelay(1 * time.Second).
		Build()

	// 계속 실패하여 Retry 소진 + CB Open → Fallback
	result, err := failsafe.With(fb, rt, cb).Get(func() (string, error) {
		return "", ErrExternal
	})

	require.NoError(t, err)
	assert.Equal(t, "fallback", result)
	assert.True(t, cb.IsOpen())
}
