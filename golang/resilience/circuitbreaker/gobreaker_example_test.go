package circuitbreaker

import (
	"errors"
	"testing"
	"time"

	"github.com/sony/gobreaker/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var errService = errors.New("service unavailable")

func newTestBreaker() *gobreaker.CircuitBreaker[[]byte] {
	settings := gobreaker.Settings{
		Name:        "test",
		MaxRequests: 1,               // Half-Open에서 1개 요청만 허용
		Interval:    0,               // 카운터 자동 초기화 안 함
		Timeout:     1 * time.Second, // Open → Half-Open 빠르게 전환 (테스트용)
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 3
		},
	}
	return gobreaker.NewCircuitBreaker[[]byte](settings)
}

func TestClosedState_SuccessfulRequests(t *testing.T) {
	cb := newTestBreaker()

	result, err := cb.Execute(func() ([]byte, error) {
		return []byte("ok"), nil
	})

	require.NoError(t, err)
	assert.Equal(t, []byte("ok"), result)
	assert.Equal(t, gobreaker.StateClosed, cb.State())
}

func TestClosedToOpen_ConsecutiveFailures(t *testing.T) {
	cb := newTestBreaker()

	// 3회 연속 실패 → Open 전환
	for i := 0; i < 3; i++ {
		_, _ = cb.Execute(func() ([]byte, error) {
			return nil, errService
		})
	}

	assert.Equal(t, gobreaker.StateOpen, cb.State())

	// Open 상태에서 요청은 즉시 거부
	_, err := cb.Execute(func() ([]byte, error) {
		return []byte("ok"), nil
	})
	assert.ErrorIs(t, err, gobreaker.ErrOpenState)
}

func TestOpenToHalfOpen_AfterTimeout(t *testing.T) {
	cb := newTestBreaker()

	// Open 상태로 전환
	for i := 0; i < 3; i++ {
		_, _ = cb.Execute(func() ([]byte, error) {
			return nil, errService
		})
	}
	assert.Equal(t, gobreaker.StateOpen, cb.State())

	// Timeout(1초) 경과 대기
	time.Sleep(1200 * time.Millisecond)

	assert.Equal(t, gobreaker.StateHalfOpen, cb.State())
}

func TestHalfOpenToClosed_OnSuccess(t *testing.T) {
	cb := newTestBreaker()

	// Open 전환
	for i := 0; i < 3; i++ {
		_, _ = cb.Execute(func() ([]byte, error) {
			return nil, errService
		})
	}

	// Half-Open 대기
	time.Sleep(1200 * time.Millisecond)
	assert.Equal(t, gobreaker.StateHalfOpen, cb.State())

	// Half-Open에서 성공 → Closed 복귀
	_, err := cb.Execute(func() ([]byte, error) {
		return []byte("recovered"), nil
	})
	require.NoError(t, err)
	assert.Equal(t, gobreaker.StateClosed, cb.State())
}

func TestGobreaker_HalfOpenToOpen_OnFailure(t *testing.T) {
	cb := newTestBreaker()

	// Open 전환
	for i := 0; i < 3; i++ {
		_, _ = cb.Execute(func() ([]byte, error) {
			return nil, errService
		})
	}

	// Half-Open 대기
	time.Sleep(1200 * time.Millisecond)
	assert.Equal(t, gobreaker.StateHalfOpen, cb.State())

	// Half-Open에서 실패 → Open 재전환
	_, _ = cb.Execute(func() ([]byte, error) {
		return nil, errService
	})
	assert.Equal(t, gobreaker.StateOpen, cb.State())
}
