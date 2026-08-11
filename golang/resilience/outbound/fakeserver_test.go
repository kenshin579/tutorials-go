package outbound

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// serverStats는 가짜 서버가 관측한 호출 통계다.
type serverStats struct {
	Total    int // 총 요청 수
	Rejected int // 429로 거절한 요청 수
	MaxQPS   int // 관측된 최대 초당 요청 수
}

// rateLimitedServer는 벤더 API의 sliding window rate limit을 흉내내는 테스트 서버다.
//
// 고정 윈도우 대신 sliding window를 쓰는 이유는, 고정 윈도우가 경계에서 최대
// 2배를 허용해 클라이언트 주기와 우연히 정렬되면 버스트 설정으로도 429가 나지
// 않기 때문이다. 그러면 테스트 통과가 처방 덕인지 정렬 운인지 구분할 수 없다.
type rateLimitedServer struct {
	*httptest.Server

	mu       sync.Mutex
	limit    int         // 1초당 허용 요청 수
	arrivals []time.Time // 최근 1초 내 요청 도착 시각
	total    int
	rejected int
	maxQPS   int
}

// newRateLimitedServer는 초당 limit개까지 허용하는 가짜 서버를 시작한다.
func newRateLimitedServer(limit int) *rateLimitedServer {
	s := &rateLimitedServer{limit: limit}
	s.Server = httptest.NewServer(http.HandlerFunc(s.handle))
	return s
}

func (s *rateLimitedServer) handle(w http.ResponseWriter, _ *http.Request) {
	s.mu.Lock()
	now := time.Now()
	s.total++

	// 1초보다 오래된 도착 기록을 버린다.
	cutoff := now.Add(-time.Second)
	kept := s.arrivals[:0]
	for _, at := range s.arrivals {
		if at.After(cutoff) {
			kept = append(kept, at)
		}
	}
	s.arrivals = append(kept, now)

	if len(s.arrivals) > s.maxQPS {
		s.maxQPS = len(s.arrivals)
	}
	exceeded := len(s.arrivals) > s.limit
	if exceeded {
		s.rejected++
	}
	s.mu.Unlock()

	if exceeded {
		w.Header().Set("Retry-After", "1")
		w.WriteHeader(http.StatusTooManyRequests)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

// stats는 지금까지 관측한 통계를 반환한다.
func (s *rateLimitedServer) stats() serverStats {
	s.mu.Lock()
	defer s.mu.Unlock()
	return serverStats{Total: s.total, Rejected: s.rejected, MaxQPS: s.maxQPS}
}

func TestRateLimitedServer_한도를_넘으면_429를_반환한다(t *testing.T) {
	const (
		limit    = 20
		requests = 25
	)

	server := newRateLimitedServer(limit)
	defer server.Close()

	// 25개를 동시에 발사하면 21번째부터 한도를 넘는다.
	var wg sync.WaitGroup
	for i := 0; i < requests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := http.Get(fmt.Sprintf("%s/maps/map-%02d", server.URL, i))
			if err != nil {
				return
			}
			resp.Body.Close()
		}()
	}
	wg.Wait()

	got := server.stats()
	assert.Equal(t, requests, got.Total, "모든 요청이 서버에 도달해야 한다")
	assert.Equal(t, requests-limit, got.Rejected, "한도 초과분만 429를 받아야 한다")
	assert.Equal(t, requests, got.MaxQPS, "동시 발사이므로 최대 관측 QPS가 요청 수와 같다")
}

func TestRateLimitedServer_1초가_지나면_다시_허용한다(t *testing.T) {
	const limit = 5

	server := newRateLimitedServer(limit)
	defer server.Close()

	send := func() int {
		resp, err := http.Get(server.URL + "/maps/map-00")
		require.NoError(t, err)
		defer resp.Body.Close()
		return resp.StatusCode
	}

	// 한도를 정확히 채운다.
	for i := 0; i < limit; i++ {
		require.Equal(t, http.StatusOK, send())
	}
	// 한 개 더 보내면 거절된다.
	require.Equal(t, http.StatusTooManyRequests, send())

	// sliding window가 비워질 때까지 기다리면 다시 통과한다.
	time.Sleep(1100 * time.Millisecond)
	assert.Equal(t, http.StatusOK, send())
}
