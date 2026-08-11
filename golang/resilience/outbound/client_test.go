package outbound

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/failsafe-go/failsafe-go/ratelimiter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient_smooth모드는_요청을_균등하게_분산한다(t *testing.T) {
	const requests = 10

	// 한도를 넉넉히 잡아 429 없이 간격만 측정한다.
	server := newRateLimitedServer(1000)
	defer server.Close()

	client := NewClient(Options{
		Mode:          LimiterSmooth,
		MaxExecutions: 10, // 100ms 간격
		Period:        time.Second,
		MaxWaitTime:   10 * time.Second,
	})

	start := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < requests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := client.Get(fmt.Sprintf("%s/maps/map-%02d", server.URL, i))
			if err != nil {
				return
			}
			resp.Body.Close()
		}()
	}
	wg.Wait()
	elapsed := time.Since(start)

	got := server.stats()
	require.Equal(t, requests, got.Total)
	assert.Zero(t, got.Rejected)

	// 10개를 100ms 간격으로 흘리면 첫 요청 이후 9번의 간격이 생긴다.
	assert.GreaterOrEqual(t, elapsed, 800*time.Millisecond,
		"균등 분산이면 최소 0.8초는 걸려야 한다")
}

func TestNewClient_정책이_없으면_즉시_모두_보낸다(t *testing.T) {
	const requests = 10

	server := newRateLimitedServer(1000)
	defer server.Close()

	client := NewClient(Options{Mode: LimiterNone})

	start := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < requests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := client.Get(fmt.Sprintf("%s/maps/map-%02d", server.URL, i))
			if err != nil {
				return
			}
			resp.Body.Close()
		}()
	}
	wg.Wait()
	elapsed := time.Since(start)

	got := server.stats()
	require.Equal(t, requests, got.Total)
	assert.Less(t, elapsed, 500*time.Millisecond, "스로틀링이 없으면 즉시 끝난다")
	assert.Equal(t, requests, got.MaxQPS, "동시에 도달하므로 최대 QPS가 요청 수와 같다")
}

func TestNewClient_대기시간을_초과하면_ErrExceeded를_반환한다(t *testing.T) {
	const requests = 5

	server := newRateLimitedServer(1000)
	defer server.Close()

	// 주기당 1개만 허용하고 대기 상한을 짧게 잡아, 첫 요청 외에는 permit을
	// 얻기 전에 반드시 시간 초과가 나도록 만든다.
	client := NewClient(Options{
		Mode:          LimiterSmooth,
		MaxExecutions: 1,
		Period:        time.Second,
		MaxWaitTime:   20 * time.Millisecond,
	})

	errs := make([]error, requests)
	var wg sync.WaitGroup
	for i := 0; i < requests; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			resp, err := client.Get(fmt.Sprintf("%s/maps/map-%02d", server.URL, i))
			if err != nil {
				errs[i] = err
				return
			}
			resp.Body.Close()
		}(i)
	}
	wg.Wait()

	var succeeded, exceeded int
	for _, err := range errs {
		switch {
		case err == nil:
			succeeded++
		case errors.Is(err, ratelimiter.ErrExceeded):
			exceeded++
		}
	}

	assert.GreaterOrEqual(t, succeeded, 1, "적어도 하나는 즉시 permit을 얻어 성공해야 한다")
	assert.GreaterOrEqual(t, exceeded, 1, "적어도 하나는 대기 상한을 넘어 ErrExceeded를 반환해야 한다")
}
