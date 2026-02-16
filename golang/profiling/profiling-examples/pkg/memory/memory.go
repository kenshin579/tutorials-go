// memory 패키지는 힙 메모리 프로파일링을 위한 메모리 할당 함수를 제공한다.
// pprof heap 프로파일에서 alloc1000 함수의 할당을 확인할 수 있다.
package memory

import (
	"time"
)

// AllocMemory는 1000바이트를 힙에 할당하고 프로그램이 종료되지 않도록 대기한다.
// 할당된 메모리가 해제되지 않으므로 pprof inuse_space에서 확인할 수 있다.
func AllocMemory() {
	bytes1000 := alloc1000()
	bytes1000[0] = '0' // 할당된 메모리를 실제로 사용하여 컴파일러 최적화로 제거되지 않도록 한다

	for {
		time.Sleep(1 * time.Second) // 프로세스를 유지하여 프로파일링 가능하도록 대기
	}
}

// alloc1000은 1000바이트 슬라이스를 힙에 할당하여 반환한다.
// pprof heap 프로파일에서 이 함수가 할당 출처(allocation site)로 표시된다.
func alloc1000() []byte {
	return make([]byte, 1000)
}
