# Resilience Patterns in Go

Go에서 서비스 안정성을 위한 Resilience 패턴 예제 코드 모음.

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
├── circuitbreaker/
│   ├── gobreaker_example.go       # sony/gobreaker v2 기본 사용법
│   ├── gobreaker_example_test.go
│   ├── gobreaker_http.go          # HTTP 클라이언트를 Circuit Breaker로 래핑
│   ├── gobreaker_http_test.go
│   ├── failsafe_example.go        # failsafe-go Count-based/Time-based Circuit Breaker
│   ├── failsafe_example_test.go
│   ├── failsafe_composed.go       # Fallback + Retry + Circuit Breaker 정책 조합
│   └── failsafe_composed_test.go
└── README.md
```

## 사용 라이브러리

- [sony/gobreaker v2](https://github.com/sony/gobreaker) - 경량 Circuit Breaker
- [failsafe-go](https://github.com/failsafe-go/failsafe-go) - 통합 Resilience 라이브러리

## 테스트

```bash
# 단위 테스트만 실행
go test -short ./golang/resilience/...

# 전체 테스트 (Docker 필요 - Redis 통합 테스트 포함)
go test ./golang/resilience/...
```

## 관련 블로그
- [Go Rate Limiting 완벽 가이드](https://blog-v2.advenoh.pe.kr)
