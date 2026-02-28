package handler

import (
	"net/http"

	"github.com/kenshin579/tutorials-go/web/sns-login/backend/service"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	authService *service.AuthService
}

func NewUserHandler(authService *service.AuthService) *UserHandler {
	return &UserHandler{authService: authService}
}

// GetMe는 현재 로그인한 사용자 정보를 반환한다
// GET /api/user/me
func (h *UserHandler) GetMe(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	user, err := h.authService.GetUser(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "사용자를 찾을 수 없습니다")
	}

	return c.JSON(http.StatusOK, user)
}
