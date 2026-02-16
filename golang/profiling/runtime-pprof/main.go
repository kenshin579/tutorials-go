// runtime/pprof를 사용한 파일 기반 프로파일링 예제
//
// net/http/pprof와 달리 HTTP 서버 없이 프로파일 데이터를 파일로 직접 저장한다.
// CLI 프로그램이나 배치 작업처럼 HTTP 서버가 없는 환경에서 사용한다.
//
// 실행 방법:
//
//	go run main.go
//	go tool pprof cpu.prof       # CPU 프로파일 분석
//	go tool pprof mem.prof       # 힙 프로파일 분석
//	go tool pprof -http=:8080 cpu.prof  # 웹 UI로 분석
package main

import (
	"fmt"
	"log"
	"os"
	"runtime/pprof"
)

func main() {
	// CPU 프로파일 파일 생성
	cpuFile, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal(err)
	}
	defer cpuFile.Close()

	// CPU 프로파일링 시작 - StartCPUProfile과 StopCPUProfile 사이의 코드가 프로파일링된다
	if err := pprof.StartCPUProfile(cpuFile); err != nil {
		log.Fatal(err)
	}
	defer pprof.StopCPUProfile()

	// CPU 부하 생성 - 프로파일에서 이 함수의 CPU 사용량을 확인할 수 있다
	result := heavyComputation()
	fmt.Printf("computation result: %d\n", result)

	// 메모리 할당 부하 생성 - 힙 프로파일에서 할당 출처를 확인할 수 있다
	data := allocateMemory()
	fmt.Printf("allocated %d items\n", len(data))

	// 힙 프로파일 저장 - 현재 시점의 힙 메모리 스냅샷을 파일로 기록한다
	writeHeapProfile()

	fmt.Println("프로파일 파일 생성 완료:")
	fmt.Println("  - cpu.prof  (go tool pprof cpu.prof)")
	fmt.Println("  - mem.prof  (go tool pprof mem.prof)")
}

// heavyComputation은 1억 회 반복 연산으로 CPU 부하를 생성한다.
// pprof CPU 프로파일에서 flat 시간이 높게 나타나는 함수이다.
func heavyComputation() int {
	result := 0
	for i := 0; i < 100_000_000; i++ {
		result += i * i
	}
	return result
}

// allocateMemory는 10KB 버퍼 1000개(총 ~10MB)를 할당하여 힙 메모리 부하를 생성한다.
// pprof heap 프로파일의 inuse_space에서 이 함수의 할당을 확인할 수 있다.
func allocateMemory() [][]byte {
	var data [][]byte
	for i := 0; i < 1000; i++ {
		buf := make([]byte, 10_000)
		for j := range buf {
			buf[j] = byte(j % 256)
		}
		data = append(data, buf)
	}
	return data
}

// writeHeapProfile은 현재 시점의 힙 메모리 상태를 파일로 저장한다.
// WriteHeapProfile은 runtime.GC()를 호출하지 않으므로, 정확한 측정이 필요하면
// 호출 전에 runtime.GC()를 수동으로 실행하는 것이 좋다.
func writeHeapProfile() {
	f, err := os.Create("mem.prof")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if err := pprof.WriteHeapProfile(f); err != nil {
		log.Fatal(err)
	}
}
