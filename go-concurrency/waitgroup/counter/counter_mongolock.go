package counter

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	lock "github.com/square/mongo-lock"
)

type CounterMongolock struct {
	Num    int64
	Locker *lock.Client
}

func (c *CounterMongolock) Increment() {
	ctx := context.Background()
	lockID := uuid.Must(uuid.NewRandom()).String()
	err := c.Locker.XLock(ctx, "counter_lock_mongo", lockID, lock.LockDetails{})
	if err != nil {
		fmt.Printf("err:%+v\n", err)
	}

	c.Num += 1 // 공유데이터 변경
	_, err = c.Locker.Unlock(ctx, lockID)
	if err != nil {
		fmt.Printf("err:%+v\n", err)
	}
}

func (c *CounterMongolock) GetNum() int64 {
	return c.Num
}
