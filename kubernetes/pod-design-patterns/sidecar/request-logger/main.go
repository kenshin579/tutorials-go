package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"
)

func main() {
	targetURL := os.Getenv("TARGET_URL")
	if targetURL == "" {
		targetURL = "http://localhost:3000"
	}
	logFile := os.Getenv("LOG_FILE")
	if logFile == "" {
		logFile = "/var/log/app/access.log"
	}

	target, err := url.Parse(targetURL)
	if err != nil {
		log.Fatalf("invalid TARGET_URL: %v", err)
	}
	proxy := httputil.NewSingleHostReverseProxy(target)

	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}
	defer f.Close()
	fileLogger := log.New(f, "", 0)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		proxy.ServeHTTP(w, r)
		duration := time.Since(start)
		fileLogger.Printf("%s %s %s %v", time.Now().Format(time.RFC3339), r.Method, r.URL.Path, duration)
	})

	log.Println("request-logger sidecar listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
