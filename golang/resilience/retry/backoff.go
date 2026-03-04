package retry

import (
	"context"
	"time"

	"github.com/cenkalti/backoff/v5"
)

// RetryWithExponentialBackoff retries the operation using exponential backoff.
// It stops retrying when maxElapsed time passes or the context is cancelled.
func RetryWithExponentialBackoff(ctx context.Context, operation func() error, maxElapsed time.Duration) error {
	_, err := backoff.Retry(ctx, func() (struct{}, error) {
		return struct{}{}, operation()
	},
		backoff.WithBackOff(backoff.NewExponentialBackOff()),
		backoff.WithMaxElapsedTime(maxElapsed),
	)
	return err
}
