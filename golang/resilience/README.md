# Resilience Patterns in Go

Go에서 서비스 안정성을 위한 Resilience 패턴 예제 코드 모음.

## Circuit Breaker (`circuitbreaker/`)

외부 서비스 장애 시 Cascading Failure를 방지하는 Circuit Breaker 패턴 구현.

### 구현 내용

| 파일 | 설명 |
|---|---|
| `gobreaker_example.go` | sony/gobreaker v2 기본 사용법 |
| `gobreaker_http.go` | HTTP 클라이언트를 Circuit Breaker로 래핑 |
| `failsafe_example.go` | failsafe-go Count-based/Time-based Circuit Breaker |
| `failsafe_composed.go` | Fallback + Retry + Circuit Breaker 정책 조합 |

### 사용 라이브러리

- [sony/gobreaker v2](https://github.com/sony/gobreaker) - 경량 Circuit Breaker
- [failsafe-go](https://github.com/failsafe-go/failsafe-go) - 통합 Resilience 라이브러리

### 테스트 실행

```bash
go test -v ./golang/resilience/circuitbreaker/
```
