package custom

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var testSigningKey = []byte("test-secret-key")

func generateTestToken(t *testing.T, claims *Claims) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(testSigningKey)
	assert.NoError(t, err)
	return tokenString
}

func TestJWTAuth_ValidToken(t *testing.T) {
	tokenString := generateTestToken(t, &Claims{
		UserID:   "user-123",
		Username: "testuser",
		Role:     "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	e := echo.New()
	e.Use(JWTAuth(JWTConfig{SigningKey: testSigningKey}))
	e.GET("/test", func(c echo.Context) error {
		claims := c.Get("user").(*Claims)
		return c.JSON(http.StatusOK, map[string]string{
			"user_id":  claims.UserID,
			"username": claims.Username,
			"role":     claims.Role,
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "user-123")
	assert.Contains(t, rec.Body.String(), "testuser")
}

func TestJWTAuth_ExpiredToken(t *testing.T) {
	tokenString := generateTestToken(t, &Claims{
		UserID: "user-123",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // 만료됨
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	})

	e := echo.New()
	e.Use(JWTAuth(JWTConfig{SigningKey: testSigningKey}))
	e.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestJWTAuth_MissingToken(t *testing.T) {
	e := echo.New()
	e.Use(JWTAuth(JWTConfig{SigningKey: testSigningKey}))
	e.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestJWTAuth_InvalidFormat(t *testing.T) {
	e := echo.New()
	e.Use(JWTAuth(JWTConfig{SigningKey: testSigningKey}))
	e.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "InvalidFormat token123")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestJWTAuth_Skipper(t *testing.T) {
	e := echo.New()
	e.Use(JWTAuth(JWTConfig{
		SigningKey: testSigningKey,
		Skipper: func(c echo.Context) bool {
			return c.Path() == "/health"
		},
	}))
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	// 토큰 없이 /health 접근 → Skipper로 통과
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "ok", rec.Body.String())
}

func TestJWTAuth_CustomContextKey(t *testing.T) {
	tokenString := generateTestToken(t, &Claims{
		UserID:   "user-456",
		Username: "admin",
		Role:     "superadmin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	})

	e := echo.New()
	e.Use(JWTAuth(JWTConfig{
		SigningKey:  testSigningKey,
		ContextKey: "auth_user",
	}))
	e.GET("/test", func(c echo.Context) error {
		claims := c.Get("auth_user").(*Claims)
		return c.String(http.StatusOK, fmt.Sprintf("hello %s", claims.Username))
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "hello admin", rec.Body.String())
}

func TestJWTAuth_WrongSigningKey(t *testing.T) {
	// 다른 키로 서명한 토큰
	wrongKey := []byte("wrong-key")
	claims := &Claims{
		UserID: "user-123",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(wrongKey)

	e := echo.New()
	e.Use(JWTAuth(JWTConfig{SigningKey: testSigningKey}))
	e.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}
