package counter

import (
	"context"
	"fmt"
	"time"

	"github.com/bsm/redislock"
)

type CounterRedisLock struct {
	Num    int64
	Locker *redislock.Client
}

// CounterMutex 값을 1씩 증가시킴
func (c *CounterRedisLock) Increment() {
	ctx := context.TODO()
	lock, err := c.Locker.Obtain(ctx, "counter", 100*time.Millisecond, &redislock.Options{
		RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(100*time.Millisecond), 3),
	})
	fmt.Println(err)
	c.Num += 1 // 공유데이터 변경
	//c.Mutex.Unlock() // Num 값 변경 완료 후 뮤텍스 잠금 해제
	lock.Release(ctx)
}

// counter의 값을 출력
func (c *CounterRedisLock) Display() {
	fmt.Println(c.Num)
}
