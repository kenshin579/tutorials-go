package outbound

import (
	"context"
	"fmt"
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
