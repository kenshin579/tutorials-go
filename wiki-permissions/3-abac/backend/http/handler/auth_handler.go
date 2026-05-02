package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/kenshin579/tutorials-go/wiki-permissions/3-abac/backend/usecase"
)

// AuthHandler는 로그인 엔드포인트를 처리한다.
type AuthHandler struct{ auth *usecase.AuthUsecase }

func NewAuthHandler(auth *usecase.AuthUsecase) *AuthHandler { return &AuthHandler{auth: auth} }

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type departmentSummary struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type loginUser struct {
	ID             uint               `json:"id"`
	Email          string             `json:"email"`
	Name           string             `json:"name"`
	Department     *departmentSummary `json:"department,omitempty"`
	EmploymentType string             `json:"employment_type"`
}

type loginResponse struct {
	Token string    `json:"token"`
	User  loginUser `json:"user"`
}

// Login은 POST /auth/login: email/password 검증 후 token + 사용자(속성 포함)를 반환한다.
//
// 1편(ACL): token + user. 2편(RBAC): + permissions/roles. 3편(ABAC): user 자체에 department/employment_type
// 속성이 들어있어 별도 권한 응답 불필요. 권한 평가는 매 요청 시 ABAC 정책 평가기가 user/page 속성을 결합 평가.
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
	var dept *departmentSummary
	if u.Department != nil {
		dept = &departmentSummary{ID: u.Department.ID, Name: u.Department.Name}
	}
	return c.JSON(http.StatusOK, loginResponse{
		Token: tok,
		User: loginUser{
			ID:             u.ID,
			Email:          u.Email,
			Name:           u.Name,
			Department:     dept,
			EmploymentType: string(u.EmploymentType),
		},
	})
}
