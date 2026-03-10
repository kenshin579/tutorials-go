package handler

import (
	"net/http"
	"strconv"

	"github.com/kenshin579/tutorials-go/rbac/backend/domain"
	"github.com/kenshin579/tutorials-go/rbac/backend/usecase"
	"github.com/labstack/echo/v4"
)

type RbacHandler struct {
	rbacUsecase usecase.RbacUsecase
}

func NewRbacHandler(ru usecase.RbacUsecase) *RbacHandler {
	return &RbacHandler{rbacUsecase: ru}
}

// Role handlers

func (h *RbacHandler) ListRoles(c echo.Context) error {
	roles, err := h.rbacUsecase.GetAllRoles()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, roles)
}

type createRoleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (h *RbacHandler) CreateRole(c echo.Context) error {
	var req createRoleRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	role := &domain.Role{
		Name:        req.Name,
		Description: req.Description,
	}
	if err := h.rbacUsecase.CreateRole(role); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, role)
}

type updateRoleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (h *RbacHandler) UpdateRole(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	role, err := h.rbacUsecase.GetRoleByID(uint(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "role not found")
	}

	var req updateRoleRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	if req.Name != "" {
		role.Name = req.Name
	}
	if req.Description != "" {
		role.Description = req.Description
	}

	if err := h.rbacUsecase.UpdateRole(role); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, role)
}

func (h *RbacHandler) DeleteRole(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	if err := h.rbacUsecase.DeleteRole(uint(id)); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "role deleted"})
}

type assignPermissionRequest struct {
	PermissionID uint `json:"permission_id"`
}

func (h *RbacHandler) AssignPermission(c echo.Context) error {
	roleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid role id")
	}

	var req assignPermissionRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	if err := h.rbacUsecase.AssignPermission(uint(roleID), req.PermissionID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "permission assigned"})
}

func (h *RbacHandler) RemovePermission(c echo.Context) error {
	roleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid role id")
	}

	permID, err := strconv.ParseUint(c.Param("permId"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid permission id")
	}

	if err := h.rbacUsecase.RemovePermission(uint(roleID), uint(permID)); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "permission removed"})
}

// Permission handlers

func (h *RbacHandler) ListPermissions(c echo.Context) error {
	permissions, err := h.rbacUsecase.GetAllPermissions()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, permissions)
}
