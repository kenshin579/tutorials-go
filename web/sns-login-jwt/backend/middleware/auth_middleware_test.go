package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kenshin579/tutorials-go/web/sns-login-jwt/backend/service"
	"github.com/labstack/echo/v4"
)

func TestJWTAuth_RejectsRefreshToken(t *testing.T) {
	ts := service.NewTokenService("test-secret")
	pair, _ := ts.GenerateTokenPair(7)

	e := echo.New()
	handler := JWTAuth(ts)(func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	// refresh 토큰으로 보호 API 호출 → 거부되어야 함
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+pair.RefreshToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler(c)
	if err == nil {
		t.Fatal("refresh 토큰이 거부되지 않음 (에러 nil)")
	}
	he, ok := err.(*echo.HTTPError)
	if !ok || he.Code != http.StatusUnauthorized {
		t.Fatalf("401 기대, 실제 %v", err)
	}
}

func TestJWTAuth_AcceptsAccessToken(t *testing.T) {
	ts := service.NewTokenService("test-secret")
	pair, _ := ts.GenerateTokenPair(7)

	e := echo.New()
	handler := JWTAuth(ts)(func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+pair.AccessToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := handler(c); err != nil {
		t.Fatalf("access 토큰이 거부됨: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("200 기대, 실제 %d", rec.Code)
	}
}
