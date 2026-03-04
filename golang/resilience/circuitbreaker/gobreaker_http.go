package circuitbreaker

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/sony/gobreaker/v2"
)

// ProtectedHTTPClient는 Circuit Breaker로 보호되는 HTTP 클라이언트이다.
type ProtectedHTTPClient struct {
	client  *http.Client
	breaker *gobreaker.CircuitBreaker[*http.Response]
}

// NewProtectedHTTPClient는 기본 설정으로 보호된 HTTP 클라이언트를 생성한다.
func NewProtectedHTTPClient() *ProtectedHTTPClient {
	settings := gobreaker.Settings{
		Name:        "http-client",
		MaxRequests: 2,
		Interval:    10 * time.Second,
		Timeout:     30 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 3
		},
		OnStateChange: func(name string, from, to gobreaker.State) {
			log.Printf("HTTP CB %s: %s → %s", name, from, to)
		},
	}

	return &ProtectedHTTPClient{
		client:  &http.Client{Timeout: 10 * time.Second},
		breaker: gobreaker.NewCircuitBreaker[*http.Response](settings),
	}
}

// Do는 Circuit Breaker로 보호된 HTTP 요청을 실행한다.
// 5xx 응답은 실패로 처리하여 Circuit Breaker에 반영된다.
func (c *ProtectedHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return c.breaker.Execute(func() (*http.Response, error) {
		resp, err := c.client.Do(req)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode >= 500 {
			return resp, fmt.Errorf("server error: %d", resp.StatusCode)
		}
		return resp, nil
	})
}

// State는 현재 Circuit Breaker 상태를 반환한다.
func (c *ProtectedHTTPClient) State() gobreaker.State {
	return c.breaker.State()
}
