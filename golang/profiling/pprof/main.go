// net/http/pprof를 이용한 메모리 누수 프로파일링 예제
//
// 실행 후 브라우저에서 http://localhost:6060/debug/pprof/ 접속하여
// 힙 메모리 프로파일을 확인할 수 있다.
//
// 프로파일 수집:
//
//	go tool pprof http://localhost:6060/debug/pprof/heap
package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	_ "net/http/pprof" // 임포트만으로 pprof HTTP 핸들러가 DefaultServeMux에 자동 등록된다
)

func main() {
	// pprof 엔드포인트를 제공하는 HTTP 서버를 백그라운드에서 시작
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	fmt.Println("hello world")
	var wg sync.WaitGroup
	wg.Add(1)
	go leakyFunction(wg)
	wg.Wait()
}

// leakyFunction은 슬라이스에 문자열을 계속 추가하여 메모리 누수를 시뮬레이션한다.
// append()가 반복될 때마다 슬라이스의 내부 배열이 재할당되면서 메모리 사용량이 지속적으로 증가한다.
// pprof heap 프로파일에서 이 함수의 메모리 할당을 확인할 수 있다.
func leakyFunction(wg sync.WaitGroup) {
	defer wg.Done()
	s := make([]string, 3)
	for i := 0; i < 10000000; i++ {
		s = append(s, "magical pandas") // 슬라이스가 무한히 커지며 메모리 사용량 증가
		if (i % 100000) == 0 {
			time.Sleep(500 * time.Millisecond) // 프로파일링 시간 확보를 위해 의도적 지연
		}
	}
}
