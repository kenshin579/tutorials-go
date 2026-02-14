package go1_25

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CrossOriginProtection_동일출처_허용(t *testing.T) {
	cop := http.NewCrossOriginProtection()

	handler := cop.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "OK")
	}))

	ts := httptest.NewServer(handler)
	defer ts.Close()

	// 동일 출처 요청 (Origin 헤더 없음 = 동일 출처로 간주)
	resp, err := http.Post(ts.URL+"/api/data", "application/json", strings.NewReader("{}"))
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func Test_CrossOriginProtection_크로스오리진_차단(t *testing.T) {
	cop := http.NewCrossOriginProtection()

	handler := cop.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "OK")
	}))

	ts := httptest.NewServer(handler)
	defer ts.Close()

	// 크로스 오리진 POST 요청 (Sec-Fetch-Site: cross-site)
	req, err := http.NewRequest("POST", ts.URL+"/api/data", strings.NewReader("{}"))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://evil.com")
	req.Header.Set("Sec-Fetch-Site", "cross-site")

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusForbidden, resp.StatusCode, "크로스 오리진 POST는 차단되어야 한다")
}

func Test_CrossOriginProtection_GET_허용(t *testing.T) {
	cop := http.NewCrossOriginProtection()

	handler := cop.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "OK")
	}))

	ts := httptest.NewServer(handler)
	defer ts.Close()

	// GET 요청은 안전한 메서드로 항상 허용
	req, err := http.NewRequest("GET", ts.URL+"/api/data", nil)
	assert.NoError(t, err)
	req.Header.Set("Origin", "https://evil.com")
	req.Header.Set("Sec-Fetch-Site", "cross-site")

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode, "GET 요청은 항상 허용되어야 한다")
}

func Test_CrossOriginProtection_신뢰_출처_허용(t *testing.T) {
	cop := http.NewCrossOriginProtection()
	err := cop.AddTrustedOrigin("https://trusted.com")
	assert.NoError(t, err)

	handler := cop.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "OK")
	}))

	ts := httptest.NewServer(handler)
	defer ts.Close()

	// 신뢰 출처에서의 POST 요청 → 허용
	req, err := http.NewRequest("POST", ts.URL+"/api/data", strings.NewReader("{}"))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://trusted.com")
	req.Header.Set("Sec-Fetch-Site", "cross-site")

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode, "신뢰 출처는 허용되어야 한다")
}
