# 시간 의존 코드 테스트 예제 구현 계획

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** `time.Time` 의존 코드를 테스트하는 4가지 패턴 예제(`golang/testing/clock/`)와 blog-v2 블로그 초안을 작성한다.

**Architecture:** 하나의 package `clock` 안에 패턴별 구현+테스트 쌍 4개(쿠폰/주문/TTL캐시/synctest)를 두고, 폴더 README와 블로그 초안(blog-v2 `docs/start/`)이 같은 내용을 공유한다. 스펙: `docs/superpowers/specs/2026-08-11-time-dependent-testing-design.md`

**Tech Stack:** Go 1.26, testify/assert, jonboulle/clockwork v0.4.0, testing/synctest (Go 1.25+)

**작업 위치:** `tutorials-go` 저장소, 브랜치 `feature/time-dependent-testing-example` (이미 생성됨). Task 6만 `blog-v2.advenoh.pe.kr` 저장소에서 수행.

**공통 규칙:**
- 테스트 함수명 `Test_대상_설명` (기존 `go1_25/synctest_test.go` 스타일)
- testify/assert 사용, 테이블 드리븐 선호
- 커밋 메시지: conventional commits prefix + 한국어 설명
- 모든 명령은 `tutorials-go/` 루트에서 실행

---

### Task 1: 패턴 ① — 시간을 파라미터로 받기 (Coupon)

**Files:**
- Create: `golang/testing/clock/coupon.go`
- Test: `golang/testing/clock/coupon_test.go`

- [ ] **Step 1: 실패하는 테스트 작성**

`golang/testing/clock/coupon_test.go`:

```go
package clock

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Coupon_IsExpiredAt_경계값(t *testing.T) {
	expiresAt := time.Date(2026, 8, 12, 0, 0, 0, 0, time.UTC)
	coupon := Coupon{Code: "WELCOME10", ExpiresAt: expiresAt}

	tests := []struct {
		name    string
		now     time.Time
		expired bool
	}{
		{"만료 1초 전", expiresAt.Add(-time.Second), false},
		{"만료 시각과 동일", expiresAt, false},
		{"만료 1초 후", expiresAt.Add(time.Second), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expired, coupon.IsExpiredAt(tt.now))
		})
	}
}
```

- [ ] **Step 2: 테스트 실패 확인**

Run: `go test ./golang/testing/clock/... -run Test_Coupon -v`
Expected: FAIL (컴파일 에러 — `undefined: Coupon`)

- [ ] **Step 3: 최소 구현 작성**

`golang/testing/clock/coupon.go`:

```go
package clock

import "time"

// Coupon은 만료 시각을 가진 쿠폰이다.
type Coupon struct {
	Code      string
	ExpiresAt time.Time
}

// IsExpiredAt은 주어진 시각 기준으로 쿠폰 만료 여부를 반환한다.
// 시간을 파라미터로 받는 순수 함수라서 실행 시점과 무관하게 테스트할 수 있다.
func (c Coupon) IsExpiredAt(now time.Time) bool {
	return now.After(c.ExpiresAt)
}
```

- [ ] **Step 4: 테스트 통과 확인**

Run: `go test ./golang/testing/clock/... -run Test_Coupon -v`
Expected: PASS (서브테스트 3개 모두)

- [ ] **Step 5: 커밋**

```bash
git add golang/testing/clock/coupon.go golang/testing/clock/coupon_test.go
git commit -m "feat: 시간을 파라미터로 받는 쿠폰 만료 검사 예제 추가"
```

---

### Task 2: 패턴 ② — nowFunc 필드 주입 (OrderService)

**Files:**
- Create: `golang/testing/clock/order.go`
- Test: `golang/testing/clock/order_test.go`

- [ ] **Step 1: 실패하는 테스트 작성**

`golang/testing/clock/order_test.go`:

```go
package clock

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_OrderService_고정_시간_주입(t *testing.T) {
	fixed := time.Date(2026, 8, 12, 9, 30, 0, 0, time.UTC)
	svc := &OrderService{nowFunc: func() time.Time { return fixed }}

	order := svc.CreateOrder("order-1")

	assert.Equal(t, fixed, order.CreatedAt)
}

func Test_OrderService_주입_없이_WithinDuration_검증(t *testing.T) {
	svc := NewOrderService()

	order := svc.CreateOrder("order-2")

	// 시간을 주입할 수 없을 때의 차선책: 정확한 일치 대신 오차 허용 비교
	assert.WithinDuration(t, time.Now(), order.CreatedAt, time.Second)
}
```

- [ ] **Step 2: 테스트 실패 확인**

Run: `go test ./golang/testing/clock/... -run Test_OrderService -v`
Expected: FAIL (컴파일 에러 — `undefined: OrderService`)

- [ ] **Step 3: 최소 구현 작성**

`golang/testing/clock/order.go`:

```go
package clock

import "time"

// Order는 생성 시각이 기록되는 주문이다.
type Order struct {
	ID        string
	CreatedAt time.Time
}

// OrderService는 주문을 생성하는 서비스다.
// nowFunc 필드로 현재 시각 함수를 주입받아 테스트에서 시간을 고정할 수 있다.
// 인터페이스 없이 함수 필드 하나로 해결하는 가장 가벼운 주입 패턴이다.
type OrderService struct {
	nowFunc func() time.Time
}

// NewOrderService는 실제 시간(time.Now)을 사용하는 OrderService를 생성한다.
func NewOrderService() *OrderService {
	return &OrderService{nowFunc: time.Now}
}

// CreateOrder는 현재 시각을 생성 시각으로 기록한 주문을 만든다.
func (s *OrderService) CreateOrder(id string) Order {
	return Order{ID: id, CreatedAt: s.nowFunc()}
}
```

- [ ] **Step 4: 테스트 통과 확인**

Run: `go test ./golang/testing/clock/... -run Test_OrderService -v`
Expected: PASS (테스트 2개 모두)

- [ ] **Step 5: 커밋**

```bash
git add golang/testing/clock/order.go golang/testing/clock/order_test.go
git commit -m "feat: nowFunc 필드 주입으로 시간을 고정하는 주문 서비스 예제 추가"
```

---

### Task 3: 패턴 ③ — Clock 인터페이스 주입 (TTLCache + clockwork)

**Files:**
- Create: `golang/testing/clock/cache.go`
- Test: `golang/testing/clock/cache_test.go`
- Modify: `go.mod` (clockwork indirect → direct, `go mod tidy`로 자동 처리)

- [ ] **Step 1: 실패하는 테스트 작성**

`golang/testing/clock/cache_test.go`:

```go
package clock

import (
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
)

func Test_TTLCache_가짜_시계로_만료_검증(t *testing.T) {
	fakeClock := clockwork.NewFakeClock()
	cache := NewTTLCache(fakeClock, 10*time.Minute)

	cache.Set("session", "user-42")

	// TTL 직전: 아직 살아있다
	fakeClock.Advance(10*time.Minute - time.Second)
	v, ok := cache.Get("session")
	assert.True(t, ok)
	assert.Equal(t, "user-42", v)

	// TTL 경과: 만료된다 (실제로 10분을 기다리지 않는다)
	fakeClock.Advance(time.Second)
	_, ok = cache.Get("session")
	assert.False(t, ok)
}

func Test_TTLCache_없는_키(t *testing.T) {
	cache := NewTTLCache(clockwork.NewFakeClock(), time.Minute)

	_, ok := cache.Get("missing")

	assert.False(t, ok)
}
```

- [ ] **Step 2: 테스트 실패 확인**

Run: `go test ./golang/testing/clock/... -run Test_TTLCache -v`
Expected: FAIL (컴파일 에러 — `undefined: NewTTLCache`)

- [ ] **Step 3: 최소 구현 작성**

`golang/testing/clock/cache.go`:

```go
package clock

import (
	"sync"
	"time"

	"github.com/jonboulle/clockwork"
)

type cacheItem struct {
	value     string
	expiresAt time.Time
}

// TTLCache는 TTL이 지나면 항목이 만료되는 인메모리 캐시다.
// clockwork.Clock을 주입받아 테스트에서 가짜 시계로 시간을 "진행"시킬 수 있다.
// nowFunc 주입(패턴 ②)과 달리 시간이 흘러야 검증되는 로직에 적합하다.
type TTLCache struct {
	clock clockwork.Clock
	ttl   time.Duration

	mu    sync.Mutex
	items map[string]cacheItem
}

// NewTTLCache는 주어진 시계와 TTL로 캐시를 생성한다.
// 프로덕션에서는 clockwork.NewRealClock()을 전달한다.
func NewTTLCache(clock clockwork.Clock, ttl time.Duration) *TTLCache {
	return &TTLCache{
		clock: clock,
		ttl:   ttl,
		items: make(map[string]cacheItem),
	}
}

// Set은 key에 value를 저장하고 TTL 이후 만료되도록 기록한다.
func (c *TTLCache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = cacheItem{
		value:     value,
		expiresAt: c.clock.Now().Add(c.ttl),
	}
}

// Get은 key의 값을 반환한다. 항목이 없거나 만료됐으면 false를 반환한다.
func (c *TTLCache) Get(key string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	it, ok := c.items[key]
	if !ok {
		return "", false
	}
	// now >= expiresAt이면 만료
	if !c.clock.Now().Before(it.expiresAt) {
		delete(c.items, key)
		return "", false
	}
	return it.value, true
}
```

- [ ] **Step 4: 의존성 정리 및 테스트 통과 확인**

```bash
go mod tidy
go test ./golang/testing/clock/... -run Test_TTLCache -v
```

Expected: PASS (테스트 2개 모두). `go.mod`에서 `github.com/jonboulle/clockwork v0.4.0`의 `// indirect` 주석이 제거됨 (`git diff go.mod`로 확인).

- [ ] **Step 5: 커밋**

```bash
git add golang/testing/clock/cache.go golang/testing/clock/cache_test.go go.mod go.sum
git commit -m "feat: clockwork Clock 인터페이스 주입으로 시간 진행을 시뮬레이션하는 TTL 캐시 예제 추가"
```

---

### Task 4: 패턴 ④ — testing/synctest로 같은 TTLCache 테스트

**Files:**
- Test: `golang/testing/clock/cache_synctest_test.go` (구현 파일 없음 — Task 3의 `TTLCache` 재사용)

- [ ] **Step 1: synctest 테스트 작성**

`golang/testing/clock/cache_synctest_test.go`:

```go
package clock

import (
	"testing"
	"testing/synctest"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
)

// 패턴 ③(가짜 시계 주입)과 달리, synctest는 코드 수정 없이
// 실제 시계(RealClock)를 쓰는 캐시를 가상 시간 버블 안에서 테스트한다.
// 버블 안에서는 time.Now()/time.Sleep()이 가상 시간으로 동작한다 (Go 1.25+).
func Test_TTLCache_Synctest_실제_시계로_만료_검증(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		cache := NewTTLCache(clockwork.NewRealClock(), 10*time.Minute)

		cache.Set("session", "user-42")

		// 버블 안의 time.Sleep은 실제로 기다리지 않고 가상 시간을 진행시킨다
		time.Sleep(10*time.Minute - time.Second)
		v, ok := cache.Get("session")
		assert.True(t, ok)
		assert.Equal(t, "user-42", v)

		time.Sleep(time.Second)
		_, ok = cache.Get("session")
		assert.False(t, ok)
	})
}
```

- [ ] **Step 2: 테스트 통과 확인 (즉시 완료되는지 확인)**

Run: `go test ./golang/testing/clock/... -run Test_TTLCache_Synctest -v`
Expected: PASS, 실행 시간 1초 미만 (가상 시간이므로 10분을 기다리지 않음). 만약 실제로 오래 걸리면 synctest 버블이 동작하지 않는 것이므로 실패로 간주.

- [ ] **Step 3: 커밋**

```bash
git add golang/testing/clock/cache_synctest_test.go
git commit -m "feat: testing/synctest로 코드 수정 없이 TTL 캐시를 테스트하는 예제 추가"
```

---

### Task 5: 폴더 README + 전체 검증

**Files:**
- Create: `golang/testing/clock/README.md`

- [ ] **Step 1: README 작성**

`golang/testing/clock/README.md`:

````markdown
# 시간 의존 코드 테스트하기

`time.Now()`를 직접 호출하는 코드는 실행 시점마다 결과가 달라 테스트가 불안정하다(flaky).
이 예제는 시간 의존 코드를 테스트하는 4가지 패턴을 점진적으로 보여준다.

## 패턴 요약

| # | 패턴 | 예제 파일 | 핵심 아이디어 |
|---|------|----------|--------------|
| ① | 시간을 파라미터로 받기 | `coupon.go` | `IsExpiredAt(now time.Time)` — 순수 함수로 만들어 테스트에서 임의 시각 전달 |
| ② | nowFunc 필드 주입 | `order.go` | 구조체에 `nowFunc func() time.Time` 필드(기본값 `time.Now`), 테스트에서 고정 시간 주입 |
| ③ | Clock 인터페이스 주입 | `cache.go` | `clockwork.Clock` 주입, 테스트에서 `fakeClock.Advance()`로 시간 진행 |
| ④ | testing/synctest | `cache_synctest_test.go` | 코드 수정 없이 가상 시간 버블에서 테스트 (Go 1.25+) |

## 어떤 패턴을 언제 쓸까

| 상황 | 추천 패턴 |
|------|----------|
| 특정 시각 기준 판단만 필요 (만료 검사 등) | ① 파라미터 전달 |
| 구조체가 `time.Now()`를 저장/기록 | ② nowFunc 주입 — 실무의 80%는 이걸로 충분 |
| 시간이 "흘러야" 검증되는 로직 (TTL, 재시도, 스케줄러) | ③ Clock 인터페이스 (clockwork) |
| 동시성 + 시간 조합, 기존 코드 수정 불가 | ④ synctest (Go 1.25+) |
| 주입이 불가능한 레거시 코드 | `assert.WithinDuration`으로 오차 허용 비교 (`order_test.go` 참고) |

## 트레이드오프

- **① 파라미터 전달**: 가장 단순하고 순수하다. 다만 깊은 콜스택이면 모든 호출자가 시간을 넘겨야 한다.
- **② nowFunc 주입**: 인터페이스 없이 가볍고 관용적이다. `Now()`만 대체 가능하다.
- **③ Clock 인터페이스**: `Sleep`/`Ticker`/`After`까지 제어하고 시간 진행을 시뮬레이션한다. 의존성과 인터페이스 추가 비용이 든다.
- **④ synctest**: 표준 라이브러리이고 코드 수정이 필요 없다. 동시성+시간 조합에 특화되어 있고 Go 1.25 이상이 필요하다.

## 실행

```bash
go test ./golang/testing/clock/... -v
```

## 관련 예제

- [`golang/go1_25/synctest_test.go`](../../go1_25/synctest_test.go) — synctest 기초 (goroutine + 가상 시간)
- [`golang/time/`](../../time/) — time 패키지 API 사용법 (parse, zone, duration)
````

- [ ] **Step 2: 전체 검증**

```bash
go test ./golang/testing/clock/... -v
go vet ./golang/testing/clock/...
gofmt -l golang/testing/clock/
file -I golang/testing/clock/README.md
```

Expected: 테스트 전체 PASS(총 6개 함수), vet 출력 없음, gofmt 출력 없음(포맷 위반 없음), README는 `charset=utf-8`.

- [ ] **Step 3: 커밋**

```bash
git add golang/testing/clock/README.md
git commit -m "docs: 시간 의존 테스트 패턴 선택 가이드 README 추가"
```

---

### Task 6: 블로그 초안 작성 (blog-v2 저장소)

**Files:**
- Create: `blog-v2.advenoh.pe.kr/docs/start/go-시간-의존-코드-테스트하기-testing-time-dependent-code-in-go/index.md`

**주의:** 이 Task만 `blog-v2.advenoh.pe.kr` 저장소에서 수행한다. main 직접 커밋 금지.

- [ ] **Step 1: blog-v2 feature 브랜치 생성**

```bash
cd ../blog-v2.advenoh.pe.kr
git checkout main && git pull origin main
git checkout -b feature/time-testing-blog-draft
```

- [ ] **Step 2: 블로그 초안 작성**

`docs/start/go-시간-의존-코드-테스트하기-testing-time-dependent-code-in-go/index.md` 생성.

frontmatter (정확히 이 내용):

```yaml
---
title: "Go에서 시간 의존 코드 테스트하기 (Testing Time-Dependent Code in Go)"
description: "time.Now()를 직접 호출하는 코드는 테스트하기 어렵다. 시간을 파라미터로 받기, nowFunc 필드 주입, clockwork Clock 인터페이스, Go 1.25의 testing/synctest까지 시간 의존 코드를 테스트하는 4가지 패턴을 예제 코드와 함께 알아봅니다."
date: 2026-08-12
update: 2026-08-12
tags:
  - golang
  - go
  - testing
  - time
  - clockwork
  - synctest
  - unit-test
  - 테스트
  - 시간
---
```

본문은 스펙의 확정 목차(`docs/superpowers/specs/2026-08-11-time-dependent-testing-design.md`의 "목차 (확정)" 절)를 그대로 따른다. 작성 규칙:

- 문체: 한국어 문어체("~이다/~한다"), 기존 글(`go-1-25-변경사항-whats-new-in-go-1-25/index.md`) 스타일
- 헤딩: `# 1.제목` / `## 1.1 소제목` 번호 체계
- 코드 스니펫: **Task 1~4에서 실제 구현된 파일 내용을 그대로 복사** (요약/변형 금지 — 코드와 글의 불일치 방지)
- 서론 뒤 참고 자료 블록:

```markdown
> 참고 자료
> - [예제 전체 코드 (tutorials-go)](https://github.com/kenshin579/tutorials-go/tree/master/golang/testing/clock)
> - [jonboulle/clockwork](https://github.com/jonboulle/clockwork)
> - [testing/synctest 패키지 문서](https://pkg.go.dev/testing/synctest)
> - [Testing Time (Go Blog: synctest)](https://go.dev/blog/synctest)
```

- 6장 선택 가이드 표: Task 5 README의 "어떤 패턴을 언제 쓸까" 표와 동일 내용 사용
- 각 장의 마무리에 다음 패턴이 필요한 이유를 한 문장으로 연결 (점진적 서사)

- [ ] **Step 3: 인코딩 및 형식 검증**

```bash
file -I "docs/start/go-시간-의존-코드-테스트하기-testing-time-dependent-code-in-go/index.md"
```

Expected: `charset=utf-8`. 추가로 frontmatter가 `---`로 정확히 열리고 닫히는지, 코드 블록 언어 태그(```go)가 있는지 눈으로 확인.

- [ ] **Step 4: 커밋**

```bash
git add "docs/start/go-시간-의존-코드-테스트하기-testing-time-dependent-code-in-go/index.md"
git commit -m "docs: Go 시간 의존 코드 테스트하기 블로그 초안 추가"
```

---

### Task 7: 최종 확인

- [ ] **Step 1: tutorials-go 전체 상태 확인**

```bash
cd ../tutorials-go
go test ./golang/testing/clock/... && go vet ./golang/testing/clock/...
git log --oneline master..HEAD
git status
```

Expected: 테스트/vet 클린, 커밋 5개(coupon/order/cache/synctest/README) + 스펙 커밋들, working tree 클린.

- [ ] **Step 2: 사용자에게 결과 보고**

PR 생성 여부는 사용자에게 확인 후 진행 (superpowers:finishing-a-development-branch 스킬 사용).
