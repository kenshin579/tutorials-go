package handler

import (
	"net/http"
	"time"

	customMiddleware "github.com/kenshin579/tutorials-go/web/sns-login-session/backend/middleware"
	"github.com/kenshin579/tutorials-go/web/sns-login-session/backend/service"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService *service.AuthService
	frontendURL string
}

func NewAuthHandler(authService *service.AuthService, frontendURL string) *AuthHandler {
	return &AuthHandler{authService: authService, frontendURL: frontendURL}
}

// GET /api/auth/:provider/url
func (h *AuthHandler) GetAuthURL(c echo.Context) error {
	url, err := h.authService.GetAuthURL(c.Param("provider"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]string{"url": url})
}

// GET /api/auth/session/callback — Google이 직접 redirect (redirect_uri = 백엔드)
func (h *AuthHandler) HandleCallback(c echo.Context) error {
	const providerName = "google" // 세션 버전은 Google 전용 서버 redirect 콜백
	code := c.QueryParam("code")
	state := c.QueryParam("state")
	if code == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "code 파라미터가 필요합니다")
	}

	sess, _, err := h.authService.HandleCallback(c.Request().Context(), providerName, code, state)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// HttpOnly 세션 쿠키 설정 후 프론트로 redirect
	c.SetCookie(&http.Cookie{
		Name:     customMiddleware.SessionCookieName,
		Value:    sess.ID,
		Path:     "/",
		Expires:  sess.ExpiresAt,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		// 개발(동일 사이트, localhost)에서는 Lax로 충분.
		// 프로덕션에서 프론트/백엔드가 다른 사이트면 SameSite=None + Secure 필요 (HTTPS).
		// Secure: true,
	})
	return c.Redirect(http.StatusFound, h.frontendURL)
}

// POST /api/auth/logout — 세션 삭제 + 쿠키 만료 (서버측 무효화)
func (h *AuthHandler) Logout(c echo.Context) error {
	if cookie, err := c.Cookie(customMiddleware.SessionCookieName); err == nil {
		if err := h.authService.Logout(cookie.Value); err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, "로그아웃 처리 실패")
		}
	}
	c.SetCookie(&http.Cookie{
		Name:     customMiddleware.SessionCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	return c.JSON(http.StatusOK, map[string]string{"message": "로그아웃 성공"})
}
