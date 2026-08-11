# 외부 API 버스트 호출 방지 예제 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 캐시 갱신 시 외부 벤더 API를 버스트로 호출해 429를 받는 문제를 재현하고, smooth rate limiter로 해결되는 것을 테스트로 증명하는 예제를 만든다.

**Architecture:** `failsafehttp.NewRoundTripper`로 rate limiter·bulkhead·재시도 정책을 `http.RoundTripper` 층에 조립한다. HTTP 클라이언트 종류(net/http, resty 등)와 무관하게 적용된다. sliding window rate limit을 흉내내는 `httptest` 서버가 429 수와 관측 최대 QPS를 세어, 정책별 효과를 숫자로 비교한다.

**Tech Stack:** Go 1.26, `github.com/failsafe-go/failsafe-go v0.9.6` (ratelimiter, bulkhead, failsafehttp), `golang.org/x/sync/singleflight`, `testify/assert`

**설계 문서:** `docs/superpowers/specs/2026-08-10-outbound-burst-ratelimit-design.md`

**브랜치:** `feature/outbound-burst-ratelimit` (이미 생성됨)

---

## 스펙과 다른 점 (의도된 것)

- 스펙의 `refresher.go`는 `errgroup`을 쓴다고 했으나 **`sync.WaitGroup` + `atomic.Int64`** 를 쓴다. 개별 키의 갱신 실패가 형제 goroutine을 취소하면 안 되므로(stale 서빙) 모든 콜백이 `nil`을 반환하게 되는데, 그러면 `errgroup`이 `WaitGroup` 대비 아무것도 더 해주지 않는다.
- 스펙의 `Stats`에 `Fetched` 필드가 있었으나 **`Unique`** 로 이름을 바꾼다. 실제 서버 호출 수는 가짜 서버가 세므로, `Refresher`가 보고할 값은 "중복 제거 후 키 개수"다.

---

## 파일 구조

| 파일 | 책임 |
|---|---|
| `golang/resilience/outbound/client.go` | `Options`/`LimiterMode` 정의, 정책을 조립한 `*http.Client` 생성 |
| `golang/resilience/outbound/refresher.go` | 키 중복 제거, singleflight, 병렬 갱신, stale 캐시 유지 |
| `golang/resilience/outbound/fakeserver_test.go` | sliding window rate limit을 흉내내는 `httptest` 서버 (측정 도구) |
| `golang/resilience/outbound/client_test.go` | 정책 없음 vs bursty vs smooth 비교, 429 재시도 검증 |
| `golang/resilience/outbound/refresher_test.go` | 키 중복 제거로 호출 수가 줄어드는지 검증 |
| `golang/resilience/outbound/README.md` | 문제 정의, Best Practice 계층, 테스트 결과 |
| `golang/resilience/README.md` | 구조 트리에 `outbound/` 추가 (수정) |

---

## Task 1: 가짜 벤더 서버

측정 도구부터 만든다. 이 서버가 세는 429 수와 최대 QPS가 이후 모든 테스트의 근거이므로, 서버 자체가 올바른지 먼저 검증한다.

**Files:**
- Create: `golang/resilience/outbound/fakeserver_test.go`

- [ ] **Step 1: 실패하는 테스트 작성**

`golang/resilience/outbound/fakeserver_test.go` 파일을 만들고 아래 내용을 넣는다. (구현체는 Step 3에서 같은 파일에 추가한다.)

```go
package outbound

import (
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
```

- [ ] **Step 2: 테스트가 실패하는지 확인**

Run:
```bash
go test ./golang/resilience/outbound/... -run TestRateLimitedServer -v
```

Expected: 컴파일 실패. `undefined: newRateLimitedServer`

- [ ] **Step 3: 가짜 서버 구현**

`golang/resilience/outbound/fakeserver_test.go`의 `import` 블록 바로 아래(테스트 함수들 위)에 아래 내용을 추가한다. `import` 블록에 `"net/http/httptest"`를 추가한다.

```go
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
```

- [ ] **Step 4: 테스트 통과 확인**

Run:
```bash
go test ./golang/resilience/outbound/... -run TestRateLimitedServer -v
```

Expected: PASS 2건. 두 번째 테스트는 1.1초가 걸린다.

- [ ] **Step 5: 커밋**

```bash
git add golang/resilience/outbound/fakeserver_test.go
git commit -m "test: 벤더 rate limit을 흉내내는 sliding window 가짜 서버 추가"
```

---

## Task 2: 정책을 조립한 HTTP 클라이언트

**Files:**
- Create: `golang/resilience/outbound/client.go`
- Create: `golang/resilience/outbound/client_test.go`

- [ ] **Step 1: 실패하는 테스트 작성**

`golang/resilience/outbound/client_test.go`를 만들고 아래 내용을 넣는다.

```go
package outbound

import (
	"fmt"
	"net/http"
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
```

- [ ] **Step 2: 테스트가 실패하는지 확인**

Run:
```bash
go test ./golang/resilience/outbound/... -run TestNewClient -v
```

Expected: 컴파일 실패. `undefined: NewClient`, `undefined: Options`, `undefined: LimiterSmooth`, `undefined: LimiterNone`

- [ ] **Step 3: 클라이언트 구현**

`golang/resilience/outbound/client.go`를 만들고 아래 내용을 넣는다.

```go
// Package outbound는 외부 벤더 API를 rate limit 안에서 안전하게 호출하는
// 예제를 담는다.
//
// 핵심은 평균 QPS가 아니라 순간 QPS를 제한하는 것이다. 1분에 40회는 평균
// 0.67 QPS지만, 40개를 동시에 발사하면 순간 40 QPS가 되어 20 QPS 한도를
// 넘는다.
package outbound

import (
	"net/http"
	"time"

	"github.com/failsafe-go/failsafe-go"
	"github.com/failsafe-go/failsafe-go/bulkhead"
	"github.com/failsafe-go/failsafe-go/failsafehttp"
	"github.com/failsafe-go/failsafe-go/ratelimiter"
)

// LimiterMode는 rate limiter의 버스트 허용 방식을 정한다.
type LimiterMode int

const (
	// LimiterNone은 rate limiting을 적용하지 않는다.
	LimiterNone LimiterMode = iota
	// LimiterBursty는 주기 시작 시점에 MaxExecutions개를 지연 없이 통과시킨다.
	// 평균 QPS는 맞지만 순간 QPS가 튄다.
	LimiterBursty
	// LimiterSmooth는 Period/MaxExecutions 간격으로 실행을 균등 분산한다.
	// 토큰을 모아두지 않으므로 버스트가 구조적으로 불가능하다.
	LimiterSmooth
)

// Options는 NewClient의 동작을 설정한다.
type Options struct {
	Mode           LimiterMode   // rate limiter 방식
	MaxExecutions  int           // 주기당 허용 실행 수
	Period         time.Duration // 주기
	MaxWaitTime    time.Duration // limiter/bulkhead 대기 상한. 초과 시 ratelimiter.ErrExceeded
	MaxConcurrency int           // 동시 실행 수 상한. 0이면 bulkhead 없음
	MaxRetries     int           // 재시도 횟수. 0이면 재시도 없음
	Base           http.RoundTripper
}

// NewClient는 재시도, rate limiter, bulkhead 정책을 조합한 http.Client를 만든다.
//
// 정책을 http.RoundTripper 층에 두므로 net/http, resty 등 클라이언트 종류와
// 무관하게 같은 방식으로 적용된다. rate limit은 보통 계정 단위이므로 클라이언트
// 하나(=정책 하나)를 모든 호출 지점이 공유해야 한다. 호출 지점마다 limiter를
// 따로 두면 합산해서 다시 한도를 넘는다.
func NewClient(opts Options) *http.Client {
	// 정책 순서 주의: failsafe는 첫 인자가 가장 바깥 정책이다.
	// retry를 가장 바깥에 두어야 재시도 요청도 rate limiter를 다시 통과한다.
	// 순서를 뒤집으면 재시도가 limiter를 우회해 429 폭풍을 만든다.
	var policies []failsafe.Policy[*http.Response]

	if opts.MaxRetries > 0 {
		// NewRetryPolicyBuilder는 429와 대부분의 5xx를 재시도하고,
		// Retry-After 헤더가 있으면 그 값만큼 기다린다.
		policies = append(policies, failsafehttp.NewRetryPolicyBuilder().
			WithMaxRetries(opts.MaxRetries).
			Build())
	}

	if limiter := newRateLimiter(opts); limiter != nil {
		policies = append(policies, limiter)
	}

	if opts.MaxConcurrency > 0 {
		// rate limiter는 실행 시작 시점만 제어하고 in-flight 수는 제한하지
		// 않는다. 응답이 느려지면 동시 연결이 쌓이므로 별도 층이 필요하다.
		policies = append(policies, bulkhead.NewBuilder[*http.Response](uint(opts.MaxConcurrency)).
			WithMaxWaitTime(opts.MaxWaitTime).
			Build())
	}

	return &http.Client{
		Transport: failsafehttp.NewRoundTripper(opts.Base, policies...),
	}
}

// newRateLimiter는 Mode에 맞는 rate limiter를 만든다.
// LimiterNone이면 nil을 반환한다.
func newRateLimiter(opts Options) ratelimiter.RateLimiter[*http.Response] {
	var builder ratelimiter.Builder[*http.Response]

	switch opts.Mode {
	case LimiterSmooth:
		builder = ratelimiter.NewSmoothBuilder[*http.Response](uint(opts.MaxExecutions), opts.Period)
	case LimiterBursty:
		builder = ratelimiter.NewBurstyBuilder[*http.Response](uint(opts.MaxExecutions), opts.Period)
	default:
		return nil
	}

	return builder.WithMaxWaitTime(opts.MaxWaitTime).Build()
}
```

- [ ] **Step 4: 의존성 정리 후 테스트 통과 확인**

Run:
```bash
go mod tidy && go test ./golang/resilience/outbound/... -run TestNewClient -v
```

Expected: PASS 2건. `go.mod`에서 `github.com/failsafe-go/failsafe-go`가 `// indirect` 주석을 잃고 direct 블록으로 이동한다.

- [ ] **Step 5: 커밋**

```bash
git add golang/resilience/outbound/client.go golang/resilience/outbound/client_test.go go.mod go.sum
git commit -m "feat: rate limiter/bulkhead/재시도를 조립한 HTTP 클라이언트 추가"
```

---

## Task 3: 키 중복 제거와 캐시 갱신

**Files:**
- Create: `golang/resilience/outbound/refresher.go`
- Create: `golang/resilience/outbound/refresher_test.go`

- [ ] **Step 1: 실패하는 테스트 작성**

`golang/resilience/outbound/refresher_test.go`를 만들고 아래 내용을 넣는다.

```go
package outbound

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRefresher_키_중복제거로_호출수가_줄어든다(t *testing.T) {
	const (
		robots      = 4
		maps        = 10
		serverLimit = 20
	)

	server := newRateLimitedServer(serverLimit)
	defer server.Close()

	client := NewClient(Options{
		Mode:          LimiterSmooth,
		MaxExecutions: 10,
		Period:        time.Second,
		MaxWaitTime:   10 * time.Second,
	})
	refresher := NewRefresher(client, server.URL)

	// 로봇마다 같은 지도 목록을 요청한다.
	// 캐시 키에 로봇 ID가 섞이면 4 x 10 = 40회를 호출하지만,
	// 지도 단위로 정규화하면 10회면 된다.
	var keys []string
	for robot := 0; robot < robots; robot++ {
		for m := 0; m < maps; m++ {
			keys = append(keys, fmt.Sprintf("map-%02d", m))
		}
	}

	stats := refresher.RefreshAll(context.Background(), keys)
	got := server.stats()

	t.Logf("입력 키 %d, 중복 제거 후 %d, 서버 호출 %d, 429 %d, 소요 %v",
		stats.Requested, stats.Unique, got.Total, got.Rejected, stats.Elapsed)

	assert.Equal(t, robots*maps, stats.Requested)
	assert.Equal(t, maps, stats.Unique, "지도 단위로 중복이 제거되어야 한다")
	assert.Equal(t, maps, got.Total, "서버 호출은 지도 개수만큼만 일어나야 한다")
	assert.Zero(t, got.Rejected)
	assert.Zero(t, stats.Failed)
}

func TestRefresher_갱신에_실패해도_이전_캐시값을_유지한다(t *testing.T) {
	// 한도 0이면 모든 요청이 429를 받는다.
	server := newRateLimitedServer(0)
	defer server.Close()

	client := NewClient(Options{Mode: LimiterNone})
	refresher := NewRefresher(client, server.URL)

	refresher.set("map-00", "이전값")

	stats := refresher.RefreshAll(context.Background(), []string{"map-00"})
	assert.Equal(t, 1, stats.Failed, "갱신은 실패해야 한다")

	value, ok := refresher.Get("map-00")
	assert.True(t, ok, "실패해도 캐시 항목이 남아야 한다")
	assert.Equal(t, "이전값", value, "이전 값이 덮어써지면 안 된다")
}

func TestDedupe_순서를_유지하며_중복을_제거한다(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want []string
	}{
		{
			name: "중복 없음",
			in:   []string{"a", "b", "c"},
			want: []string{"a", "b", "c"},
		},
		{
			name: "연속 중복",
			in:   []string{"a", "a", "b"},
			want: []string{"a", "b"},
		},
		{
			name: "떨어진 중복은 첫 등장 순서를 유지한다",
			in:   []string{"b", "a", "b", "c", "a"},
			want: []string{"b", "a", "c"},
		},
		{
			name: "빈 입력",
			in:   nil,
			want: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, dedupe(tt.in))
		})
	}
}
```

- [ ] **Step 2: 테스트가 실패하는지 확인**

Run:
```bash
go test ./golang/resilience/outbound/... -run 'TestRefresher|TestDedupe' -v
```

Expected: 컴파일 실패. `undefined: NewRefresher`, `undefined: dedupe`

- [ ] **Step 3: Refresher 구현**

`golang/resilience/outbound/refresher.go`를 만들고 아래 내용을 넣는다.

```go
package outbound

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/sync/singleflight"
)

// Stats는 RefreshAll 한 번의 결과를 요약한다.
type Stats struct {
	Requested int           // 입력된 키 개수 (중복 포함)
	Unique    int           // 중복 제거 후 키 개수
	Failed    int           // 갱신에 실패한 키 개수
	Elapsed   time.Duration // 전체 소요 시간
}

// Refresher는 키 목록에 대한 캐시 갱신을 수행한다.
type Refresher struct {
	client  *http.Client
	baseURL string
	group   singleflight.Group

	mu    sync.RWMutex
	cache map[string]string
}

// NewRefresher는 client로 baseURL을 호출하는 Refresher를 만든다.
func NewRefresher(client *http.Client, baseURL string) *Refresher {
	return &Refresher{
		client:  client,
		baseURL: baseURL,
		cache:   make(map[string]string),
	}
}

// RefreshAll은 keys의 중복을 제거한 뒤 병렬로 캐시를 갱신한다.
//
// 갱신에 실패한 키는 기존 캐시 값을 그대로 유지한다(stale 서빙). 캐시 갱신
// 실패가 서비스 응답 실패로 번지지 않게 하기 위함이다. 같은 이유로 개별 키의
// 실패가 다른 키의 갱신을 취소하지 않는다.
func (r *Refresher) RefreshAll(ctx context.Context, keys []string) Stats {
	start := time.Now()
	unique := dedupe(keys)

	var failed atomic.Int64
	var wg sync.WaitGroup
	for _, key := range unique {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := r.refresh(ctx, key); err != nil {
				failed.Add(1)
			}
		}()
	}
	wg.Wait()

	return Stats{
		Requested: len(keys),
		Unique:    len(unique),
		Failed:    int(failed.Load()),
		Elapsed:   time.Since(start),
	}
}

// refresh는 키 하나를 갱신한다.
// 같은 키에 대해 이미 진행 중인 호출이 있으면 그 결과를 공유한다.
func (r *Refresher) refresh(ctx context.Context, key string) error {
	_, err, _ := r.group.Do(key, func() (any, error) {
		value, err := r.fetch(ctx, key)
		if err != nil {
			return nil, err
		}
		r.set(key, value)
		return value, nil
	})
	return err
}

// fetch는 벤더 API를 호출한다.
func (r *Refresher) fetch(ctx context.Context, key string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, r.baseURL+"/maps/"+key, nil)
	if err != nil {
		return "", err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status %d for key %s", resp.StatusCode, key)
	}

	return string(body), nil
}

// Get은 캐시된 값을 반환한다.
func (r *Refresher) Get(key string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	value, ok := r.cache[key]
	return value, ok
}

// set은 캐시 값을 저장한다.
func (r *Refresher) set(key, value string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cache[key] = value
}

// dedupe는 순서를 유지하며 중복 키를 제거한다.
//
// 로봇 4대 x 지도 10개처럼 캐시 키에 로봇 ID가 섞여 들어가면 같은 지도를 4번
// 호출하게 된다. 지도 단위로 정규화하면 호출 수가 1/4로 줄어든다. 호출 수를
// 줄이는 것이 스로틀링보다 근본적인 개선이다.
func dedupe(keys []string) []string {
	seen := make(map[string]struct{}, len(keys))
	unique := make([]string, 0, len(keys))
	for _, key := range keys {
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		unique = append(unique, key)
	}
	return unique
}
```

- [ ] **Step 4: 테스트 통과 확인**

Run:
```bash
go test ./golang/resilience/outbound/... -run 'TestRefresher|TestDedupe' -v
```

Expected: PASS. `TestRefresher_키_중복제거로_호출수가_줄어든다`는 약 1초 걸린다.

- [ ] **Step 5: 커밋**

```bash
git add golang/resilience/outbound/refresher.go golang/resilience/outbound/refresher_test.go
git commit -m "feat: 키 중복 제거와 singleflight로 캐시 갱신 호출 수를 줄이는 Refresher 추가"
```

---

## Task 4: 버스트 vs 스로틀링 비교 (핵심 테스트)

이 계획의 결과물 중 리뷰 논의에 직접 답하는 부분이다.

**Files:**
- Modify: `golang/resilience/outbound/client_test.go` (테스트 함수 추가)

- [ ] **Step 1: 실패하는 테스트 작성**

`golang/resilience/outbound/client_test.go` 끝에 아래 테스트를 추가한다. `import` 블록에 `"context"`를 추가한다.

```go
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
			name: "정책 없음 - 동시에 발사하면 429가 발생한다",
			opts: Options{Mode: LimiterNone},
			expectReject: true,
		},
		{
			// 40회 / 2초 = 평균 20 QPS로 한도에 정확히 맞춘 설정이다.
			// 그런데도 주기 시작에 40개를 한꺼번에 통과시켜 순간 QPS가 튄다.
			// "1분에 1200회까지 가능하다"는 계산이 안전을 보장하지 않는 이유다.
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
```

- [ ] **Step 2: 테스트 실행**

Run:
```bash
go test ./golang/resilience/outbound/... -run TestClient_버스트 -v
```

Expected: PASS 3건. smooth 케이스는 약 4초, bursty는 약 2초 걸린다.

이 테스트는 Task 2·3에서 만든 코드만으로 통과한다. 새 구현이 필요 없는 검증 테스트이므로, 실패를 먼저 보는 대신 **출력 수치가 기대와 맞는지** 확인한다. `t.Logf` 출력이 아래와 비슷해야 한다.

```
정책_없음: 총 호출 40, 429 20, 최대 관측 QPS 40, 소요 ...ms
bursty:    총 호출 40, 429 20, 최대 관측 QPS 40, 소요 ...ms
smooth:    총 호출 40, 429 0,  최대 관측 QPS 10 또는 11, 소요 ~3.9s
```

- [ ] **Step 3: 실패 시 조치**

`bursty` 케이스에서 429가 0이면 limiter가 버스트를 허용하지 않은 것이다. `MaxExecutions`가 `requests`(40)와 같고 `Period`가 2초인지 확인한다. `MaxWaitTime`이 너무 짧아 `ratelimiter.ErrExceeded`가 나면 `got.Total`이 40보다 작아지므로, 그 경우 `MaxWaitTime`을 늘린다.

`smooth` 케이스에서 `got.MaxQPS`가 20을 넘으면 `MaxExecutions`/`Period`가 10/1초인지 확인한다.

- [ ] **Step 4: 전체 테스트 통과 확인**

Run:
```bash
go test ./golang/resilience/outbound/... -v
```

Expected: 전체 PASS. 총 소요 10초 내외.

- [ ] **Step 5: 커밋**

```bash
git add golang/resilience/outbound/client_test.go
git commit -m "test: 버스트/bursty/smooth 세 방식의 429 발생 여부 비교 테스트 추가"
```

---

## Task 5: 429 재시도 안전망 검증

**Files:**
- Modify: `golang/resilience/outbound/client_test.go` (테스트 함수 추가)

- [ ] **Step 1: 테스트 작성**

`golang/resilience/outbound/client_test.go` 끝에 아래를 추가한다.

```go
func TestClient_429를_받으면_RetryAfter만큼_기다린다(t *testing.T) {
	if testing.Short() {
		t.Skip("Retry-After 대기로 수 초가 걸린다")
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
```

- [ ] **Step 2: 테스트 실행**

Run:
```bash
go test ./golang/resilience/outbound/... -run TestClient_429 -v
```

Expected: PASS. 2~3초 걸린다.

- [ ] **Step 3: 실패 시 조치**

`stats.Failed`가 0이 아니면 재시도 횟수가 부족한 것이다. 재시도가 모두 같은 순간에 몰리면 다시 429를 받을 수 있다. `MaxRetries`를 8로 올린다.

- [ ] **Step 4: short 모드에서 건너뛰는지 확인**

Run:
```bash
go test ./golang/resilience/outbound/... -short -run TestClient_429 -v
```

Expected: `SKIP` 출력

- [ ] **Step 5: 커밋**

```bash
git add golang/resilience/outbound/client_test.go
git commit -m "test: 429 발생 시 Retry-After 기반 재시도로 복구되는지 검증"
```

---

## Task 6: 문서

**Files:**
- Create: `golang/resilience/outbound/README.md`
- Modify: `golang/resilience/README.md:7-31` (구조 트리)

- [ ] **Step 1: 실제 테스트 수치 수집**

Run:
```bash
go test ./golang/resilience/outbound/... -v 2>&1 | grep -E "총 호출|입력 키"
```

출력된 수치를 README의 표에 그대로 넣는다. 추정치를 쓰지 말고 실제 출력을 옮긴다.

- [ ] **Step 2: `golang/resilience/outbound/README.md` 작성**

아래 내용을 넣되, `<실측>` 자리에는 Step 1에서 얻은 실제 수치를 채운다.

````markdown
# 외부 API 버스트 호출 방지

캐시 갱신을 위해 외부 벤더 API를 주기적으로 호출할 때, 요청을 한꺼번에 발사해
rate limit(429)에 걸리는 문제와 그 해결책을 다룬다.

## 문제

- 벤더 제한: 20 QPS
- 호출량: 로봇 4대 x 지도 10개 = 40 호출/분
- 평균 QPS: **0.67** (한도의 3%)
- 순간 QPS: **최대 40** (동시 발사 시)

총량 문제가 아니라 **분산 문제**다. 40개를 몇 초에 걸쳐 흘려보내기만 해도 한도
안에 들어간다.

## 핵심: 평균 QPS를 맞추는 것으로는 부족하다

"20 QPS면 1분에 1200회까지 가능하다"는 계산은 평균만 본 것이다. 제약은 순간
QPS에 걸린다. 아래는 40개 요청을 20 QPS 한도(sliding window) 서버에 보낸 결과다.

| 방식 | 설정 | 평균 QPS | 429 | 최대 관측 QPS |
|---|---|---|---|---|
| 정책 없음 | — | — | <실측> | <실측> |
| bursty | 40회 / 2초 | 20 | <실측> | <실측> |
| smooth | 10회 / 1초 | 10 | <실측> | <실측> |

`bursty`는 **평균을 한도에 정확히 맞췄는데도 429를 받는다.** 주기 시작 시점에
40개를 지연 없이 통과시키기 때문이다. `smooth`는 `period / maxExecutions`
간격을 강제해 버스트 자체가 불가능하다.

## Best Practice 계층

| 계층 | 처방 | 효과 |
|---|---|---|
| ① 호출 수 줄이기 | 키 정규화·중복 제거, singleflight, TTL 연장 + 지터, 조건부 요청(ETag) | 40 → 10 |
| ② 남은 호출 평탄화 | smooth rate limiter(한도의 50~70%만), bulkhead, 클라이언트 1개에 정책 공유 | 순간 QPS 상한 보장 |
| ③ 그래도 429면 | `Retry-After` 존중 백오프, adaptive throttling, stale 값 서빙 | 실패 흡수 |
| ④ 관측 | 429 수, 스로틀 대기 시간, 주기 내 완료 여부 | 증설 시점 조기 감지 |
| ⑤ 다중 인스턴스 | 분산 limiter 또는 인스턴스 수로 rate 분할 | 프로세스별 limiter 합산 초과 방지 |

①이 가장 효과가 크고 비용이 싸다. ②만 적용하면 대상이 늘어날 때 대기 시간이
길어져 갱신 주기를 지키지 못하게 된다.

`4 x 10 = 40`이라는 곱셈이 나온다는 것은 **캐시 키에 로봇 ID가 들어갔다는
신호**다. 지도 데이터가 로봇에 종속적이지 않다면 키를 지도 단위로 정규화해
호출을 1/4로 줄일 수 있다.

## 구성

```go
client := outbound.NewClient(outbound.Options{
    Mode:           outbound.LimiterSmooth,
    MaxExecutions:  10,              // 한도 20 QPS의 50%
    Period:         time.Second,
    MaxWaitTime:    30 * time.Second,
    MaxConcurrency: 5,
    MaxRetries:     3,
})
```

정책은 `http.RoundTripper` 층에 붙으므로 `net/http`, resty 등 클라이언트 종류와
무관하다. resty를 쓴다면 `client.SetTransport(...)`로 같은 트랜스포트를 넣으면
된다.

**정책 순서가 중요하다.** `failsafe`는 첫 인자가 가장 바깥 정책이므로
`retry → rateLimiter → bulkhead` 순으로 넘긴다. 이렇게 해야 재시도 요청도 rate
limiter를 다시 통과한다. 순서를 뒤집으면 재시도가 limiter를 우회해 429 폭풍을
만든다.

**클라이언트 하나를 공유해야 한다.** rate limit은 보통 계정 단위이므로, 호출
지점마다 limiter를 따로 두면 합산해서 다시 한도를 넘는다.

### bulkhead가 별도로 필요한 이유

rate limiter는 실행 *시작 시점*만 제어하고 in-flight 수는 제한하지 않는다.
응답이 느려지면 동시 연결이 쌓인다. 벤더가 동시 연결 수도 제한하는 경우가 많아
별도 층이 필요하다. (이 예제는 구성만 제공하고 전용 테스트는 없다. 로컬
`httptest` 서버는 응답이 즉시 돌아와 동시 실행 수가 쌓이지 않는다.)

## resty 내장 rate limiter를 쓰지 않는 이유

- **resty v2** (`v2.8.0`+): `SetRateLimiter`의 인터페이스가 `Allow() bool`
  하나뿐이다. 논블로킹으로 초과분을 **즉시 버리고**, `wrapNoRetryErr`로 감싸서
  재시도조차 하지 않는다. 캐시 갱신처럼 "늦게라도 다 보내야 하는" 요청에는
  부적합하다.
- **resty v3**: `SetRateLimiter(resty.NewRateLimitTokenBucket(rate, burst))`가
  토큰이 생길 때까지 **대기**한다. 시맨틱은 적합하나 아직 `v3.0.0-rc.3`으로
  정식 릴리스 전이다.

`failsafehttp`는 `RoundTripper` 층에서 동작하므로 클라이언트 라이브러리에
종속되지 않는다.

## 테스트

```bash
# 전체 (약 10초)
go test ./golang/resilience/outbound/... -v

# Retry-After 대기 테스트 제외
go test ./golang/resilience/outbound/... -short -v
```

가짜 벤더 서버는 최근 1초 **sliding window**로 요청을 센다. 고정 윈도우는
경계에서 최대 2배를 허용해 클라이언트 주기와 우연히 정렬되면 버스트 설정으로도
429가 나지 않는다. 그러면 테스트 통과가 처방 덕인지 정렬 운인지 구분할 수 없다.

## 다음 단계

- **분산 환경**: 프로세스별 limiter는 인스턴스 수만큼 합산되어 한도를 넘는다.
  `../ratelimit/redis_limiter.go`의 Redis 기반 분산 limiter를 참고한다.
- **adaptive throttling**: 429 비율을 보고 rate를 자동으로 낮추는 방식.
  `failsafe-go`의 `adaptivethrottler` 패키지.
- **호출 자체 없애기**: 벤더에 배치 API, webhook/push, ETag 조건부 요청 지원
  여부를 문의한다. 조건부 요청의 304가 rate limit에 계산되는지도 확인이 필요하다.
````

- [ ] **Step 3: `golang/resilience/README.md` 구조 트리 수정**

`circuitbreaker/` 블록과 `└── README.md` 사이에 아래를 추가하고, 기존
`├── circuitbreaker/`의 마지막 항목 뒤 트리 문자를 맞춘다.

```
├── outbound/
│   ├── client.go                  # rate limiter + bulkhead + 재시도 조립
│   ├── client_test.go             # 버스트 vs bursty vs smooth 429 비교
│   ├── refresher.go               # 키 중복 제거 + singleflight 캐시 갱신
│   ├── refresher_test.go
│   └── fakeserver_test.go         # sliding window rate limit 가짜 서버
```

같은 파일의 "사용 라이브러리" 섹션에 아래 줄을 추가한다.

```markdown
- [golang.org/x/sync/singleflight](https://pkg.go.dev/golang.org/x/sync/singleflight) - 진행 중인 중복 호출 합치기
```

- [ ] **Step 4: 인코딩과 빌드 확인**

Run:
```bash
file -I golang/resilience/outbound/README.md golang/resilience/README.md
go build ./golang/resilience/... && go vet ./golang/resilience/outbound/...
```

Expected: 두 파일 모두 `charset=utf-8`. 빌드와 vet 모두 출력 없이 성공.

- [ ] **Step 5: 커밋**

```bash
git add golang/resilience/outbound/README.md golang/resilience/README.md
git commit -m "docs: outbound 예제 README와 resilience 구조 트리 추가"
```

---

## Task 7: 최종 검증

**Files:** 없음 (검증만)

- [ ] **Step 1: 포맷 확인**

Run:
```bash
gofmt -l golang/resilience/outbound/
```

Expected: 출력 없음. 출력이 있으면 `gofmt -w golang/resilience/outbound/` 실행 후 커밋한다.

- [ ] **Step 2: 전체 테스트 3회 반복 (flaky 확인)**

Run:
```bash
go test ./golang/resilience/outbound/... -count=3
```

Expected: `ok`. 시간 의존 테스트이므로 3회 모두 통과해야 안심할 수 있다.

- [ ] **Step 3: 저장소 전체 빌드 영향 확인**

Run:
```bash
go build ./... && go vet ./golang/resilience/...
```

Expected: 출력 없이 성공. `go mod tidy`로 인한 다른 패키지 영향이 없는지 본다.

- [ ] **Step 4: 변경 범위 확인**

Run:
```bash
git status --short
git log --oneline master..HEAD
```

Expected: 작업 시작 전부터 있던 변경(`cloud/docker/redis/Makefile` 수정, 로그 파일 등)만 uncommitted로 남아 있어야 한다. 커밋 목록은 스펙 1건 + 구현 6건.

- [ ] **Step 5: PR 생성**

```bash
git push -u origin feature/outbound-burst-ratelimit

gh pr create --title "feat: 외부 API 버스트 호출 방지 예제 추가" --body "$(cat <<'EOF'
## Summary
캐시 갱신 시 외부 벤더 API를 버스트로 호출해 429를 받는 문제를 재현하고, smooth rate limiter로 해결되는 것을 테스트로 증명하는 예제입니다.

핵심은 **평균 QPS를 한도에 맞추는 것만으로는 부족하다**는 점입니다. 40회/2초(평균 20 QPS)로 설정한 bursty limiter는 한도에 정확히 맞췄는데도 429를 받습니다. 주기 시작에 40개를 한꺼번에 통과시키기 때문입니다.

- `failsafehttp.NewRoundTripper`로 정책을 트랜스포트 층에 조립 (클라이언트 라이브러리 무관)
- `ratelimiter` smooth/bursty 비교, `bulkhead`, `Retry-After` 존중 재시도
- 키 중복 제거 + singleflight로 호출 수 자체를 40 → 10으로 감소
- sliding window rate limit을 흉내내는 `httptest` 서버로 429 수와 최대 관측 QPS 측정

신규 의존성 없음. `failsafe-go`가 indirect에서 direct로 승격됩니다.

## Test plan
- [ ] `go test ./golang/resilience/outbound/... -v` 전체 통과
- [ ] `go test ./golang/resilience/outbound/... -count=3` 반복 통과 (flaky 확인)
- [ ] `go build ./...` 저장소 전체 빌드 성공

🤖 Generated with [Claude Code](https://claude.com/claude-code)
EOF
)"
```

---

## 자체 검토 메모

**스펙 커버리지**

| 스펙 항목 | 담당 Task |
|---|---|
| `client.go` — `LimiterMode`, `Options`, 정책 순서 | Task 2 |
| `refresher.go` — 중복 제거, singleflight, stale 서빙 | Task 3 |
| `fakeserver_test.go` — sliding window, 429 + Retry-After, 통계 | Task 1 |
| 비교 테스트 (정책 없음 / bursty / smooth), `MaxRetries: 0` | Task 4 |
| 키 중복 제거 테스트 | Task 3 |
| 429 재시도 테스트, `testing.Short()` 가드 | Task 5 |
| README 2건 | Task 6 |
| `go mod tidy`로 failsafe-go direct 승격 | Task 2 Step 4 |
| bulkhead 전용 테스트 없음 (의도) | Task 6 README에 명시 |

**타입 일관성**: `Options`(Task 2)의 필드명이 Task 3·4·5의 호출부와 일치한다. `Stats`(Task 3)의 `Requested`/`Unique`/`Failed`/`Elapsed`가 Task 3·4·5의 단정과 일치한다. `serverStats`(Task 1)의 `Total`/`Rejected`/`MaxQPS`가 Task 2·3·4·5의 단정과 일치한다. `newRateLimitedServer`/`stats`/`dedupe`/`NewRefresher`/`NewClient` 이름이 전 Task에서 동일하다.

**Task 4·5가 TDD 형태가 아닌 이유**: 두 Task는 새 구현을 요구하지 않는 검증 테스트다. 실패를 먼저 보는 대신 출력 수치가 기대와 맞는지 확인하고, 어긋날 때의 조치를 Step 3에 적었다.
