package circuitbreaker

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sony/gobreaker/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestHTTPClient() *ProtectedHTTPClient {
	settings := gobreaker.Settings{
		Name:        "test-http",
		MaxRequests: 1,
		Interval:    0,
		Timeout:     1 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 3
		},
	}
	return &ProtectedHTTPClient{
		client:  &http.Client{Timeout: 5 * time.Second},
		breaker: gobreaker.NewCircuitBreaker[*http.Response](settings),
	}
}

func TestHTTP_SuccessfulRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}))
	defer server.Close()

	client := newTestHTTPClient()
	req, _ := http.NewRequest(http.MethodGet, server.URL, nil)

	resp, err := client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()
	assert.Equal(t, gobreaker.StateClosed, client.State())
}

func TestHTTP_ServerError_TripsBreaker(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := newTestHTTPClient()

	// 3회 5xx 응답 → Open 전환
	for i := 0; i < 3; i++ {
		req, _ := http.NewRequest(http.MethodGet, server.URL, nil)
		resp, err := client.Do(req)
		assert.Error(t, err)
		if resp != nil {
			resp.Body.Close()
		}
	}

	assert.Equal(t, gobreaker.StateOpen, client.State())

	// Open 상태에서 즉시 거부
	req, _ := http.NewRequest(http.MethodGet, server.URL, nil)
	_, err := client.Do(req)
	assert.ErrorIs(t, err, gobreaker.ErrOpenState)
}

func TestHTTP_RecoveryAfterServerFix(t *testing.T) {
	var shouldFail atomic.Bool
	shouldFail.Store(true)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if shouldFail.Load() {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("recovered"))
	}))
	defer server.Close()

	client := newTestHTTPClient()

	// 서버 장애로 Open 전환
	for i := 0; i < 3; i++ {
		req, _ := http.NewRequest(http.MethodGet, server.URL, nil)
		resp, _ := client.Do(req)
		if resp != nil {
			resp.Body.Close()
		}
	}
	assert.Equal(t, gobreaker.StateOpen, client.State())

	// 서버 복구
	shouldFail.Store(false)

	// Half-Open 대기
	time.Sleep(1200 * time.Millisecond)
	assert.Equal(t, gobreaker.StateHalfOpen, client.State())

	// Half-Open에서 성공 → Closed 복귀
	req, _ := http.NewRequest(http.MethodGet, server.URL, nil)
	resp, err := client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()
	assert.Equal(t, gobreaker.StateClosed, client.State())
}
