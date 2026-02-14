package shutdown

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

// Server - Graceful Shutdown을 지원하는 HTTP 서버
type Server struct {
	httpServer *http.Server
	mux        *http.ServeMux
}

// NewServer - 새 서버 생성
func NewServer(addr string) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, World!")
	})

	mux.HandleFunc("/slow", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second) // 느린 요청 시뮬레이션
		fmt.Fprintln(w, "Slow response done")
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})

	return &Server{
		httpServer: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
		mux: mux,
	}
}

// Start - 서버 시작 (blocking)
func (s *Server) Start(listener net.Listener) error {
	return s.httpServer.Serve(listener)
}

// Shutdown - Graceful Shutdown 수행
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
