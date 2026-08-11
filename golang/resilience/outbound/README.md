# 외부 API 버스트 호출 방지

캐시 갱신을 위해 외부 벤더 API를 주기적으로 호출할 때, 요청을 한꺼번에 발사해
rate limit(429)에 걸리는 문제와 그 해결책을 다룬다.

## 문제

- 벤더 제한: 20 QPS
- 호출량: 로봇 4대 x 지도 10개 = 40 호출/분
- 평균 QPS: **0.67** (한도의 3%)
- 순간 QPS: **최대 40** (동시 발사 시)

총량 문제가 아니라 **분산 문제**다. 40개를 몇 초에 걸쳐 흘려보내기만 해도 한도
안에 들어간다. 로봇 대수가 늘면 호출량은 비례해서 늘어나므로 이 구조는 그대로
두면 재현된다.

## 핵심: 평균 QPS를 맞추는 것으로는 부족하다

"20 QPS면 1분에 1200회까지 가능하다"는 계산은 평균만 본 것이다. 제약은 순간
QPS에 걸린다. 아래는 `TestClient_버스트_요청과_스로틀링_비교`가 40개 요청을
20 QPS 한도(sliding window) 서버에 보낸 실측 결과다.

| 방식 | 설정 | 총 호출 | 429 | 최대 관측 QPS | 소요 |
|---|---|---|---|---|---|
| 정책 없음 | — | 40 | 20 | 40 | 8.5ms |
| bursty | 40회 / 2초 (평균 20 QPS) | 40 | 20 | 40 | 11.8ms |
| smooth | 10회 / 1초 (평균 10 QPS) | 40 | 0 | 11 | 3.90s |

`bursty`는 **평균을 한도에 정확히 맞췄는데도 429를 받는다.** 주기 시작 시점에
40개를 지연 없이 통과시키기 때문이다. `smooth`는 `Period / MaxExecutions`
간격을 강제해 버스트 자체가 불가능하다.

### `bursty`와 `정책 없음`의 수치가 같은 것은 버그가 아니다

두 행의 총 호출·429·최대 관측 QPS가 완전히 같다. 이 예제는 그렇게 되도록
설계됐다: 40개짜리 단발 버스트 앞에서는 "주기 시작에 허용량을 한꺼번에 다
쓰는" bursty limiter와 "limiter가 없는" 상태가 구분되지 않는다. **설정되어
있고 동작 중인 limiter가 limiter가 없는 것과 똑같이 행동할 수 있다**는 것이
이 비교의 요점이다.

그래서 이 비교 테스트 하나만으로는 bursty limiter가 실제로 붙어 있는지, 아니면
우연히 코드에서 빠졌는지 구분할 수 없다. 그 증명은
`TestNewClient_bursty모드는_주기_안의_버스트를_허용하고_초과분은_ErrExceeded로_막는다`가
담당한다. 이 테스트는 주기 안의 첫 `MaxExecutions`개가 지연 없이 통과하고
(smooth와 구분되는 속성), 허용량을 넘긴 다음 요청은 `ratelimiter.ErrExceeded`로
막힌다는 것(정책 없음과 구분되는 속성)을 직접 검증한다.

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
호출을 1/4로 줄일 수 있다. `TestRefresher_키_중복제거로_호출수가_줄어든다`가
이 효과를 그대로 실측한다: 입력 키 40 → 중복 제거 후 10 → 서버 호출 10, 429 0,
소요 901.7ms.

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

**`MaxWaitTime`의 정직한 상한은 `2 × MaxWaitTime × (1 + MaxRetries)`다.** 정책이
`Limiter(Bulkhead(func))`로 중첩되므로 한 요청이 limiter에서 `MaxWaitTime`,
bulkhead에서 다시 `MaxWaitTime`을 기다릴 수 있다(시도 1회당 2배). 게다가 retry가
가장 바깥이라 실제 429에 대한 재시도는 시도마다 limiter를 다시 통과한다 — 의도된
동작이다([정책 순서가 중요한 이유](#정책-순서가-중요한-이유)). 위 예제
(`MaxWaitTime: 30s`, `MaxRetries: 3`)면 최악 240초다. `MaxRetries`를 올릴 때 이
곱셈을 함께 감안해야 한다.

이 상한은 permit/slot 대기만 센 값이다. 벤더가 `Retry-After`로 지정한 대기는
재시도 사이에 별도 타이머로 더해지므로 여기에 포함되지 않는다. `Retry-After`가
10초인 벤더라면 위 예제에서 `3 × 10초`가 240초 위에 얹힌다.

**단, 우리 쪽 limiter/bulkhead가 포기한 경우는 재시도하지 않는다.** 재시도 정책이
`ratelimiter.ErrExceeded`와 `bulkhead.ErrFull`을 `AbortOnErrors`로 제외한다.
"벤더가 429를 줬다"와 "우리 limiter가 `MaxWaitTime`을 다 기다려도 permit을 못
얻었다"는 다른 신호다. 전자는 `Retry-After`를 존중해 재시도할 가치가 있지만,
후자에 지연 0으로 즉시 재진입하면 성공할 리가 없다. 이걸 걸러주지 않으면 그것이
이 예제가 경고하는 재시도 폭풍 그 자체가 된다.

**잘못된 설정은 조용히 넘어가지 않고 panic한다.** `Mode`가 `LimiterNone`이 아닌데
`MaxExecutions <= 0`이거나 `Period <= 0`이면 `NewClient`가 panic한다. `int`를
`uint`로 바꾸기 전에 걸러내지 않으면 `0`은 division-by-zero panic으로, 음수는
`uint`로 wrap되어 거의 무한대의 실행을 허용하는 "스로틀링이 조용히 사라지는"
결과로 이어진다. 이 예제 전체가 경고하는 실패 모드이므로 설정 단계에서 막는다.

### 정책 순서가 중요한 이유

`failsafe`는 첫 인자가 가장 바깥 정책이므로 `retry → rateLimiter → bulkhead`
순으로 넘긴다. 이렇게 해야 재시도 요청도 rate limiter를 다시 통과한다. 순서를
뒤집어 limiter가 바깥에 오면 재시도가 limiter를 우회해 429 폭풍을 만든다.

### 클라이언트 하나를 공유해야 하는 이유

rate limit은 보통 계정 단위다. 호출 지점마다 limiter를 따로 두면 지점별로는
한도 안이어도 합산하면 다시 한도를 넘는다. `NewClient`가 반환하는
`*http.Client` 하나를 모든 호출 지점이 공유해야 한다.

### bulkhead가 별도 계층인 이유

rate limiter는 실행 *시작 시점*만 제어하고 in-flight 수는 제한하지 않는다.
응답이 느려지면 동시 연결이 쌓인다. 벤더가 동시 연결 수도 제한하는 경우가 많아
별도 층이 필요하다. 로컬 `httptest` 서버는 보통 응답이 즉시 돌아와 동시 실행
수가 쌓이지 않으므로, bulkhead 동작을 의미 있게 검증하려면 인위적 지연을 넣어야
한다. `TestClient_bulkhead가_고갈되면_재시도하지_않고_즉시_실패한다`가 응답을
일부러 지연시키는 서버로 이 지연 계층을 직접 검증한다(아래
[재시도가 limiter/bulkhead 고갈로 새지 않는지 검증하기](#재시도가-limiterbulkhead-고갈로-새지-않는지-검증하기)
참고).

## singleflight가 실제로 하는 일

`Refresher.RefreshAll`은 `dedupe`로 키 중복을 제거한 **뒤에** goroutine을
띄운다. 그래서 위 40 → 10 감소는 전부 `dedupe`가 만든 것이고, 같은 `RefreshAll`
호출 안에서는 같은 키를 가진 두 goroutine이 동시에 `group.Do`에 들어올 일이
없다 — singleflight가 합칠 대상이 아예 없다.

singleflight(`Refresher.group`)가 실제로 막는 것은 **RefreshAll 호출들이
겹치는 경우**다: 이전 주기의 갱신이 아직 끝나지 않았는데 다음 주기가 같은 키를
다시 요청하면(예: 갱신이 느려져 티커 주기를 넘기는 경우), singleflight가 그
요청을 진행 중인 호출에 합류시켜 벤더 API를 다시 부르지 않는다. 이 보장은
`TestRefresher_겹치는_요청은_singleflight로_합류하여_벤더_호출이_한번만_일어난다`가
느린 벤더 응답을 채널로 흉내내어 직접 증명한다.

## goroutine fan-out과 스케일링 가정

`RefreshAll`은 중복 제거된 키마다 goroutine을 하나씩 띄우고 상한을 두지 않는다.
지도 수십 개 수준에서는 문제가 없지만, 키가 훨씬 많아지면(수백~수천) goroutine
생성 자체가 부담이 될 수 있어 worker pool로 상한을 두는 구조가 필요해진다. rate
limiter와 bulkhead는 **HTTP 호출**의 시작 시점과 동시 실행 수를 제한할 뿐,
goroutine 생성 자체를 제한하지 않는다.

## resty 내장 rate limiter를 쓰지 않는 이유

- **resty v2** (`v2.9.0`+): `SetRateLimiter`의 인터페이스가 `Allow() bool`
  하나뿐이다. 논블로킹으로 초과분을 **즉시 버리고**, `wrapNoRetryErr`로 감싸서
  재시도조차 하지 않는다. 캐시 갱신처럼 "늦게라도 다 보내야 하는" 요청에는
  부적합하다. 참고로 이 저장소는 `v2.7.0`이라 이 API 자체가 없다.
- **resty v3**: `SetRateLimiter(resty.NewRateLimitTokenBucket(rate, burst))`가
  토큰이 생길 때까지 **대기**한다. 시맨틱은 적합하나 아직 `v3.0.0-rc.3`으로
  정식 릴리스 전이다.

`failsafehttp`는 `RoundTripper` 층에서 동작하므로 resty를 쓰든 안 쓰든 그대로
적용되고, 클라이언트 라이브러리에 종속되지 않는다.

## 테스트

```bash
# 전체 (약 5~10초)
go test ./golang/resilience/outbound/... -v

# Retry-After 대기(약 1초) 테스트 제외
go test ./golang/resilience/outbound/... -short -v
```

가짜 벤더 서버(`fakeserver_test.go`)는 최근 1초 **sliding window**로 요청을
센다. 고정 윈도우는 경계에서 최대 2배를 허용해 클라이언트 주기와 우연히
정렬되면 버스트 설정으로도 429가 나지 않는다. 그러면 테스트 통과가 처방 덕인지
정렬 운인지 구분할 수 없다. sliding window는 정렬에 무관하게 순간 QPS를
판정하므로 결과가 결정론적이다.

`Retry-After` 재시도 안전망은 `TestClient_429를_받으면_RetryAfter만큼_기다린다`가
검증한다: 서버 한도를 5/s로 낮추고 10건을 동시 발사해 429를 유발하면, 429 5건이
발생하고 재시도로 서버 호출이 10건에서 15건으로 늘어나며, 소요 1.00s 뒤 최종적으로
모든 키가 성공한다.

### 재시도가 limiter/bulkhead 고갈로 새지 않는지 검증하기

[구성](#구성)의 `AbortOnErrors`는 두 테스트가 다른 각도로 검증한다.

- `TestClient_bulkhead가_고갈되면_재시도하지_않고_즉시_실패한다` — **판별력 있음.**
  bulkhead는 슬롯이 빌 때까지 실제 타이머로 기다리므로, `AbortOnErrors`를 걷어내면
  재시도마다 그 대기가 반복되어 시간이 눈에 보이게 늘어난다(`MaxWaitTime: 20ms`,
  `MaxRetries: 3`에서 약 20ms → 약 84ms). 이 시간 차이를 직접 잰다.
- `TestClient_limiter가_고갈되면_ErrExceeded가_그대로_표면화된다` — **판별력 없음.**
  rate limiter는 대기가 `MaxWaitTime`을 넘길지 미리 계산해 즉시 실패하므로, 재시도가
  새더라도 시간이 늘지 않는다(`MaxRetries`를 200으로 올려도 수십 마이크로초). 이
  테스트는 에러 종류와 서버 미도달만 확인한다.

두 정책이 같은 `AbortOnErrors` 한 줄로 함께 걸러지므로, 판별력 있는 bulkhead
테스트가 통과하면 rate limiter 쪽도 같은 메커니즘으로 보호된다. 재시도 폭풍이
bulkhead에서는 지연으로, rate limiter에서는 낭비되는 재시도 횟수로 나타나는 차이일
뿐이다.

## 다음 단계

- **분산 환경**: 프로세스별 limiter는 인스턴스 수만큼 합산되어 한도를 넘는다.
  `../ratelimit/redis_limiter.go`의 Redis 기반 분산 limiter를 참고한다.
- **adaptive throttling**: 429 비율을 보고 rate를 자동으로 낮추는 방식.
  `failsafe-go`의 `adaptivethrottler` 패키지.
- **호출 자체 없애기**: 벤더에 배치 API, webhook/push, ETag 조건부 요청 지원
  여부를 문의한다. 조건부 요청의 304가 rate limit에 계산되는지도 확인이
  필요하다.
