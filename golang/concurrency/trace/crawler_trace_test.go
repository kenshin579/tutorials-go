package _trace

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	rttrace "runtime/trace"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestCrawler_Task_Region_계측 - Worker Pool 기반 크롤러에 Task/Region/Log 계측 적용
// httptest 서버를 사용하여 외부 의존성 없이 자체 포함(self-contained) 테스트를 구성한다.
//
// trace 분석 포인트:
//  1. Task별 latency: 각 URL 크롤링에 걸린 시간
//  2. Region별 시간: fetch vs parse 단계 비교
//  3. Worker 활용률: goroutine이 idle인 시간 비율
//  4. Blocking 분석: channel 대기, HTTP 응답 대기 등
func TestCrawler_Task_Region_계측(t *testing.T) {
	// trace 수집 설정
	f, err := os.CreateTemp(t.TempDir(), "trace_crawler_*.out")
	assert.NoError(t, err)
	defer f.Close()

	err = rttrace.Start(f)
	assert.NoError(t, err)
	defer rttrace.Stop()

	// 테스트용 웹 서버 (가변 응답 시간)
	pages := map[string]string{
		"/":        `<a href="/about">About</a><a href="/blog">Blog</a>`,
		"/about":   `<h1>About Us</h1><a href="/contact">Contact</a>`,
		"/blog":    `<h1>Blog</h1><a href="/blog/post1">Post 1</a>`,
		"/contact": `<h1>Contact</h1>`,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 페이지별 다른 응답 시간 시뮬레이션
		time.Sleep(time.Duration(len(r.URL.Path)) * time.Millisecond)
		if body, ok := pages[r.URL.Path]; ok {
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprint(w, body)
		} else {
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	// 크롤러 실행
	ctx := context.Background()
	ctx, mainTask := rttrace.NewTask(ctx, "web-crawler")
	defer mainTask.End()

	results := crawl(ctx, server.URL, server.Client(), 3)

	assert.NotEmpty(t, results, "크롤링 결과가 있어야 한다")
	t.Logf("크롤링 완료: %d 페이지", len(results))
	for _, r := range results {
		t.Logf("  %s: %d bytes (%.1fms)", r.url, r.bodyLen, r.elapsed.Seconds()*1000)
	}

	info, _ := f.Stat()
	assert.Positive(t, info.Size())
	t.Logf("trace 파일: %s (%d bytes)", f.Name(), info.Size())
}

// crawlResult는 크롤링 결과를 담는 구조체
type crawlResult struct {
	url     string
	bodyLen int
	elapsed time.Duration
}

// crawl은 Worker Pool 패턴으로 URL을 크롤링하고 Task/Region으로 계측한다
func crawl(ctx context.Context, baseURL string, client *http.Client, numWorkers int) []crawlResult {
	urls := make(chan string, 20)
	results := make(chan crawlResult, 20)

	// 방문 추적 (중복 방지)
	visited := &sync.Map{}

	// Worker Pool 시작
	var wg sync.WaitGroup
	for w := range numWorkers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			workerName := fmt.Sprintf("crawler-worker-%d", w)

			for url := range urls {
				// 각 URL 크롤링을 Task로 계측
				taskCtx, task := rttrace.NewTask(ctx, "crawl-page")
				rttrace.Log(taskCtx, "url", url)
				rttrace.Log(taskCtx, "worker", workerName)

				result := fetchPage(taskCtx, client, url)
				if result != nil {
					results <- *result
				}

				task.End()
			}
		}()
	}

	// 시드 URL 투입
	seedPaths := []string{"/", "/about", "/blog", "/contact"}
	go func() {
		for _, path := range seedPaths {
			url := baseURL + path
			if _, loaded := visited.LoadOrStore(url, true); !loaded {
				urls <- url
			}
		}
		close(urls)
	}()

	// Worker 완료 후 results 닫기
	go func() {
		wg.Wait()
		close(results)
	}()

	// 결과 수집
	var collected []crawlResult
	for r := range results {
		collected = append(collected, r)
	}
	return collected
}

// fetchPage는 단일 페이지를 가져오고 Region으로 각 단계를 계측한다
func fetchPage(ctx context.Context, client *http.Client, url string) *crawlResult {
	start := time.Now()

	var body []byte

	// Region: HTTP 요청
	rttrace.WithRegion(ctx, "http-fetch", func() {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			rttrace.Log(ctx, "error", err.Error())
			return
		}

		resp, err := client.Do(req)
		if err != nil {
			rttrace.Log(ctx, "error", err.Error())
			return
		}
		defer resp.Body.Close()

		// 간단한 크기 제한 읽기
		body = make([]byte, 0, 4096)
		buf := make([]byte, 1024)
		for {
			n, readErr := resp.Body.Read(buf)
			if n > 0 {
				body = append(body, buf[:n]...)
			}
			if readErr != nil {
				break
			}
		}

		rttrace.Log(ctx, "status", fmt.Sprintf("%d", resp.StatusCode))
	})

	// Region: 파싱 (간단한 바이트 길이 계산으로 대체)
	var bodyLen int
	rttrace.WithRegion(ctx, "parse-html", func() {
		bodyLen = len(body)
		rttrace.Log(ctx, "bodySize", fmt.Sprintf("%d bytes", bodyLen))
	})

	elapsed := time.Since(start)
	rttrace.Log(ctx, "elapsed", elapsed.String())

	return &crawlResult{
		url:     url,
		bodyLen: bodyLen,
		elapsed: elapsed,
	}
}
