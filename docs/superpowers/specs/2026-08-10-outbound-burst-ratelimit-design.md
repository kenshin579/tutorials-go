# 외부 API 버스트 호출 방지 설계

작성일: 2026-08-10

## 배경

캐시 갱신을 위해 외부 벤더 API(Gausium)를 1분 주기로 호출하는 구조에서, 코드 리뷰 중 rate limit 위반 우려가 제기되었다.

- 벤더 제한: 20 QPS
- 현재 호출량: 로봇 4대 × 지도 10개 = 40 호출/분
- 평균 QPS: 0.67 (한도의 3%)
- 문제: 40개를 동시에 발사하면 **순간 30~40 QPS**가 되어 뒤쪽 요청이 429를 받는다

즉 총량 문제가 아니라 **분산(버스트) 문제**다. 40개를 몇 초에 걸쳐 흘려보내기만 해도 한도 안에 들어가고, 신선도는 1분 주기 대비 그 몇 초가 밀리는 것이 전부다.

로봇 대수가 늘면 호출량이 비례해 늘어나므로, 재현 가능한 형태로 문제와 처방을 검증하는 예제가 필요하다.

## 목표

이 저장소에 **문제를 재현하고 처방의 효과를 숫자로 증명하는** 실행 가능한 예제를 만든다.

검증할 명제 3개:

1. 버스트로 발사하면 실제로 429가 발생한다
2. **평균 QPS를 한도에 맞추는 것만으로는 부족하다** — 평균이 같아도 버스트를 허용하는 limiter는 여전히 429를 받는다
3. smooth(균등 간격) limiter를 쓰면 429가 사라진다

2번이 이번 리뷰 논의의 핵심이다. "1분에 1200회까지 가능하다"는 계산은 평균만 본 것이고, 제약은 순간 QPS에 걸린다.

## 비목표

- 프로덕션 코드에 대한 직접 수정
- 분산 환경(다중 인스턴스) limiter 구현 — `ratelimit/redis_limiter.go`에 이미 있음
- adaptive throttling — `failsafe-go`의 `adaptivethrottler`로 가능하나 이번 범위 밖
- 메트릭 익스포터 연동

README에 "다음 단계"로만 언급한다.

## Best Practice 계층

| 계층 | 처방 | 효과 |
|---|---|---|
| ① 호출 수 줄이기 | 키 정규화·중복 제거, singleflight, TTL 연장 + 지터, 조건부 요청(ETag) | 40 → 10 |
| ② 남은 호출 평탄화 | smooth rate limiter(한도의 50~70%만), bulkhead, 클라이언트 1개에 정책 공유 | 순간 QPS 상한 보장 |
| ③ 그래도 429면 | `Retry-After` 존중 백오프, adaptive throttling, stale 값 서빙 | 실패 흡수 |
| ④ 관측 | 429 수, 스로틀 대기 시간, 주기 내 완료 여부 | 증설 시점 조기 감지 |
| ⑤ 다중 인스턴스 | 분산 limiter 또는 인스턴스 수로 rate 분할 | 프로세스별 limiter 합산 초과 방지 |

①이 가장 효과가 크고 비용이 싸다. ②만 적용하면 로봇이 늘 때 대기 시간이 길어져 1분 주기를 지키지 못하게 된다.

## 라이브러리 선택

**신규 의존성 추가 없음.** `failsafe-go v0.9.6`이 이미 저장소에 있고, 필요한 기능이 모두 들어 있다.

### 사용할 API (소스 확인 완료)

- `failsafehttp.NewRoundTripper(inner http.RoundTripper, policies ...failsafe.Policy[*http.Response]) http.RoundTripper`
  `http.RoundTripper`를 감싸므로 HTTP 클라이언트 종류와 무관하다. `http.Client`, resty(`SetTransport`) 모두 동일하게 적용된다.

- `failsafehttp.NewRetryPolicyBuilder() retrypolicy.Builder[*http.Response]`
  429와 대부분의 5xx를 재시도하고, `Retry-After` 헤더를 기본으로 존중한다(`failsafehttp.DelayFunc`). 직접 구현이 불필요하다.

- `ratelimiter.NewSmoothBuilder[R](maxExecutions uint, period time.Duration) Builder[R]`
  실행 간격을 `period / maxExecutions`로 균등하게 강제한다. 토큰을 모아두지 않으므로 버스트가 구조적으로 불가능하다. `WithMaxWaitTime(d)`로 무한 대기를 막는다.

- `ratelimiter.NewBurstyBuilder[R](maxExecutions uint, period time.Duration) Builder[R]`
  주기당 `maxExecutions`개를 허용하되, 주기 시작 시점에 **지연 없이 전부 통과시킨다**. 대조군으로 사용한다. 문서 인용: "Executions are performed with no delay up until the maxExecutions are reached for the current period."

- `bulkhead.NewBuilder[R](maxConcurrency uint) Builder[R]`
  동시 실행 수를 제한한다. smooth limiter는 실행 *시작 시점*만 제어하고 in-flight 수는 제한하지 않으므로, 응답이 느려지면 동시 연결이 쌓인다. 벤더가 동시 연결 수도 제한하는 경우가 많아 별도 층이 필요하다.

### `x/time/rate`를 대조군으로 쓰지 않는 이유

`x/time/rate`는 토큰 버킷이라 설계상 `burst`만큼 버스트를 **허용**한다(`burst=1`로 두면 균등 간격에 가까워진다). 대조군으로 적합한 성질이지만, failsafe 정책 체인 밖에 있어 `RoundTripper` 조립 방식이 케이스마다 달라진다. 비교 대상이 "limiter 종류"가 아니라 "조립 방식"으로 오염된다.

`ratelimiter.NewBurstyBuilder`가 동일한 성질을 같은 정책 체인 안에서 제공하므로, **smooth vs bursty를 같은 조건에서 비교**할 수 있다. 대조군은 이쪽을 쓴다.

### resty 내장 rate limiter를 쓰지 않는 이유

조사 결과를 기록해 둔다.

- **resty v2** (`v2.9.0`+): `SetRateLimiter(RateLimiter)`의 인터페이스가 `Allow() bool` 하나뿐이다. 논블로킹으로 초과분을 **즉시 버리고**, `wrapNoRetryErr`로 감싸서 재시도조차 하지 않는다. 캐시 갱신처럼 "늦게라도 다 보내야 하는" 요청에는 부적합하다. 참고로 이 저장소는 `v2.7.0`이라 이 API 자체가 없다.
- **resty v3**: `SetRateLimiter(resty.NewRateLimitTokenBucket(rate, burst))`가 토큰이 생길 때까지 **대기**한다. 시맨틱은 적합하나 아직 `v3.0.0-rc.3`으로 정식 릴리스 전이다.
- **결론**: `failsafehttp`는 `RoundTripper` 층에서 동작하므로 resty를 쓰든 안 쓰든 그대로 적용된다. 클라이언트 라이브러리에 종속되지 않는 쪽을 택한다.

## 구조

**위치**: `golang/resilience/outbound/`

기존 `golang/resilience/ratelimit/`은 limiter를 하나씩 보여주는 예제다. 이 예제는 여러 정책을 조합해 외부 API를 안전하게 호출하는 통합 예제이므로 분리한다.

```
golang/resilience/outbound/
├── client.go             # 정책을 조합한 http.Client 생성
├── client_test.go        # 정책 없음 vs bursty vs smooth 비교
├── refresher.go          # 키 중복 제거 + singleflight + 병렬 갱신
├── refresher_test.go
├── fakeserver_test.go    # QPS 제한을 흉내내는 httptest 서버
└── README.md
```

### `client.go`

```go
// LimiterMode는 rate limiter의 버스트 허용 방식을 정한다.
type LimiterMode int

const (
    LimiterNone   LimiterMode = iota // rate limiting 없음
    LimiterBursty                    // 주기 시작에 maxExecutions개를 한꺼번에 허용
    LimiterSmooth                    // period/maxExecutions 간격으로 균등 분산
)

type Options struct {
    Mode           LimiterMode
    MaxExecutions  int           // 주기당 허용 실행 수
    Period         time.Duration // 주기
    MaxWaitTime    time.Duration // limiter 대기 상한 (초과 시 ratelimiter.ErrExceeded)
    MaxConcurrency int           // 0이면 bulkhead 없음
    MaxRetries     int           // 0이면 재시도 없음
    Base           http.RoundTripper
}

func NewClient(opts Options) *http.Client
```

**정책 순서가 중요하다.** `failsafe`는 첫 인자가 가장 바깥 정책이므로 `retry` → `rateLimiter` → `bulkhead` 순으로 넘긴다. 이렇게 하면 재시도 요청도 rate limiter를 **다시** 통과한다. 순서를 뒤집어 limiter가 바깥에 오면 재시도가 limiter를 우회해 429 폭풍을 만든다. 이 이유를 코드 주석으로 남긴다.

### `refresher.go`

```go
type Stats struct {
    Requested   int           // 호출 요청된 키 개수 (중복 포함)
    Fetched     int           // 실제 서버 호출 수
    Failed      int
    Elapsed     time.Duration
}

func (r *Refresher) RefreshAll(ctx context.Context, keys []string) Stats
```

처리 순서:

1. 키 중복 제거 — `4 robots × 10 maps = 40`이 나오는 것은 캐시 키에 로봇 ID가 들어갔다는 신호다. 지도 단위로 정규화하면 10개가 된다.
2. `singleflight.Group`으로 in-flight 중복 호출 합치기
3. `errgroup`으로 병렬 호출 (동시성은 bulkhead가 제한)
4. 실패 시 기존 캐시 값을 유지 (stale 서빙) — 캐시 갱신 실패가 서비스 응답 실패로 번지지 않게 한다

### `fakeserver_test.go`

`httptest.Server`로 벤더 rate limit을 흉내낸다.

- 최근 1초 **sliding window**로 요청 수를 세고, 한도 초과분에 `429 + Retry-After: 1` 반환
- 구현: 요청 도착 시각을 슬라이스에 기록하고 1초 이전 항목을 버린 뒤 남은 개수로 판정
- 리포트 항목: 총 호출 수, 429 응답 수, **관측된 최대 초당 요청 수**

**왜 sliding window인가.** 고정 윈도우는 경계에서 최대 2배를 허용하므로, 클라이언트 주기와 서버 윈도우가 우연히 정렬되면 버스트 설정으로도 429가 나지 않는다. 그러면 테스트 통과가 처방 덕인지 정렬 운인지 구분할 수 없다. sliding window는 정렬에 무관하게 순간 QPS를 판정하므로 결과가 결정론적이다. 실제 API 게이트웨이도 sliding window나 GCRA를 쓰는 경우가 많아 보수적인 모델이기도 하다.

## 테스트

프로젝트 컨벤션에 따라 테이블 드리븐 + `testify/assert`, 테스트명은 `TestXxx_설명` 형식.

### `TestClient_버스트_요청과_스로틀링_비교`

서버 한도 20/s(sliding window), 요청 40개 동시 발사.

| 케이스 | 설정 | 평균 QPS | 기대 |
|---|---|---|---|
| 정책 없음 | `LimiterNone` | — | 429 ≥ 1, 최대 관측 QPS > 20 |
| bursty | `NewBurstyBuilder(40, 2s)` | 20 | **429 ≥ 1** (평균은 맞지만 여전히 발생) |
| smooth | `NewSmoothBuilder(10, 1s)` | 10 | 429 == 0, 최대 관측 QPS ≤ 20 |

두 번째 행이 이 테이블의 존재 이유다. **평균 20 QPS로 한도에 정확히 맞췄는데도 429가 난다.** 40개를 주기 시작에 한꺼번에 통과시키기 때문이다. "1분에 1200회까지 가능하다"는 계산이 왜 안전을 보장하지 않는지를 숫자로 고정한다.

세 번째 행은 한도의 50%(10/s)를 쓴다. 20/s로 두면 균등 간격이라도 sliding window 경계에서 21개로 셀 여지가 있어 마진이 없다. 벤더 한도의 50~70%만 쓰라는 권고안과도 일치한다.

**세 케이스 모두 `MaxRetries: 0`으로 둔다.** 재시도를 켜면 429를 받은 요청이 다시 성공하면서 서버가 센 429 수가 늘어나고, 최종 성공 건수로는 케이스 간 차이가 사라진다. 이 테스트는 "429가 났는가"를 비교하는 것이므로 안전망을 꺼야 한다. 재시도의 효과는 별도 테스트에서 확인한다.

### `TestRefresher_키_중복제거로_호출수가_줄어든다`

로봇 4대 × 지도 10개 = 40키(지도 중복)를 입력하면 서버 호출은 10회, 429는 0.

### `TestClient_429를_받으면_RetryAfter만큼_기다린다`

서버 한도를 5/s로 낮추고 rate limiter 없이 **10건**을 동시 발사해 일부러 429를 유발한다. 재시도 정책을 켜면 최종적으로 10건 전부 성공한다. 안전망이 동작함을 확인한다.

40건이 아니라 10건인 이유는 `Retry-After: 1` 대기가 누적되어 테스트가 길어지기 때문이다. 안전망 동작 확인에는 10건으로 충분하다.

### 테스트로 덮지 않는 것

- **bulkhead**: `Options.MaxConcurrency`로 구성은 제공하지만 전용 테스트는 없다. 로컬 `httptest` 서버는 응답이 즉시 돌아와 동시 실행 수가 쌓이지 않으므로, 의미 있는 검증을 하려면 인위적 지연을 넣어야 하고 그만큼 테스트가 느려진다. README에 "왜 필요한가"만 설명한다.
- **stale 서빙**: `Refresher`에 구현하되 실패 주입 테스트는 넣지 않는다. 갱신 실패 경로는 이 예제의 주제(버스트)와 별개 축이다.

### 실행 시간

smooth 케이스는 40개 / 10 QPS = **약 4초**가 걸린다. bursty 케이스는 약 2초, 나머지는 1초 미만이다. 파일 전체가 10초 안팎이 될 것으로 예상한다.

이 저장소에는 testcontainers 기반의 더 느린 테스트가 이미 있으므로 그대로 둔다. 다만 `TestClient_429를_받으면_RetryAfter만큼_기다린다`는 `Retry-After: 1` 대기가 누적되므로 `testing.Short()` 가드를 붙인다.

### 결정론성

QPS 측정은 시간에 의존하므로 flaky 가능성이 있다. 완화 방법:

- 서버를 sliding window로 두어 클라이언트 주기와의 정렬 운에 결과가 좌우되지 않게 한다
- 단정은 정확한 개수가 아니라 **방향**으로 한다 (`429 >= 1` / `429 == 0`, `maxQPS > 20` / `maxQPS <= 20`)
- smooth 케이스는 한도의 50%만 사용해 경계에서 흔들릴 여지를 없앤다
- `httptest` 서버는 로컬 루프백이라 40개 goroutine 요청이 수 ms 내에 도달한다

## 문서

- `golang/resilience/outbound/README.md` — 문제 정의(평균 0.67 QPS vs 순간 40 QPS), Best Practice 계층 표, 테스트 결과 수치, resty 조사 결과, 다음 단계
- `golang/resilience/README.md` — 구조 트리에 `outbound/` 항목 추가

## 의존성 변경

신규 추가 없음. `github.com/failsafe-go/failsafe-go v0.9.6`이 `// indirect`에서 direct로 승격된다 (`go mod tidy`가 처리).

`golang.org/x/sync`(singleflight, errgroup)는 이미 direct다. `golang.org/x/time`은 대조군에서 빠졌으므로 indirect로 남는다.

## 열린 질문

없음.
