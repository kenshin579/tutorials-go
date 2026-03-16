package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"sync/atomic"
	"time"
)

type StatsResponse struct {
	RequestCount int64   `json:"request_count"`
	Uptime       float64 `json:"uptime_seconds"`
	GoRoutines   int     `json:"goroutines"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	start := time.Now()
	var requestCount atomic.Int64

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		count := requestCount.Add(1)
		fmt.Fprintf(w, "Hello from main-app! (request #%d)\n", count)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})

	// Adapter 패턴용: 커스텀 JSON 메트릭
	http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		stats := StatsResponse{
			RequestCount: requestCount.Load(),
			Uptime:       time.Since(start).Seconds(),
			GoRoutines:   runtime.NumGoroutine(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)
	})

	log.Printf("main-app listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
