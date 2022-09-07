package counter

import (
	"fmt"

	lock "github.com/square/mongo-lock"
)

type CounterMongolock struct {
	Num        int64
	lockClient *lock.Client
}

// CounterMutex 값을 1씩 증가시킴
func (c *CounterMongolock) Increment() {
	c.Num += 1 // 공유데이터 변경
}

// counter의 값을 출력
func (c *CounterMongolock) Display() {
	fmt.Println(c.Num)
}
