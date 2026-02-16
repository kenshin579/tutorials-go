// mutex 패키지는 뮤텍스 경합(contention) 프로파일링 예제를 제공한다.
// 3개의 고루틴이 동일한 뮤텍스를 두고 경쟁하며, pprof mutex 프로파일에서
// 각 함수가 Lock 획득을 위해 대기한 시간을 확인할 수 있다.
//
// 뮤텍스 프로파일을 수집하려면 프로그램 시작 시 아래 설정이 필요하다:
//
//	runtime.SetMutexProfileFraction(1)
package mutex

import (
	"fmt"
	"sync"
)

// 3개의 고루틴이 공유하는 뮤텍스. 동시 접근 시 경합이 발생한다.
var mutex = sync.Mutex{}

// Mutex01은 공유 뮤텍스를 반복적으로 획득/해제하며 경합을 발생시킨다.
// Mutex02, Mutex03과 동시에 실행되면 Lock 대기 시간이 프로파일에 기록된다.
func Mutex01() {
	for {
		mutex.Lock()
		fmt.Printf("Mutex01\n")
		mutex.Unlock()
	}
}

// Mutex02는 Mutex01과 동일한 뮤텍스를 사용하여 경합을 발생시킨다.
func Mutex02() {
	for {
		mutex.Lock()
		fmt.Printf("Mutex02\n")
		mutex.Unlock()
	}
}

// Mutex03은 Mutex01, Mutex02와 동일한 뮤텍스를 사용하여 3-way 경합을 발생시킨다.
func Mutex03() {
	for {
		mutex.Lock()
		fmt.Printf("Mutex03\n")
		mutex.Unlock()
	}
}
