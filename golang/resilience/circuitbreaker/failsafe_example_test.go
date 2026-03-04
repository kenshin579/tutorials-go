package circuitbreaker

import (
	"testing"
	"time"

	"github.com/failsafe-go/failsafe-go"
	"github.com/failsafe-go/failsafe-go/circuitbreaker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCountBased_FailureThreshold(t *testing.T) {
	cb := circuitbreaker.NewBuilder[string]().
		HandleErrors(ErrExternal).
		WithFailureThreshold(3).
		WithDelay(1 * time.Second).
		Build()

	// 3회 실패 → Open
	for i := 0; i < 3; i++ {
		_, _ = failsafe.With(cb).Get(func() (string, error) {
			return "", ErrExternal
		})
	}

	assert.True(t, cb.IsOpen())

	// Open 상태에서 요청 거부
	_, err := failsafe.With(cb).Get(func() (string, error) {
		return "ok", nil
	})
	assert.Error(t, err)
}

func TestCountBased_FullCycle(t *testing.T) {
	cb := circuitbreaker.NewBuilder[string]().
		HandleErrors(ErrExternal).
		WithFailureThreshold(3).
		WithSuccessThreshold(1).
		WithDelay(1 * time.Second).
		Build()

	assert.True(t, cb.IsClosed())

	// Closed → Open (3회 실패)
	for i := 0; i < 3; i++ {
		_, _ = failsafe.With(cb).Get(func() (string, error) {
			return "", ErrExternal
		})
	}
	assert.True(t, cb.IsOpen())

	// Open → Half-Open (delay 경과 후 다음 요청 시 전환)
	time.Sleep(1200 * time.Millisecond)

	// Half-Open에서 성공 → Closed 복귀
	result, err := failsafe.With(cb).Get(func() (string, error) {
		return "recovered", nil
	})
	require.NoError(t, err)
	assert.Equal(t, "recovered", result)
	assert.True(t, cb.IsClosed())
}

func TestTimeBased_FailureThresholdPeriod(t *testing.T) {
	cb := circuitbreaker.NewBuilder[string]().
		HandleErrors(ErrExternal).
		WithFailureThresholdPeriod(3, 1*time.Minute).
		WithDelay(1 * time.Second).
		Build()

	assert.True(t, cb.IsClosed())

	// 1분 윈도우 내 3회 실패 → Open
	for i := 0; i < 3; i++ {
		_, _ = failsafe.With(cb).Get(func() (string, error) {
			return "", ErrExternal
		})
	}
	assert.True(t, cb.IsOpen())
}

func TestHalfOpenToOpen_OnFailure(t *testing.T) {
	cb := circuitbreaker.NewBuilder[string]().
		HandleErrors(ErrExternal).
		WithFailureThreshold(3).
		WithSuccessThreshold(1).
		WithDelay(1 * time.Second).
		Build()

	// Open 전환
	for i := 0; i < 3; i++ {
		_, _ = failsafe.With(cb).Get(func() (string, error) {
			return "", ErrExternal
		})
	}
	assert.True(t, cb.IsOpen())

	// delay 경과 대기 → 다음 요청 시 Half-Open 전환
	time.Sleep(1200 * time.Millisecond)

	// Half-Open에서 실패 → Open 재전환
	_, _ = failsafe.With(cb).Get(func() (string, error) {
		return "", ErrExternal
	})
	assert.True(t, cb.IsOpen())
}
