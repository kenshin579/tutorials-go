# 시간 의존 코드 테스트 예제 설계

- 날짜: 2026-08-11
- 위치: `golang/testing/clock/` (신규 폴더, package `clock`)
- 목적: 구조체/함수에 `time.Time` 의존성이 있을 때 테스트하는 4가지 패턴을 점진적으로 보여주는 튜토리얼 예제. 블로그(blog-v2) 글 소재로도 활용.

## 배경

`time.Now()`를 직접 호출하는 코드는 실행 시점마다 결과가 달라 테스트가 불안정하다. 이 예제는 "왜 다음 단계가 필요한가"를 보여주는 점진적 구성으로 4가지 해법을 다룬다.

기존 예제와의 관계:

- `golang/time/` — time 패키지 API 사용법 (parse, zone, duration). 이 예제와 주제가 다름.
- `golang/testing/` — testify, benchmark 등 테스트 기법. 이 예제는 그 하위 폴더로 들어감.
- `golang/go1_25/synctest_test.go` — synctest 기초 (goroutine + 가상 시간). README에서 상호 링크하고 중복하지 않음.

## 파일 구성

하나의 일관된 도메인(쿠폰 만료 / 주문 생성 / TTL 캐시)으로 네 패턴을 관통한다.

| 파일 | 패턴 | 내용 |
|------|------|------|
| `coupon.go` / `coupon_test.go` | ① 시간을 파라미터로 받기 | `func (c Coupon) IsExpiredAt(now time.Time) bool` — 순수 함수. 테이블 드리븐 테스트로 경계값(만료 직전/직후/동일 시각) 검증 |
| `order.go` / `order_test.go` | ② nowFunc 필드 주입 | `OrderService` 구조체에 `nowFunc func() time.Time` 필드(기본값 `time.Now`). 테스트에서 고정 시간 주입. 주입 없이 테스트할 때의 차선책으로 `assert.WithinDuration` 예제 포함 |
| `cache.go` / `cache_test.go` | ③ Clock 인터페이스 주입 | `TTLCache`가 `clockwork.Clock`을 받음. 테스트에서 `fakeClock.Advance(ttl)`로 시간을 진행시켜 만료 검증 — 시간이 "흘러야" 하는 로직은 ②로 불가능함을 보여줌 |
| `cache_synctest_test.go` | ④ testing/synctest (Go 1.25+) | 같은 `TTLCache`를 clock 주입 없이 실제 `time.Sleep`으로 synctest 버블 안에서 테스트. 코드 수정 없이 시간 의존 테스트가 가능한 표준 라이브러리 접근법 |
| `README.md` | — | 각 패턴의 트레이드오프 표, 어떤 상황에 어떤 패턴을 쓰는지 가이드, go1_25 synctest 예제 링크 |

## 블로그 글 초안 (추가 산출물)

- 위치: `blog-v2.advenoh.pe.kr/docs/start/<슬러그>/index.md` (blog-v2 저장소, `contents/`에 바로 넣지 않고 초안 폴더에 작성)
- 형식: 기존 초안 컨벤션 준수 — frontmatter(title, description, date, update, tags), 한국어 문어체("~이다"), 참고 자료 블록
- 슬러그(가안): `go-시간-의존-코드-테스트하기-testing-time-dependent-code-in-go`
- 내용: 4가지 패턴을 점진적 서사로 소개, tutorials-go 예제 코드 인용 및 저장소 링크
- 작성 시점: 코드 구현·테스트 통과 완료 후 (코드가 글의 근거이므로)

## 패턴별 트레이드오프 (README에 수록)

1. **파라미터 전달**: 가장 단순, 순수 함수. 호출자마다 시간을 넘겨야 해 깊은 콜스택엔 부적합.
2. **nowFunc 필드 주입**: 인터페이스 없이 가볍고 관용적. 실무의 80%는 이걸로 충분. `Now()`만 필요할 때 한정.
3. **Clock 인터페이스 (clockwork)**: `Sleep`/`Ticker`/`After`까지 제어, 시간 진행(Advance) 시뮬레이션 가능. 의존성/인터페이스 추가 비용.
4. **synctest**: 코드 수정 불필요, 표준 라이브러리. 동시성+시간 조합에 특화, Go 1.25+ 필요.

## 기술 결정

- **의존성**: `github.com/jonboulle/clockwork` v0.4.0 — 이미 go.mod에 indirect로 존재, direct로 승격 (`go mod tidy`).
- **synctest 위치**: `go1_25/`가 아닌 `testing/clock/`에 둔다. "같은 코드를 4가지 방법으로 테스트"라는 튜토리얼 자체 완결성이 버전별 폴더 컨벤션보다 우선. go.mod가 go 1.26이므로 빌드 문제 없음.
- **테스트 컨벤션**: `TestXxx_설명` 네이밍, 테이블 드리븐, testify/assert (`.claude/rules/testing.md` 준수).

## 범위 제외 (YAGNI)

- mockery로 Clock 목 생성 — clockwork의 fake clock으로 충분
- benbjohnson/clock 등 대안 라이브러리 비교
- 타임존 이슈 — `golang/time/time_zone_test.go`에 이미 존재

## 검증 기준

- `go test ./golang/testing/clock/...` 전체 통과
- `go vet ./golang/testing/clock/...` 클린
- `gofmt` 적용 확인
- 블로그 초안: frontmatter 형식 유효, 본문 코드 스니펫이 실제 예제 코드와 일치, UTF-8 인코딩 확인
