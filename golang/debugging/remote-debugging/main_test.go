package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestHealthHandler(t *testing.T) {
	worker := NewWorker(1, 1*time.Minute)
	srv := NewServer(worker)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	srv.healthHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}

	if resp["status"] != "ok" {
		t.Errorf("expected status ok, got %v", resp["status"])
	}
}

func TestProcessHandler(t *testing.T) {
	worker := NewWorker(1, 1*time.Minute)
	srv := NewServer(worker)

	req := httptest.NewRequest(http.MethodPost, "/process", nil)
	w := httptest.NewRecorder()

	srv.processHandler(w, req)

	if w.Code != http.StatusAccepted {
		t.Fatalf("expected status 202, got %d", w.Code)
	}

	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}

	if resp["status"] != "queued" {
		t.Errorf("expected status queued, got %v", resp["status"])
	}

	jobID, ok := resp["job_id"].(float64)
	if !ok || jobID != 1 {
		t.Errorf("expected job_id 1, got %v", resp["job_id"])
	}
}

func TestProcessHandler_MethodNotAllowed(t *testing.T) {
	worker := NewWorker(1, 1*time.Minute)
	srv := NewServer(worker)

	req := httptest.NewRequest(http.MethodGet, "/process", nil)
	w := httptest.NewRecorder()

	srv.processHandler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", w.Code)
	}
}

func TestWorkerProcessesJobs(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	worker := NewWorker(1, 1*time.Minute)

	wg.Add(1)
	go worker.Run(ctx, &wg)

	// Enqueue a job
	worker.Enqueue(Job{ID: 1, Payload: "test-job", CreatedAt: time.Now()})

	// Wait for the job to be processed
	deadline := time.After(3 * time.Second)
	for {
		select {
		case <-deadline:
			t.Fatal("timed out waiting for job to be processed")
		default:
			if worker.ProcessedCount() >= 1 {
				cancel()
				wg.Wait()
				return
			}
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func TestRegisterRoutes(t *testing.T) {
	worker := NewWorker(1, 1*time.Minute)
	srv := NewServer(worker)
	mux := http.NewServeMux()
	srv.RegisterRoutes(mux)

	// Test /health via mux
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	// Test /process via mux
	req = httptest.NewRequest(http.MethodPost, "/process", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusAccepted {
		t.Fatalf("expected status 202, got %d", w.Code)
	}
}
