package retry

import (
	"context"
	"time"

	retryx "github.com/avast/retry-go/v4"
)

// RetryWithJitter retries fn with exponential backoff and jitter.
func RetryWithJitter(ctx context.Context, fn retryx.RetryableFunc, maxAttempts uint, opts ...retryx.Option) error {
	defaultOpts := []retryx.Option{
		retryx.Attempts(maxAttempts),
		retryx.DelayType(retryx.BackOffDelay),
		retryx.Context(ctx),
		retryx.MaxJitter(1 * time.Second),
	}
	return retryx.Do(fn, append(defaultOpts, opts...)...)
}
