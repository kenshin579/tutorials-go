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
	lock, _ := c.Locker.Obtain(ctx, "counter", 30*time.Second, &redislock.Options{
		RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(100*time.Millisecond), 100),
	})

	//if err != nil {
	//	//fmt.Printf("err:%+v\n", err)
	//}
	c.Num += 1 // 공유데이터 변경

	if lock != nil {
		defer lock.Release(ctx)
	}

}

func (c *CounterRedisLock) GetNum() int64 {
	return c.Num
}

// counter의 값을 출력
func (c *CounterRedisLock) Display() {
	fmt.Println(c.Num)
}
