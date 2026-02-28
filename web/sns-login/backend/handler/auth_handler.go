package handler

import (
	"net/http"

	"github.com/kenshin579/tutorials-go/web/sns-login/backend/service"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// GetAuthURL은 OAuth 인증 URL을 반환한다
// GET /api/auth/:provider/url
func (h *AuthHandler) GetAuthURL(c echo.Context) error {
	providerName := c.Param("provider")

	url, err := h.authService.GetAuthURL(providerName)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{
		"url": url,
	})
}

// HandleCallback은 OAuth 콜백을 처리한다
// GET /api/auth/:provider/callback
func (h *AuthHandler) HandleCallback(c echo.Context) error {
	providerName := c.Param("provider")
	code := c.QueryParam("code")
	state := c.QueryParam("state")

	if code == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "code 파라미터가 필요합니다")
	}

	tokens, user, err := h.authService.HandleCallback(c.Request().Context(), providerName, code, state)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"tokens": tokens,
		"user":   user,
	})
}

// RefreshToken은 새 토큰 쌍을 발급한다
// POST /api/auth/refresh
func (h *AuthHandler) RefreshToken(c echo.Context) error {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "잘못된 요청")
	}

	tokens, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	return c.JSON(http.StatusOK, tokens)
}

// Logout은 로그아웃을 처리한다
// POST /api/auth/logout
func (h *AuthHandler) Logout(c echo.Context) error {
	// 클라이언트 측에서 토큰 삭제로 처리
	return c.JSON(http.StatusOK, map[string]string{
		"message": "로그아웃 성공",
	})
}
