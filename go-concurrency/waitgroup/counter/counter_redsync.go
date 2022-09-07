package counter

import (
	"fmt"

	"github.com/go-redsync/redsync/v4"
)

type CounterRedSync struct {
	Num   int64
	Mutex *redsync.Mutex
}

func (c *CounterRedSync) Increment() {
	if err := c.Mutex.Lock(); err != nil {
		fmt.Println(err)
	}

	c.Num += 1 // 공유데이터 변경
	if ok, err := c.Mutex.Unlock(); !ok || err != nil {
		fmt.Println(err)
	}
}

func (c *CounterRedSync) Display() {
	fmt.Println(c.Num)
}
