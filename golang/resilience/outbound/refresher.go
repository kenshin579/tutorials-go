package outbound

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/sync/singleflight"
)

// Stats는 RefreshAll 한 번의 결과를 요약한다.
type Stats struct {
	Requested int           // 입력된 키 개수 (중복 포함)
	Unique    int           // 중복 제거 후 키 개수
	Failed    int           // 갱신에 실패한 키 개수
	Elapsed   time.Duration // 전체 소요 시간
}

// Refresher는 키 목록에 대한 캐시 갱신을 수행한다.
type Refresher struct {
	client  *http.Client
	baseURL string
	group   singleflight.Group

	mu    sync.RWMutex
	cache map[string]string
}

// NewRefresher는 client로 baseURL을 호출하는 Refresher를 만든다.
func NewRefresher(client *http.Client, baseURL string) *Refresher {
	return &Refresher{
		client:  client,
		baseURL: baseURL,
		cache:   make(map[string]string),
	}
}

// RefreshAll은 keys의 중복을 제거한 뒤 병렬로 캐시를 갱신한다.
//
// 갱신에 실패한 키는 기존 캐시 값을 그대로 유지한다(stale 서빙). 캐시 갱신
// 실패가 서비스 응답 실패로 번지지 않게 하기 위함이다. 같은 이유로 개별 키의
// 실패가 다른 키의 갱신을 취소하지 않는다.
func (r *Refresher) RefreshAll(ctx context.Context, keys []string) Stats {
	start := time.Now()
	unique := dedupe(keys)

	var failed atomic.Int64
	var wg sync.WaitGroup
	for _, key := range unique {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := r.refresh(ctx, key); err != nil {
				failed.Add(1)
			}
		}()
	}
	wg.Wait()

	return Stats{
		Requested: len(keys),
		Unique:    len(unique),
		Failed:    int(failed.Load()),
		Elapsed:   time.Since(start),
	}
}

// refresh는 키 하나를 갱신한다.
// 같은 키에 대해 이미 진행 중인 호출이 있으면 그 결과를 공유한다.
func (r *Refresher) refresh(ctx context.Context, key string) error {
	_, err, _ := r.group.Do(key, func() (any, error) {
		value, err := r.fetch(ctx, key)
		if err != nil {
			return nil, err
		}
		r.set(key, value)
		return value, nil
	})
	return err
}

// fetch는 벤더 API를 호출한다.
func (r *Refresher) fetch(ctx context.Context, key string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, r.baseURL+"/maps/"+key, nil)
	if err != nil {
		return "", err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status %d for key %s", resp.StatusCode, key)
	}

	return string(body), nil
}

// Get은 캐시된 값을 반환한다.
func (r *Refresher) Get(key string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	value, ok := r.cache[key]
	return value, ok
}

// set은 캐시 값을 저장한다.
func (r *Refresher) set(key, value string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cache[key] = value
}

// dedupe는 순서를 유지하며 중복 키를 제거한다.
//
// 로봇 4대 x 지도 10개처럼 캐시 키에 로봇 ID가 섞여 들어가면 같은 지도를 4번
// 호출하게 된다. 지도 단위로 정규화하면 호출 수가 1/4로 줄어든다. 호출 수를
// 줄이는 것이 스로틀링보다 근본적인 개선이다.
func dedupe(keys []string) []string {
	seen := make(map[string]struct{}, len(keys))
	unique := make([]string, 0, len(keys))
	for _, key := range keys {
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		unique = append(unique, key)
	}
	return unique
}
