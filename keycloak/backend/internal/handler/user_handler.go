package handler

import (
	"net/http"
	"strings"

	"github.com/kenshin579/tutorials-go/keycloak/backend/internal/domain"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userUseCase domain.UserUseCase
}

func NewUserHandler(userUseCase domain.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

func (h *UserHandler) GetUserInfo(c echo.Context) error {
	token := extractToken(c)
	if token == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Authorization header required")
	}

	ctx := c.Request().Context()
	user, err := h.userUseCase.GetUserInfo(ctx, token)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
	}

	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) ValidateToken(c echo.Context) error {
	token := extractToken(c)
	if token == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Authorization header required")
	}

	ctx := c.Request().Context()
	valid, err := h.userUseCase.ValidateToken(ctx, token)
	if err != nil || !valid {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
	}

	return c.JSON(http.StatusOK, map[string]bool{"valid": true})
}

func extractToken(c echo.Context) string {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}

	return authHeader
}
