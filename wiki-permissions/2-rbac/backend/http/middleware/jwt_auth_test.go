package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	jwthelper "github.com/kenshin579/tutorials-go/wiki-permissions/2-rbac/backend/pkg/jwt"
)

func TestJWTAuth_ValidToken_InjectsUserID(t *testing.T) {
	secret := "s"
	tok, err := jwthelper.Issue(42, secret, time.Hour)
	require.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	var got uint
	handler := func(c echo.Context) error {
		got = UserIDFrom(c)
		return c.String(http.StatusOK, "")
	}

	mw := JWTAuth(secret)
	require.NoError(t, mw(handler)(c))
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, uint(42), got)
}

func TestJWTAuth_MissingHeader(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := func(c echo.Context) error { return c.String(http.StatusOK, "") }
	mw := JWTAuth("s")
	err := mw(handler)(c)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	assert.Equal(t, http.StatusUnauthorized, httpErr.Code)
}

func TestJWTAuth_InvalidToken(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer not-a-token")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := func(c echo.Context) error { return c.String(http.StatusOK, "") }
	mw := JWTAuth("s")
	err := mw(handler)(c)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	assert.Equal(t, http.StatusUnauthorized, httpErr.Code)
}
