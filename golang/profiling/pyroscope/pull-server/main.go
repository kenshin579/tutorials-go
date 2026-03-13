// Pull 모드 예제: pprof 엔드포인트만 노출하는 HTTP 서버
//
// Pyroscope SDK 없이 net/http/pprof만 사용한다.
// Grafana Alloy가 /debug/pprof/* 엔드포인트를 주기적으로 스크래핑하여
// Pyroscope 서버로 프로파일 데이터를 전송한다.
//
// 실행:
//
//	go run .
//
// pprof 엔드포인트 확인:
//
//	curl http://localhost:6060/debug/pprof/
package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	_ "net/http/pprof" // /debug/pprof/* 엔드포인트 자동 등록
	"os"
	"sync"
	"time"
)

func main() {
	port := getEnv("PORT", "6060")

	// 비즈니스 핸들러
	http.HandleFunc("/fast", handleFast)
	http.HandleFunc("/slow", handleSlow)
	http.HandleFunc("/memory", handleMemory)

	log.Printf("Pull 모드 서버 시작: http://localhost:%s", port)
	log.Printf("pprof 엔드포인트: http://localhost:%s/debug/pprof/", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleFast(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `{"message":"fast response"}`)
}

func handleSlow(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	fibonacci(38) // CPU 집약적 연산
	fmt.Fprintf(w, `{"message":"slow response","elapsed":"%s"}`, time.Since(start))
}

func handleMemory(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	allocateMemory() // 대량 메모리 할당
	fmt.Fprintf(w, `{"message":"memory response","elapsed":"%s"}`, time.Since(start))
}

func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

func allocateMemory() {
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			data := make([]byte, 5*1024*1024) // 5MB
			for j := range data {
				data[j] = byte(j % 256)
			}
			_ = math.Sqrt(float64(len(data)))
			time.Sleep(50 * time.Millisecond)
		}()
	}
	wg.Wait()
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
