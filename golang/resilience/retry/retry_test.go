package retry

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	retryx "github.com/avast/retry-go/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var errRetryable = errors.New("retryable error")
var errPermanent = errors.New("permanent error")

func TestRetryWithJitter_Success(t *testing.T) {
	var attempts atomic.Int32

	err := RetryWithJitter(context.Background(), func() error {
		if attempts.Add(1) < 3 {
			return errRetryable
		}
		return nil
	}, 5)

	require.NoError(t, err)
	assert.Equal(t, int32(3), attempts.Load())
}

func TestRetryWithJitter_MaxAttempts(t *testing.T) {
	var attempts atomic.Int32

	err := RetryWithJitter(context.Background(), func() error {
		attempts.Add(1)
		return errRetryable
	}, 3, retryx.Delay(10*time.Millisecond))

	assert.Error(t, err, "should fail after max attempts")
	assert.Equal(t, int32(3), attempts.Load())
}

func TestRetryWithJitter_RetryIf(t *testing.T) {
	var attempts atomic.Int32

	err := RetryWithJitter(context.Background(), func() error {
		attempts.Add(1)
		return errPermanent
	}, 5,
		retryx.Delay(10*time.Millisecond),
		retryx.RetryIf(func(err error) bool {
			return !errors.Is(err, errPermanent)
		}),
	)

	assert.Error(t, err, "should not retry permanent errors")
	assert.Equal(t, int32(1), attempts.Load(), "should stop after first attempt")
}

func TestRetryWithJitter_OnRetry(t *testing.T) {
	var retryCount atomic.Int32

	err := RetryWithJitter(context.Background(), func() error {
		return errRetryable
	}, 3,
		retryx.Delay(10*time.Millisecond),
		retryx.OnRetry(func(n uint, err error) {
			retryCount.Add(1)
		}),
	)

	assert.Error(t, err)
	assert.Equal(t, int32(3), retryCount.Load(), "OnRetry should be called for each failed attempt")
}

func TestRetryWithJitter_ContextCancel(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := RetryWithJitter(ctx, func() error {
		return errRetryable
	}, 100, retryx.Delay(50*time.Millisecond))

	assert.Error(t, err, "should fail when context is cancelled")
}
