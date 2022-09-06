package counter

import (
	"context"
	"fmt"
	"time"

	"github.com/bsm/redislock"
)

type CounterRedSync struct {
	Num int64
}

func (c *CounterRedSync) Increment() {
	ctx := context.TODO()
	lock, err := c.Locker.Obtain(ctx, "counter", 100*time.Millisecond, &redislock.Options{
		RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(100*time.Millisecond), 3),
	})
	fmt.Println(err)
	c.Num += 1 // 공유데이터 변경
	//c.Mutex.Unlock() // Num 값 변경 완료 후 뮤텍스 잠금 해제
	lock.Release(ctx)
}

func (c *CounterRedSync) Display() {
	fmt.Println(c.Num)
}
