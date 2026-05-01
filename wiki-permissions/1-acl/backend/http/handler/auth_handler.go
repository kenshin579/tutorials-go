package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/kenshin579/tutorials-go/wiki-permissions/1-acl/backend/usecase"
)

// AuthHandler는 로그인 엔드포인트를 처리한다.
type AuthHandler struct {
	auth *usecase.AuthUsecase
}

// NewAuthHandler는 AuthUsecase를 주입받아 AuthHandler를 생성한다.
func NewAuthHandler(auth *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{auth: auth}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginUser struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type loginResponse struct {
	Token string    `json:"token"`
	User  loginUser `json:"user"`
}

// Login는 POST /auth/login: email/password 검증 후 access token + 사용자 정보를 반환한다.
// 자격 증명 실패 시 401, 잘못된 본문은 400을 반환한다.
func (h *AuthHandler) Login(c echo.Context) error {
	var req loginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}
	tok, u, err := h.auth.Login(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidCredentials) {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, loginResponse{
		Token: tok,
		User:  loginUser{ID: u.ID, Email: u.Email, Name: u.Name},
	})
}
