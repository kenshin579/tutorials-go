package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/kenshin579/tutorials-go/wiki-permissions/2-rbac/backend/domain"
	mw "github.com/kenshin579/tutorials-go/wiki-permissions/2-rbac/backend/http/middleware"
	"github.com/kenshin579/tutorials-go/wiki-permissions/2-rbac/backend/usecase"
)

// RoleHandler는 admin 전용 사용자 role 관리 엔드포인트를 처리한다.
type RoleHandler struct{ uc *usecase.RoleUsecase }

// NewRoleHandler는 RoleUsecase를 주입받아 RoleHandler를 생성한다.
func NewRoleHandler(uc *usecase.RoleUsecase) *RoleHandler { return &RoleHandler{uc: uc} }

// ListUsers는 GET /api/users: admin이 모든 사용자(+ 각자 roles)를 조회한다.
func (h *RoleHandler) ListUsers(c echo.Context) error {
	requester := mw.UserIDFrom(c)
	users, err := h.uc.ListUsers(requester)
	if err != nil {
		return mapRoleError(err)
	}
	return c.JSON(http.StatusOK, users)
}

// ListRoles는 GET /api/roles: admin이 모든 role(+ 각 role의 permissions)을 조회한다.
func (h *RoleHandler) ListRoles(c echo.Context) error {
	requester := mw.UserIDFrom(c)
	roles, err := h.uc.ListRoles(requester)
	if err != nil {
		return mapRoleError(err)
	}
	return c.JSON(http.StatusOK, roles)
}

type assignRoleRequest struct {
	RoleID uint `json:"role_id"`
}

// Assign은 POST /api/users/:id/roles: admin이 대상 사용자에게 role을 부여한다.
func (h *RoleHandler) Assign(c echo.Context) error {
	targetUserID, err := parseUintParam(c, "id")
	if err != nil {
		return err
	}
	requester := mw.UserIDFrom(c)
	var req assignRoleRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}
	if err := h.uc.AssignRole(requester, targetUserID, req.RoleID); err != nil {
		return mapRoleError(err)
	}
	return c.NoContent(http.StatusNoContent)
}

// Revoke는 DELETE /api/users/:id/roles/:roleId: admin이 대상 사용자에게서 role을 회수한다.
func (h *RoleHandler) Revoke(c echo.Context) error {
	targetUserID, err := parseUintParam(c, "id")
	if err != nil {
		return err
	}
	roleID, err := parseUintParam(c, "roleId")
	if err != nil {
		return err
	}
	requester := mw.UserIDFrom(c)
	if err := h.uc.RevokeRole(requester, targetUserID, roleID); err != nil {
		return mapRoleError(err)
	}
	return c.NoContent(http.StatusNoContent)
}

// mapRoleError는 RoleUsecase 에러를 HTTP 상태로 매핑한다.
func mapRoleError(err error) error {
	if errors.Is(err, usecase.ErrForbidden) {
		return echo.NewHTTPError(http.StatusForbidden, "forbidden")
	}
	var nf domain.ErrNotFound
	if errors.As(err, &nf) {
		return echo.NewHTTPError(http.StatusNotFound, nf.Error())
	}
	return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
}
