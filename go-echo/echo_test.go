package go_echo

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestEchoHandler(t *testing.T) {
	e := echo.New()

	// HandlerFunc
	e.GET("/ok", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	c, b := request(http.MethodGet, "/ok", e)
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, "OK", b)
}

func TestEchoServeHTTPPathEncoding(t *testing.T) {
	e := echo.New()
	e.GET("/with/slash", func(c echo.Context) error {
		return c.String(http.StatusOK, "/with/slash")
	})
	e.GET("/:id", func(c echo.Context) error {
		return c.String(http.StatusOK, c.Param("id"))
	})

	var testCases = []struct {
		name         string
		whenURL      string
		expectURL    string
		expectStatus int
	}{
		{
			name:         "url with encoding is not decoded for routing",
			whenURL:      "/with%2Fslash",
			expectURL:    "with%2Fslash", // `%2F` is not decoded to `/` for routing
			expectStatus: http.StatusOK,
		},
		{
			name:         "url without encoding is used as is",
			whenURL:      "/with/slash",
			expectURL:    "/with/slash",
			expectStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.whenURL, nil)
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectStatus, rec.Code)
			assert.Equal(t, tc.expectURL, rec.Body.String())
		})
	}
}

func request(method, path string, e *echo.Echo) (int, string) {
	req := httptest.NewRequest(method, path, nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec.Code, rec.Body.String()
}
