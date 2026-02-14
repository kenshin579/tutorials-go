package crawler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sync"
	"time"
)

// Result - 크롤링 결과
type Result struct {
	URL   string
	Title string
	Links []string
	Err   error
}

// Crawler - 동시성 웹 크롤러
type Crawler struct {
	client      *http.Client
	maxWorkers  int
	rateLimit   time.Duration
	visited     sync.Map
	results     []Result
	resultsMu   sync.Mutex
}

// NewCrawler - 새 크롤러 생성
func NewCrawler(maxWorkers int, rateLimit time.Duration) *Crawler {
	return &Crawler{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		maxWorkers: maxWorkers,
		rateLimit:  rateLimit,
	}
}

// Crawl - 주어진 URL들을 동시에 크롤링
func (c *Crawler) Crawl(ctx context.Context, urls []string) []Result {
	sem := make(chan struct{}, c.maxWorkers) // semaphore
	var wg sync.WaitGroup

	// rate limiter
	ticker := time.NewTicker(c.rateLimit)
	defer ticker.Stop()

	for _, u := range urls {
		// 이미 방문한 URL은 건너뜀
		if _, loaded := c.visited.LoadOrStore(u, true); loaded {
			continue
		}

		select {
		case <-ctx.Done():
			break
		case <-ticker.C:
			// rate limit 대기
		}

		sem <- struct{}{} // semaphore 획득
		wg.Add(1)

		go func() {
			defer wg.Done()
			defer func() { <-sem }() // semaphore 해제

			result := c.fetch(ctx, u)
			c.resultsMu.Lock()
			c.results = append(c.results, result)
			c.resultsMu.Unlock()
		}()
	}

	wg.Wait()
	return c.results
}

// fetch - URL에서 콘텐츠 가져오기
func (c *Crawler) fetch(ctx context.Context, url string) Result {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Result{URL: url, Err: fmt.Errorf("request 생성 실패: %w", err)}
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return Result{URL: url, Err: fmt.Errorf("fetch 실패: %w", err)}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{URL: url, Err: fmt.Errorf("body 읽기 실패: %w", err)}
	}

	title := extractTitle(string(body))
	links := extractLinks(string(body))

	return Result{
		URL:   url,
		Title: title,
		Links: links,
	}
}

// extractTitle - HTML에서 title 추출
func extractTitle(html string) string {
	re := regexp.MustCompile(`<title>(.*?)</title>`)
	matches := re.FindStringSubmatch(html)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// extractLinks - HTML에서 href 링크 추출
func extractLinks(html string) []string {
	re := regexp.MustCompile(`href="(https?://[^"]+)"`)
	matches := re.FindAllStringSubmatch(html, -1)
	var links []string
	for _, m := range matches {
		if len(m) > 1 {
			links = append(links, m[1])
		}
	}
	return links
}
