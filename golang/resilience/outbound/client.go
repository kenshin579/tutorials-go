// Package outbound는 외부 벤더 API를 rate limit 안에서 안전하게 호출하는
// 예제를 담는다.
//
// 핵심은 평균 QPS가 아니라 순간 QPS를 제한하는 것이다. 1분에 40회는 평균
// 0.67 QPS지만, 40개를 동시에 발사하면 순간 40 QPS가 되어 20 QPS 한도를
// 넘는다.
package outbound

import (
	"fmt"
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
	Mode LimiterMode // rate limiter 방식

	// MaxExecutions는 주기당 허용 실행 수다. Mode가 LimiterNone이 아니면
	// 0보다 커야 한다 (그렇지 않으면 NewClient가 panic한다).
	MaxExecutions int

	// Period는 MaxExecutions를 적용할 주기다. Mode가 LimiterNone이 아니면
	// 0보다 커야 한다 (그렇지 않으면 NewClient가 panic한다).
	Period time.Duration

	// MaxWaitTime은 permit(rate limiter) 또는 slot(bulkhead) 대기 상한이다.
	// 초과 시 rate limiter는 ratelimiter.ErrExceeded, bulkhead는
	// bulkhead.ErrFull을 반환한다.
	//
	// 정책이 Limiter(Bulkhead(func)) 순으로 중첩되므로 한 요청이 최악의
	// 경우 limiter에서 MaxWaitTime, bulkhead에서 다시 MaxWaitTime을 기다릴
	// 수 있다. 즉 꼬리 지연(worst-case latency)은 MaxWaitTime의 최대 2배다.
	MaxWaitTime time.Duration

	MaxConcurrency int // 동시 실행 수 상한. 0이면 bulkhead 없음
	MaxRetries     int // 재시도 횟수. 0이면 재시도 없음

	// Base는 실제 요청을 보낼 하위 RoundTripper다. nil이면
	// http.DefaultTransport를 사용한다.
	Base http.RoundTripper
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
//
// MaxExecutions/Period가 0 이하이면 panic한다. int를 uint로 변환하기 전에
// 걸러내지 않으면 0은 division-by-zero panic으로, 음수는 uint로 wrap되어
// 거의 무한대의 실행을 허용하는 "스로틀링이 조용히 사라지는" 결과로 이어진다.
func newRateLimiter(opts Options) ratelimiter.RateLimiter[*http.Response] {
	if opts.Mode == LimiterNone {
		return nil
	}
	if opts.MaxExecutions <= 0 {
		panic(fmt.Sprintf("outbound: MaxExecutions는 0보다 커야 합니다 (got %d)", opts.MaxExecutions))
	}
	if opts.Period <= 0 {
		panic(fmt.Sprintf("outbound: Period는 0보다 커야 합니다 (got %s)", opts.Period))
	}

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
