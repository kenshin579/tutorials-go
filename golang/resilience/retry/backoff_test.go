package retry

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var errTransient = errors.New("transient error")

func TestRetryWithExponentialBackoff_Success(t *testing.T) {
	var attempts atomic.Int32

	err := RetryWithExponentialBackoff(context.Background(), func() error {
		if attempts.Add(1) < 3 {
			return errTransient
		}
		return nil
	}, 5*time.Second)

	require.NoError(t, err)
	assert.Equal(t, int32(3), attempts.Load())
}

func TestRetryWithExponentialBackoff_MaxElapsed(t *testing.T) {
	err := RetryWithExponentialBackoff(context.Background(), func() error {
		return errTransient
	}, 500*time.Millisecond)

	assert.Error(t, err, "should fail after max elapsed time")
}

func TestRetryWithExponentialBackoff_ContextCancel(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	err := RetryWithExponentialBackoff(ctx, func() error {
		return errTransient
	}, 10*time.Second)

	assert.Error(t, err, "should fail when context is cancelled")
}
