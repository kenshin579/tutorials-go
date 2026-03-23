package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// Server는 HTTP API 서버 설정을 담는 구조체이다.
type Server struct {
	Port   string
	Router *http.ServeMux
}

// Response는 API 응답 구조체이다.
type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Time    string `json:"time"`
}

// HealthResponse는 헬스체크 응답 구조체이다.
type HealthResponse struct {
	Status string `json:"status"`
	Uptime string `json:"uptime"`
}

var startTime = time.Now()

// NewServer는 새로운 Server 인스턴스를 생성한다.
func NewServer(port string) *Server {
	s := &Server{
		Port:   port,
		Router: http.NewServeMux(),
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.Router.HandleFunc("GET /api/hello", handleHello)
	s.Router.HandleFunc("GET /health", handleHealth)
}

func handleHello(w http.ResponseWriter, r *http.Request) {
	resp := Response{
		Status:  "ok",
		Message: "Hello from Ralph Loop demo!",
		Time:    time.Now().Format(time.RFC3339),
	}
	writeJSON(w, http.StatusOK, resp)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	resp := HealthResponse{
		Status: "healthy",
		Uptime: time.Since(startTime).Round(time.Second).String(),
	}
	writeJSON(w, http.StatusOK, resp)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("JSON 인코딩 실패: %v", err)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := NewServer(port)
	fmt.Printf("서버 시작: http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, srv.Router))
}
