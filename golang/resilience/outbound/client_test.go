package outbound

import (
	"fmt"
	"sync"
	"testing"
	"time"

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
