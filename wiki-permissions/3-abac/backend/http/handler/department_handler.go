package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/kenshin579/tutorials-go/wiki-permissions/3-abac/backend/domain"
)

// DepartmentHandler는 인증된 사용자에게 부서 목록을 노출한다.
// (페이지 생성 시 부서 선택 등에 참조 — 인증만 있으면 누구나 조회 가능.)
type DepartmentHandler struct{ repo domain.DepartmentRepository }

func NewDepartmentHandler(repo domain.DepartmentRepository) *DepartmentHandler {
	return &DepartmentHandler{repo: repo}
}

func (h *DepartmentHandler) List(c echo.Context) error {
	ds, err := h.repo.List()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, ds)
}
