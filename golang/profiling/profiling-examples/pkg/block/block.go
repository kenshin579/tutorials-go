// block 패키지는 블로킹 프로파일링을 위한 I/O 경합 예제를 제공한다.
// PrintHello와 PrintWorld가 동시에 실행되면 stdout에 대한 내부 Lock을 두고 경합이 발생한다.
// pprof block 프로파일에서 fmt.Printf 내부의 대기 시간을 확인할 수 있다.
//
// 블로킹 프로파일을 수집하려면 프로그램 시작 시 아래 설정이 필요하다:
//
//	runtime.SetBlockProfileRate(1)
package block

import (
	"fmt"
)

// PrintHello는 무한 루프에서 "Hello"를 출력한다.
// fmt.Printf는 내부적으로 stdout Lock을 획득하므로 PrintWorld와 동시에 실행 시 블로킹이 발생한다.
func PrintHello() {
	for {
		fmt.Printf("Hello\n")
	}
}

// PrintWorld는 무한 루프에서 "World"를 출력한다.
// PrintHello와 stdout Lock을 두고 경합하여 블로킹 프로파일에 기록된다.
func PrintWorld() {
	for {
		fmt.Printf("World\n")
	}
}
