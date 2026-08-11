package outbound

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRefresher_키_중복제거로_호출수가_줄어든다(t *testing.T) {
	const (
		robots      = 4
		maps        = 10
		serverLimit = 20
	)

	server := newRateLimitedServer(serverLimit)
	defer server.Close()

	client := NewClient(Options{
		Mode:          LimiterSmooth,
		MaxExecutions: 10,
		Period:        time.Second,
		MaxWaitTime:   10 * time.Second,
	})
	refresher := NewRefresher(client, server.URL)

	// 로봇마다 같은 지도 목록을 요청한다.
	// 캐시 키에 로봇 ID가 섞이면 4 x 10 = 40회를 호출하지만,
	// 지도 단위로 정규화하면 10회면 된다.
	var keys []string
	for robot := 0; robot < robots; robot++ {
		for m := 0; m < maps; m++ {
			keys = append(keys, fmt.Sprintf("map-%02d", m))
		}
	}

	stats := refresher.RefreshAll(context.Background(), keys)
	got := server.stats()

	t.Logf("입력 키 %d, 중복 제거 후 %d, 서버 호출 %d, 429 %d, 소요 %v",
		stats.Requested, stats.Unique, got.Total, got.Rejected, stats.Elapsed)

	assert.Equal(t, robots*maps, stats.Requested)
	assert.Equal(t, maps, stats.Unique, "지도 단위로 중복이 제거되어야 한다")
	assert.Equal(t, maps, got.Total, "서버 호출은 지도 개수만큼만 일어나야 한다")
	assert.Zero(t, got.Rejected)
	assert.Zero(t, stats.Failed)
}

func TestRefresher_갱신에_실패해도_이전_캐시값을_유지한다(t *testing.T) {
	// 한도 0이면 모든 요청이 429를 받는다.
	server := newRateLimitedServer(0)
	defer server.Close()

	client := NewClient(Options{Mode: LimiterNone})
	refresher := NewRefresher(client, server.URL)

	refresher.set("map-00", "이전값")

	stats := refresher.RefreshAll(context.Background(), []string{"map-00"})
	assert.Equal(t, 1, stats.Failed, "갱신은 실패해야 한다")

	value, ok := refresher.Get("map-00")
	assert.True(t, ok, "실패해도 캐시 항목이 남아야 한다")
	assert.Equal(t, "이전값", value, "이전 값이 덮어써지면 안 된다")
}

// TestRefresher_겹치는_요청은_singleflight로_합류하여_벤더_호출이_한번만_일어난다는
// dedupe가 손댈 수 없는 영역을 증명한다: dedupe는 RefreshAll *한 번의 호출*
// 안에서만 중복을 제거하므로, 이전 주기의 갱신이 아직 끝나지 않은 상태에서
// 다음 주기가 같은 키를 다시 요청하는 "겹치는 호출" 시나리오에는 아무 효과가
// 없다. 이 시나리오를 막는 것은 Refresher.group(singleflight)이다.
//
// golang/singleflight/loader_test.go의 기법을 따른다: 먼저 진행 중인 호출이
// 벤더 서버 핸들러 안에서 막혀 있음을 채널로 확인한 뒤(entirely
// deterministic), 두 번째 호출이 group.Do에 합류할 시간을 스케줄러에게 주고
// (bounded margin), 그 다음에야 첫 번째 호출을 풀어준다.
func TestRefresher_겹치는_요청은_singleflight로_합류하여_벤더_호출이_한번만_일어난다(t *testing.T) {
	var calls atomic.Int64
	// entered는 핸들러가 실제로 실행됐음을 알린다. 버그로 두 번째 호출이
	// 합류하지 못하고 벤더 서버를 다시 부르면, 그 신호도 여기로 들어와야
	// 하므로 용량을 2로 둔다(그래야 버그가 나도 핸들러 goroutine이 영원히
	// 막히지 않는다).
	entered := make(chan struct{}, 2)
	// release를 닫기 전까지 핸들러는 응답하지 않는다 — 느린 벤더 API를 흉내낸다.
	release := make(chan struct{})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls.Add(1)
		entered <- struct{}{}
		<-release
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()

	client := NewClient(Options{Mode: LimiterNone})
	refresher := NewRefresher(client, server.URL)

	statsCh := make(chan Stats, 2)
	go func() {
		statsCh <- refresher.RefreshAll(context.Background(), []string{"map-00"})
	}()

	// 결정적인 부분: 첫 번째 호출이 벤더 서버 핸들러 안에서 막혀 있음을
	// 확인한다. group.Do는 fn(핸들러 호출)을 부르기 전에 map-00 항목을
	// 등록하므로, 이 신호를 받은 시점에는 이미 그 항목이 singleflight의
	// map에 들어가 있다.
	<-entered

	go func() {
		statsCh <- refresher.RefreshAll(context.Background(), []string{"map-00"})
	}()

	// 결정적이지 않은 부분: 두 번째 호출이 group.Do에 진입해 이미 등록된
	// map-00 항목을 발견하고 합류할 시간을 스케줄러에게 준다. 이 창(window)
	// 동안 첫 번째 호출은 release가 닫히기 전까지 절대 끝나지 않으므로, 이
	// 시간 안에 벤더 서버가 다시 호출된다면 그것은 합류 실패를 뜻한다.
	// (golang/singleflight/loader_test.go의 waitForDuplicates와 동일한 기법 —
	// "두 번째 호출이 이미 group.Do 안에 있다"는 사실 자체는 외부에서 관찰할
	// 수 없으므로, 이 마진에 기대는 것은 이 저장소의 기존 관례다.)
	select {
	case <-entered:
		t.Fatal("두 번째 호출이 합류하지 않고 벤더 서버를 다시 호출했다")
	case <-time.After(20 * time.Millisecond):
	}
	close(release)

	first := <-statsCh
	second := <-statsCh

	assert.Zero(t, first.Failed)
	assert.Zero(t, second.Failed)
	assert.Equal(t, int64(1), calls.Load(), "진행 중인 갱신에 합류해야 하며 벤더 서버를 다시 호출하면 안 된다")

	value, ok := refresher.Get("map-00")
	assert.True(t, ok)
	assert.Equal(t, "ok", value)
}

func TestDedupe_순서를_유지하며_중복을_제거한다(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want []string
	}{
		{
			name: "중복 없음",
			in:   []string{"a", "b", "c"},
			want: []string{"a", "b", "c"},
		},
		{
			name: "연속 중복",
			in:   []string{"a", "a", "b"},
			want: []string{"a", "b"},
		},
		{
			name: "떨어진 중복은 첫 등장 순서를 유지한다",
			in:   []string{"b", "a", "b", "c", "a"},
			want: []string{"b", "a", "c"},
		},
		{
			name: "빈 입력",
			in:   nil,
			want: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, dedupe(tt.in))
		})
	}
}
