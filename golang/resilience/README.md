# Resilience Patterns in Go

Rate Limiting과 Retry 패턴 예제 코드.

## 구조

```
resilience/
├── ratelimit/
│   ├── token_bucket.go           # x/time/rate 기반 Token Bucket
│   ├── token_bucket_test.go
│   ├── middleware.go             # Echo HTTP 미들웨어 (IP별 Rate Limiting)
│   ├── middleware_test.go
│   ├── redis_limiter.go          # go-redis/redis_rate 분산 Rate Limiting
│   └── redis_limiter_test.go     # testcontainers-go 통합 테스트
├── retry/
│   ├── backoff.go                # cenkalti/backoff/v5 Exponential Backoff
│   ├── backoff_test.go
│   ├── retry.go                  # avast/retry-go/v4 Jitter Retry
│   └── retry_test.go
└── README.md
```

## 테스트

```bash
# 단위 테스트만 실행
go test -short ./golang/resilience/...

# 전체 테스트 (Docker 필요 - Redis 통합 테스트 포함)
go test ./golang/resilience/...
```

## 관련 블로그
- [Go Rate Limiting 완벽 가이드](https://blog-v2.advenoh.pe.kr)
