package account

import (
	"context"
	"time"

	"github.com/bsm/redislock"
)

type AccountRedislock struct {
	locker          *redislock.Client
	customerBalance map[string]int
}

func (a *AccountRedislock) add(customerName string) {
	ctx := context.TODO()
	lock, _ := a.locker.Obtain(ctx, "account_lock_redislock", 5*time.Second, &redislock.Options{
		//RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(100*time.Millisecond), 100),
		RetryStrategy: redislock.LimitRetry(redislock.ExponentialBackoff(10*time.Millisecond, 300*time.Millisecond), 100),
	})
	if lock != nil {
		defer lock.Release(ctx)
	}

	a.customerBalance[customerName]++
}

func (a *AccountRedislock) GetBalance(customerName string) int {
	return a.customerBalance[customerName]
}
