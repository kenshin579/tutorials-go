package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/kenshin579/tutorials-go/wiki-permissions/2-rbac/backend/usecase"
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

type roleSummary struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type loginResponse struct {
	Token       string        `json:"token"`
	User        loginUser     `json:"user"`
	Permissions []string      `json:"permissions"`
	Roles       []roleSummary `json:"roles"`
}

// Login은 POST /auth/login: email/password 검증 후 access token + 사용자 + 권한/역할을 반환한다.
// 1편(ACL) 응답에는 token + user만 있었지만, RBAC에서는 클라이언트가 PermissionGate에 사용할
// permissions와 roles까지 한 응답에 담는다.
func (h *AuthHandler) Login(c echo.Context) error {
	var req loginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}
	res, err := h.auth.Login(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidCredentials) {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	permKeys := make([]string, 0, len(res.Permissions))
	for _, p := range res.Permissions {
		permKeys = append(permKeys, p.Key())
	}
	roles := make([]roleSummary, 0, len(res.User.Roles))
	for _, r := range res.User.Roles {
		roles = append(roles, roleSummary{ID: r.ID, Name: r.Name})
	}
	return c.JSON(http.StatusOK, loginResponse{
		Token:       res.Token,
		User:        loginUser{ID: res.User.ID, Email: res.User.Email, Name: res.User.Name},
		Permissions: permKeys,
		Roles:       roles,
	})
}
