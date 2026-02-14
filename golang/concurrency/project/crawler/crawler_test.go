package crawler

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestCrawlerBasic - 기본 크롤링 테스트
func TestCrawlerBasic(t *testing.T) {
	// 테스트 서버 구성
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `<html><title>Home</title><body>
			<a href="/page1">Page 1</a>
			<a href="/page2">Page 2</a>
		</body></html>`)
	})
	mux.HandleFunc("/page1", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `<html><title>Page 1</title><body>Content 1</body></html>`)
	})
	mux.HandleFunc("/page2", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `<html><title>Page 2</title><body>Content 2</body></html>`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	c := NewCrawler(3, 1*time.Millisecond)
	ctx := context.Background()

	urls := []string{
		server.URL + "/",
		server.URL + "/page1",
		server.URL + "/page2",
	}

	results := c.Crawl(ctx, urls)

	assert.Equal(t, 3, len(results))

	for _, r := range results {
		assert.NoError(t, r.Err)
		assert.NotEmpty(t, r.Title)
		t.Logf("URL: %s, Title: %s", r.URL, r.Title)
	}
}

// TestCrawlerWithCancel - context 취소로 크롤링 중단
func TestCrawlerWithCancel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond) // 느린 서버
		fmt.Fprint(w, `<html><title>Slow</title></html>`)
	}))
	defer server.Close()

	c := NewCrawler(2, 1*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	urls := []string{
		server.URL + "/1",
		server.URL + "/2",
		server.URL + "/3",
		server.URL + "/4",
		server.URL + "/5",
	}

	results := c.Crawl(ctx, urls)

	// timeout으로 인해 일부만 처리됨 (또는 모두 에러)
	t.Logf("처리된 결과 수: %d", len(results))
	for _, r := range results {
		if r.Err != nil {
			t.Logf("URL: %s, Error: %v", r.URL, r.Err)
		}
	}
}

// TestCrawlerDuplicateURL - 중복 URL은 한 번만 크롤링
func TestCrawlerDuplicateURL(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		fmt.Fprint(w, `<html><title>Test</title></html>`)
	}))
	defer server.Close()

	c := NewCrawler(3, 1*time.Millisecond)
	ctx := context.Background()

	urls := []string{
		server.URL + "/page",
		server.URL + "/page", // 중복
		server.URL + "/page", // 중복
	}

	results := c.Crawl(ctx, urls)
	assert.Equal(t, 1, len(results)) // 1번만 크롤링
	assert.Equal(t, 1, callCount)
}

// TestCrawlerLinkExtraction - 링크 추출 테스트
func TestCrawlerLinkExtraction(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `<html><title>Links</title><body>
			<a href="https://example.com/a">A</a>
			<a href="https://example.com/b">B</a>
			<a href="/relative">Relative</a>
		</body></html>`)
	}))
	defer server.Close()

	c := NewCrawler(1, 1*time.Millisecond)
	ctx := context.Background()

	results := c.Crawl(ctx, []string{server.URL + "/"})

	assert.Equal(t, 1, len(results))
	assert.Equal(t, "Links", results[0].Title)
	assert.Equal(t, 2, len(results[0].Links)) // https:// 링크만 추출
	assert.Equal(t, "https://example.com/a", results[0].Links[0])
	assert.Equal(t, "https://example.com/b", results[0].Links[1])
}

// TestExtractTitle - title 추출 단위 테스트
func TestExtractTitle(t *testing.T) {
	tests := []struct {
		html     string
		expected string
	}{
		{`<html><title>Hello</title></html>`, "Hello"},
		{`<html><title>Go 동시성</title></html>`, "Go 동시성"},
		{`<html><body>no title</body></html>`, ""},
	}

	for _, tt := range tests {
		result := extractTitle(tt.html)
		assert.Equal(t, tt.expected, result)
	}
}

// TestExtractLinks - 링크 추출 단위 테스트
func TestExtractLinks(t *testing.T) {
	html := `<a href="https://go.dev">Go</a> <a href="http://example.com">Example</a> <a href="/local">Local</a>`
	links := extractLinks(html)

	assert.Equal(t, 2, len(links))
	assert.Equal(t, "https://go.dev", links[0])
	assert.Equal(t, "http://example.com", links[1])
}
