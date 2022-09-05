package counter

import (
	"fmt"
	"sync"
)

type CounterMutex struct {
	Num   int64
	Mutex sync.Mutex // 공유데이터 i를 보호하기 위한 뮤텍스
}

// CounterMutex 값을 1씩 증가시킴
func (c *CounterMutex) Increment() {
	c.Mutex.Lock()   // Num 값을 변경하는 부분(임계영역)을 뮤텍스로 잠금
	c.Num += 1       // 공유데이터 변경
	c.Mutex.Unlock() // Num 값 변경 완료 후 뮤텍스 잠금 해제
}

// counter의 값을 출력
func (c *CounterMutex) Display() {
	fmt.Println(c.Num)
}
