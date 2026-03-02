package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Job represents a unit of work for the background worker.
type Job struct {
	ID        int
	Payload   string
	CreatedAt time.Time
}

// Worker processes jobs in the background.
type Worker struct {
	id       int
	interval time.Duration
	jobs     chan Job
	mu       sync.Mutex
	count    int
}

// NewWorker creates a new background worker.
func NewWorker(id int, interval time.Duration) *Worker {
	return &Worker{
		id:       id,
		interval: interval,
		jobs:     make(chan Job, 10),
	}
}

// Run starts the worker loop. It processes jobs from the channel
// and performs periodic background tasks.
func (w *Worker) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	log.Printf("[Worker %d] started", w.id)

	for {
		select {
		case <-ctx.Done():
			log.Printf("[Worker %d] shutting down (processed %d jobs)", w.id, w.ProcessedCount())
			return
		case job := <-w.jobs:
			w.processJob(job)
		case <-ticker.C:
			w.doPeriodicTask()
		}
	}
}

func (w *Worker) processJob(job Job) {
	log.Printf("[Worker %d] processing job #%d: %s", w.id, job.ID, job.Payload)
	// Simulate work
	time.Sleep(time.Duration(100+rand.Intn(200)) * time.Millisecond)
	w.mu.Lock()
	w.count++
	w.mu.Unlock()
	log.Printf("[Worker %d] completed job #%d", w.id, job.ID)
}

func (w *Worker) doPeriodicTask() {
	log.Printf("[Worker %d] performing periodic health check", w.id)
}

// ProcessedCount returns the number of jobs processed.
func (w *Worker) ProcessedCount() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.count
}

// Enqueue adds a job to the worker's queue.
func (w *Worker) Enqueue(job Job) {
	w.jobs <- job
}

// Server holds HTTP server dependencies.
type Server struct {
	worker  *Worker
	jobSeq  int
	mu      sync.Mutex
}

// NewServer creates a new Server instance.
func NewServer(worker *Worker) *Server {
	return &Server{worker: worker}
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"ok","worker_jobs_processed":%d}`, s.worker.ProcessedCount())
}

func (s *Server) processHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	s.mu.Lock()
	s.jobSeq++
	id := s.jobSeq
	s.mu.Unlock()

	job := Job{
		ID:        id,
		Payload:   fmt.Sprintf("request-%d", id),
		CreatedAt: time.Now(),
	}

	// Simulate some request processing delay
	time.Sleep(50 * time.Millisecond)

	s.worker.Enqueue(job)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(w, `{"job_id":%d,"status":"queued"}`, id)
}

// RegisterRoutes sets up the HTTP routes.
func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", s.healthHandler)
	mux.HandleFunc("/process", s.processHandler)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	// Create and start background worker
	worker := NewWorker(1, 5*time.Second)
	wg.Add(1)
	go worker.Run(ctx, &wg)

	// Set up HTTP server
	srv := NewServer(worker)
	mux := http.NewServeMux()
	srv.RegisterRoutes(mux)

	addr := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}

	httpServer := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	// Start HTTP server in a goroutine
	go func() {
		log.Printf("HTTP server listening on %s", addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Println("Received shutdown signal")

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	cancel() // Signal workers to stop
	wg.Wait()
	log.Println("Server stopped gracefully")
}
