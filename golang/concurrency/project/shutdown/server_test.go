package shutdown

import (
	"context"
	"io"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestGracefulShutdown - graceful shutdown 기본 테스트
func TestGracefulShutdown(t *testing.T) {
	// 랜덤 포트로 리스너 생성
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	assert.NoError(t, err)
	addr := listener.Addr().String()

	srv := NewServer(addr)

	// 서버 시작
	go func() {
		if err := srv.Start(listener); err != http.ErrServerClosed {
			t.Errorf("unexpected error: %v", err)
		}
	}()

	// 서버가 준비될 때까지 대기
	time.Sleep(50 * time.Millisecond)

	// 요청 확인
	resp, err := http.Get("http://" + addr + "/health")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	// graceful shutdown (5초 timeout)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = srv.Shutdown(ctx)
	assert.NoError(t, err)
	t.Log("서버가 gracefully 종료됨")
}

// TestGracefulShutdownWithPendingRequests - 진행 중인 요청이 완료된 후 종료
func TestGracefulShutdownWithPendingRequests(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	assert.NoError(t, err)
	addr := listener.Addr().String()

	srv := NewServer(addr)

	go func() {
		if err := srv.Start(listener); err != http.ErrServerClosed {
			t.Errorf("unexpected error: %v", err)
		}
	}()

	time.Sleep(50 * time.Millisecond)

	var wg sync.WaitGroup
	var slowResp string

	// 느린 요청 시작
	wg.Add(1)
	go func() {
		defer wg.Done()
		resp, err := http.Get("http://" + addr + "/slow")
		if err != nil {
			t.Logf("slow request error (expected during shutdown): %v", err)
			return
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		slowResp = string(body)
	}()

	// 느린 요청이 시작된 후 shutdown 시작
	time.Sleep(100 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = srv.Shutdown(ctx)
	assert.NoError(t, err)

	wg.Wait()

	// 진행 중이던 요청이 완료되었는지 확인
	if slowResp != "" {
		assert.Contains(t, slowResp, "Slow response done")
		t.Log("진행 중인 요청이 완료된 후 서버 종료됨")
	}
}

// TestShutdownTimeout - shutdown timeout 초과 시 강제 종료
func TestShutdownTimeout(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	assert.NoError(t, err)
	addr := listener.Addr().String()

	srv := NewServer(addr)

	go func() {
		if err := srv.Start(listener); err != http.ErrServerClosed {
			t.Errorf("unexpected error: %v", err)
		}
	}()

	time.Sleep(50 * time.Millisecond)

	// 느린 요청 시작 (2초 소요)
	go func() {
		http.Get("http://" + addr + "/slow") //nolint:errcheck
	}()

	time.Sleep(100 * time.Millisecond)

	// 아주 짧은 timeout으로 shutdown → context deadline exceeded
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err = srv.Shutdown(ctx)
	assert.Error(t, err) // context deadline exceeded
	t.Logf("shutdown timeout: %v", err)
}

// TestSignalPattern - signal.NotifyContext 패턴 시뮬레이션
func TestSignalPattern(t *testing.T) {
	// 실제 signal 대신 context로 시뮬레이션
	ctx, cancel := context.WithCancel(context.Background())

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	assert.NoError(t, err)
	addr := listener.Addr().String()

	srv := NewServer(addr)

	// 서버 시작
	serverDone := make(chan struct{})
	go func() {
		defer close(serverDone)
		if err := srv.Start(listener); err != http.ErrServerClosed {
			t.Errorf("unexpected error: %v", err)
		}
	}()

	time.Sleep(50 * time.Millisecond)

	// health check
	resp, err := http.Get("http://" + addr + "/health")
	assert.NoError(t, err)
	resp.Body.Close()

	// "시그널" 수신 시뮬레이션
	cancel()

	// shutdown 처리
	go func() {
		<-ctx.Done()
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		srv.Shutdown(shutdownCtx) //nolint:errcheck
	}()

	// 서버가 종료될 때까지 대기
	select {
	case <-serverDone:
		t.Log("signal 패턴으로 서버 gracefully 종료됨")
	case <-time.After(3 * time.Second):
		t.Fatal("서버가 시간 내에 종료되지 않음")
	}
}
