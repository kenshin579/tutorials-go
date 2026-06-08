package patterns

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestRateLimitWithTicker - time.Ticker로 rate limiting
func TestRateLimitWithTicker(t *testing.T) {
	rate := time.NewTicker(20 * time.Millisecond) // 50 req/sec
	defer rate.Stop()                             // 내부 goroutine leak 방지

	start := time.Now()
	for i := range 5 {
		<-rate.C // tick을 기다림: 각 작업 사이 최소 20ms 간격 보장
		_ = i    // 작업 수행
	}

	elapsed := time.Since(start)
	t.Logf("5개 작업 소요 시간: %v", elapsed)
	assert.GreaterOrEqual(t, elapsed, 80*time.Millisecond) // 최소 4 tick 대기
}

// TestBurstRateLimit - 버스트 허용 rate limiter
func TestBurstRateLimit(t *testing.T) {
	// 버스트 크기 3, 이후 20ms 간격
	burstLimit := make(chan time.Time, 3)

	// 초기 토큰 채우기: 처음 3개 요청은 대기 없이 즉시 처리됨
	for range 3 {
		burstLimit <- time.Now()
	}

	// 토큰 보충 goroutine: 20ms마다 토큰 1개 추가
	go func() {
		ticker := time.NewTicker(20 * time.Millisecond)
		defer ticker.Stop()
		for t := range ticker.C {
			select {
			case burstLimit <- t:
			default: // 버퍼 가득 시 토큰 버림 → 토큰 무한 누적 방지
			}
		}
	}()

	// 처음 3개는 즉시, 이후는 rate limit 적용
	start := time.Now()
	for range 5 {
		<-burstLimit
	}

	elapsed := time.Since(start)
	t.Logf("5개 작업 (burst 3) 소요 시간: %v", elapsed)
	assert.GreaterOrEqual(t, elapsed, 30*time.Millisecond)
}
