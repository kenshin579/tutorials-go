package handler

import (
	"net/http"

	"github.com/kenshin579/tutorials-go/rbac/backend/domain"
	"github.com/kenshin579/tutorials-go/rbac/backend/usecase"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authUsecase usecase.AuthUsecase
	permRepo    domain.PermissionRepository
	jwtSecret   string
}

func NewAuthHandler(au usecase.AuthUsecase, permRepo domain.PermissionRepository, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		authUsecase: au,
		permRepo:    permRepo,
		jwtSecret:   jwtSecret,
	}
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type loginResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	User         interface{} `json:"user"`
}

type userInfo struct {
	ID          uint     `json:"id"`
	Email       string   `json:"email"`
	Name        string   `json:"name"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req registerRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	user, err := h.authUsecase.Register(req.Email, req.Password, req.Name)
	if err != nil {
		if err == usecase.ErrEmailAlreadyExists {
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, user)
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req loginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	pair, user, err := h.authUsecase.Login(req.Email, req.Password, h.jwtSecret)
	if err != nil {
		if err == usecase.ErrInvalidCredentials {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Get user permissions
	perms, _ := h.permRepo.FindByUserID(user.ID)
	permKeys := make([]string, len(perms))
	for i, p := range perms {
		permKeys[i] = p.Key()
	}

	roles := make([]string, len(user.Roles))
	for i, r := range user.Roles {
		roles[i] = r.Name
	}

	return c.JSON(http.StatusOK, loginResponse{
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
		User: userInfo{
			ID:          user.ID,
			Email:       user.Email,
			Name:        user.Name,
			Roles:       roles,
			Permissions: permKeys,
		},
	})
}

func (h *AuthHandler) Refresh(c echo.Context) error {
	var req refreshRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	pair, err := h.authUsecase.Refresh(req.RefreshToken, h.jwtSecret)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid refresh token")
	}

	return c.JSON(http.StatusOK, pair)
}

func (h *AuthHandler) Logout(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "logged out"})
}
