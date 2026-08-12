# 시간 의존 코드 테스트하기

`time.Now()`를 직접 호출하는 코드는 실행 시점마다 결과가 달라 테스트가 불안정하다(flaky).
이 예제는 시간 의존 코드를 테스트하는 4가지 패턴을 점진적으로 보여준다.

## 패턴 요약

| # | 패턴 | 예제 파일 | 핵심 아이디어 |
|---|------|----------|--------------|
| 1 | 시간을 파라미터로 받기 | `coupon.go` | `IsExpiredAt(now time.Time)` — 순수 함수로 만들어 테스트에서 임의 시각 전달 |
| 2 | nowFunc 필드 주입 | `order.go` | 구조체에 `nowFunc func() time.Time` 필드(기본값 `time.Now`), 테스트에서 고정 시간 주입 |
| 3 | Clock 인터페이스 주입 | `cache.go` | `clockwork.Clock` 주입, 테스트에서 `fakeClock.Advance()`로 시간 진행 |
| 4 | testing/synctest | `cache_synctest_test.go` | 코드 수정 없이 가상 시간 버블에서 테스트 (Go 1.25+) |

## 어떤 패턴을 언제 쓸까

| 상황 | 추천 패턴 |
|------|----------|
| 특정 시각 기준 판단만 필요 (만료 검사 등) | 패턴 1 (파라미터 전달) |
| 구조체가 `time.Now()`를 저장/기록 | 패턴 2 (nowFunc 주입) — 실무의 80%는 이걸로 충분 |
| 시간이 "흘러야" 검증되는 로직 (TTL, 재시도, 스케줄러) | 패턴 3 (Clock 인터페이스, clockwork) |
| 동시성 + 시간 조합, 기존 코드 수정 불가 | 패턴 4 (synctest, Go 1.25+) |
| 주입이 불가능한 레거시 코드 | `assert.WithinDuration`으로 오차 허용 비교 (`order_test.go` 참고) |

## 트레이드오프

- **패턴 1 (파라미터 전달)**: 가장 단순하고 순수하다. 다만 깊은 콜스택이면 모든 호출자가 시간을 넘겨야 한다.
- **패턴 2 (nowFunc 주입)**: 인터페이스 없이 가볍고 관용적이다. `Now()`만 대체 가능하다. 필드가 unexported라 같은 패키지 테스트(white-box)에서만 주입할 수 있다는 점에 주의.
- **패턴 3 (Clock 인터페이스)**: `Sleep`/`Ticker`/`After`까지 제어하고 시간 진행을 시뮬레이션한다. 의존성과 인터페이스 추가 비용이 든다.
- **패턴 4 (synctest)**: 표준 라이브러리이고 코드 수정이 필요 없다 — 시계 주입 설계가 없는 코드도 테스트할 수 있다. 가상 시간이 결정론적이라 `now == expiresAt` 같은 정확한 경계도 flaky 없이 검증된다. 동시성+시간 조합에 특화되어 있고 Go 1.25 이상이 필요하다.

## 경계값 처리 주의

이 예제의 두 구현은 만료 경계를 서로 다르게 처리한다. 의도적인 차이다.

- 쿠폰(패턴 1): `now.After(expiresAt)` — 만료 시각과 정확히 같은 순간은 **아직 유효**
- TTL 캐시(패턴 3/4): `!now.Before(expiresAt)` — 만료 시각 **정각부터 miss**

둘 다 유효한 선택이지만, 시간 비교 코드는 경계 포함 여부를 테스트로 고정해 두는 것이 좋다
(`coupon_test.go`의 "만료 시각과 동일" 케이스 참고).

## 실행

```bash
go test ./golang/testing/clock/... -v
```

## 관련 예제

- [`golang/go1_25/synctest_test.go`](../../go1_25/synctest_test.go) — synctest 기초 (goroutine + 가상 시간)
- [`golang/time/`](../../time/) — time 패키지 API 사용법 (parse, zone, duration)
