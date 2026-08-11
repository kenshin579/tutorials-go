package outbound

import (
	"context"
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

// bursty 모드는 다른 모드와 달리 아래 비교 테스트(TestClient_버스트_요청과_스로틀링_비교)의
// 서브테스트 하나에만 등장하고, 그 서브테스트는 정책 없음과 수치가 똑같이 나오도록
// 설계돼 있다. 그래서 이 직접 테스트가 없으면 newRateLimiter가 LimiterBursty
// 분기를 잃어버려(nil을 반환해) limiter가 완전히 빠져도 비교 테스트가 계속
// 통과한다 — limiter가 있는지 없는지 구분하지 못하는 것이다. 이 테스트는 bursty의
// 두 가지 정의적 속성을 직접 검증해 그 구멍을 막는다.
func TestNewClient_bursty모드는_주기_안의_버스트를_허용하고_초과분은_ErrExceeded로_막는다(t *testing.T) {
	const (
		maxExecutions = 3               // 버스트로 허용할 개수
		period        = 2 * time.Second // 테스트 중 주기가 넘어가지 않도록 충분히 길게
		maxWaitTime   = 20 * time.Millisecond
	)

	// 벤더 한도를 넉넉히 잡아 서버의 429가 클라이언트 측 limiter 동작과
	// 섞이지 않게 한다. 여기서 보는 것은 클라이언트 limiter의 동작이다.
	server := newRateLimitedServer(1000)
	defer server.Close()

	client := NewClient(Options{
		Mode:          LimiterBursty,
		MaxExecutions: maxExecutions,
		Period:        period,
		MaxWaitTime:   maxWaitTime,
	})

	// 속성 1 (버스트는 허용된다): 주기 안의 첫 maxExecutions개는 지연 없이
	// 통과해야 한다. smooth였다면 간격을 두고 흘렸을 것이므로, 이 속성이
	// bursty와 smooth를 구분한다.
	start := time.Now()
	for i := 0; i < maxExecutions; i++ {
		resp, err := client.Get(fmt.Sprintf("%s/maps/map-%02d", server.URL, i))
		require.NoError(t, err, "허용량 안의 요청은 성공해야 한다")
		resp.Body.Close()
	}
	elapsed := time.Since(start)
	assert.Less(t, elapsed, 500*time.Millisecond, "허용량 안의 요청은 대기 없이 즉시 통과해야 한다")

	// 속성 2 (허용량은 유한하다): 허용량을 다 쓴 다음 요청은 짧은 대기 상한
	// 안에 permit을 못 얻어 ErrExceeded로 막혀야 한다. LimiterNone이었다면
	// 그대로 통과했을 것이므로, 이 속성이 bursty와 정책 없음을 구분한다.
	_, err := client.Get(fmt.Sprintf("%s/maps/map-%02d", server.URL, maxExecutions))
	require.Error(t, err, "허용량을 넘은 요청은 실패해야 한다")
	assert.True(t, errors.Is(err, ratelimiter.ErrExceeded),
		"허용량 초과는 ErrExceeded로 나타나야 한다")
}

func TestClient_버스트_요청과_스로틀링_비교(t *testing.T) {
	const (
		serverLimit = 20 // 벤더 한도
		requests    = 40 // 로봇 4대 x 지도 10개
	)

	tests := []struct {
		name         string
		opts         Options
		expectReject bool
	}{
		{
			name:         "정책 없음 - 동시에 발사하면 429가 발생한다",
			opts:         Options{Mode: LimiterNone},
			expectReject: true,
		},
		{
			// 40회 / 2초 = 평균 20 QPS로 한도에 정확히 맞춘 설정이다.
			// 그런데도 주기 시작에 40개를 한꺼번에 통과시켜 순간 QPS가 튄다.
			// "1분에 1200회까지 가능하다"는 계산이 안전을 보장하지 않는 이유다.
			//
			// 이 서브테스트의 수치는 "정책 없음" 서브테스트와 동일하게 나오도록
			// 의도됐다 — 주기 시작에 허용량을 한꺼번에 다 써버리는 bursty
			// limiter는 이 40개짜리 단발 버스트 앞에서는 limiter가 없는 것과
			// 구분되지 않는다는 것 자체가 가르치려는 요점이다. 버그가 아니다.
			// 다만 그래서 이 서브테스트만으로는 bursty limiter가 실제로
			// 붙어 있는지(vs 우연히 삭제됐는지) 증명할 수 없다 — 그 증명은
			// TestNewClient_bursty모드는_주기_안의_버스트를_허용하고_초과분은_ErrExceeded로_막는다
			// 가 담당한다.
			name: "bursty - 평균 QPS를 한도에 맞춰도 429가 발생한다",
			opts: Options{
				Mode:          LimiterBursty,
				MaxExecutions: 40,
				Period:        2 * time.Second,
				MaxWaitTime:   30 * time.Second,
			},
			expectReject: true,
		},
		{
			// 한도의 50%인 10 QPS로 균등 분산한다.
			// 20 QPS로 두면 sliding window 경계에서 21개로 셀 여지가 있어
			// 마진이 없다. 벤더 한도의 50~70%만 쓰는 것이 권장된다.
			name: "smooth - 균등 분산하면 429가 사라진다",
			opts: Options{
				Mode:          LimiterSmooth,
				MaxExecutions: 10,
				Period:        time.Second,
				MaxWaitTime:   30 * time.Second,
			},
			expectReject: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 이 테스트는 "429가 났는가"를 비교한다. 재시도를 켜면 429를 받은
			// 요청이 다시 성공해 서버가 센 429 수가 늘고 케이스 간 차이가
			// 흐려지므로 안전망을 끈다. 재시도 효과는 별도 테스트에서 본다.
			require.Zero(t, tt.opts.MaxRetries, "비교 테스트는 재시도를 꺼야 한다")

			server := newRateLimitedServer(serverLimit)
			defer server.Close()

			refresher := NewRefresher(NewClient(tt.opts), server.URL)

			keys := make([]string, requests)
			for i := range keys {
				keys[i] = fmt.Sprintf("map-%02d", i)
			}

			stats := refresher.RefreshAll(context.Background(), keys)
			got := server.stats()

			t.Logf("총 호출 %d, 429 %d, 최대 관측 QPS %d, 소요 %v",
				got.Total, got.Rejected, got.MaxQPS, stats.Elapsed)

			assert.Equal(t, requests, got.Total, "재시도가 없으므로 키당 정확히 1회 호출된다")

			if tt.expectReject {
				assert.Positive(t, got.Rejected, "429가 발생해야 한다")
				assert.Greater(t, got.MaxQPS, serverLimit, "순간 QPS가 한도를 넘어야 한다")
				assert.Positive(t, stats.Failed, "429를 받은 키는 갱신에 실패한다")
			} else {
				assert.Zero(t, got.Rejected, "429가 없어야 한다")
				assert.LessOrEqual(t, got.MaxQPS, serverLimit, "순간 QPS가 한도 이하여야 한다")
				assert.Zero(t, stats.Failed, "모든 키가 갱신되어야 한다")
			}
		})
	}
}

func TestClient_429를_받으면_RetryAfter만큼_기다린다(t *testing.T) {
	// 이 테스트는 실측 약 1초가 걸린다("수 초"는 아니다). 그런데도 이 테스트만
	// -short로 건너뛰게 하는 이유는 시간 자체가 아니라 검증 대상이다 — 이
	// 테스트는 실제 Retry-After 대기를 시간으로 재서 확인하는 유일한 테스트라
	// 그 대기를 없앨 수 없다. 반면 위 TestClient_버스트_요청과_스로틀링_비교의
	// smooth 서브테스트는 3.9초로 더 오래 걸리지만 이 예제의 헤드라인
	// 데모이므로 항상 실행되게 남겨둔다.
	if testing.Short() {
		t.Skip("Retry-After 대기로 약 1초가 걸린다")
	}

	const (
		serverLimit = 5
		requests    = 10
	)

	server := newRateLimitedServer(serverLimit)
	defer server.Close()

	// rate limiter 없이 발사해 일부러 429를 유발하고, 재시도로 복구되는지 본다.
	// failsafehttp.NewRetryPolicyBuilder가 Retry-After 헤더(1초)를 존중한다.
	client := NewClient(Options{
		Mode:       LimiterNone,
		MaxRetries: 5,
	})
	refresher := NewRefresher(client, server.URL)

	keys := make([]string, requests)
	for i := range keys {
		keys[i] = fmt.Sprintf("map-%02d", i)
	}

	stats := refresher.RefreshAll(context.Background(), keys)
	got := server.stats()

	t.Logf("총 호출 %d, 429 %d, 소요 %v", got.Total, got.Rejected, stats.Elapsed)

	assert.Positive(t, got.Rejected, "한도를 넘겼으므로 429가 발생해야 한다")
	assert.Greater(t, got.Total, requests, "재시도만큼 서버 호출이 늘어난다")
	assert.Zero(t, stats.Failed, "재시도로 모든 키가 최종 성공해야 한다")
	assert.GreaterOrEqual(t, stats.Elapsed, time.Second, "Retry-After만큼 기다려야 한다")
}
