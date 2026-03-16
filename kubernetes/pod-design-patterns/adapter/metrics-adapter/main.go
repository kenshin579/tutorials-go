package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type AppStats struct {
	RequestCount int64   `json:"request_count"`
	Uptime       float64 `json:"uptime_seconds"`
	GoRoutines   int     `json:"goroutines"`
}

func main() {
	appURL := os.Getenv("APP_STATS_URL")
	if appURL == "" {
		appURL = "http://localhost:3000/stats"
	}

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get(appURL)
		if err != nil {
			http.Error(w, "failed to fetch stats", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		var stats AppStats
		if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
			http.Error(w, "failed to decode stats", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, "# HELP app_request_count Total number of requests\n")
		fmt.Fprintf(w, "# TYPE app_request_count counter\n")
		fmt.Fprintf(w, "app_request_count %d\n\n", stats.RequestCount)
		fmt.Fprintf(w, "# HELP app_uptime_seconds Application uptime in seconds\n")
		fmt.Fprintf(w, "# TYPE app_uptime_seconds gauge\n")
		fmt.Fprintf(w, "app_uptime_seconds %.2f\n\n", stats.Uptime)
		fmt.Fprintf(w, "# HELP app_goroutines Number of goroutines\n")
		fmt.Fprintf(w, "# TYPE app_goroutines gauge\n")
		fmt.Fprintf(w, "app_goroutines %d\n", stats.GoRoutines)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	log.Println("metrics-adapter listening on :9090")
	log.Fatal(http.ListenAndServe(":9090", nil))
}
