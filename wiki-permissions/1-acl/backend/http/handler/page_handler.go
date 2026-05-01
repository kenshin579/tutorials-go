package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/kenshin579/tutorials-go/wiki-permissions/1-acl/backend/domain"
	mw "github.com/kenshin579/tutorials-go/wiki-permissions/1-acl/backend/http/middleware"
	"github.com/kenshin579/tutorials-go/wiki-permissions/1-acl/backend/usecase"
)

// PageHandler는 페이지 목록/상세/수정 엔드포인트를 처리한다.
type PageHandler struct{ uc *usecase.PageUsecase }

// NewPageHandler는 PageUsecase를 주입받아 PageHandler를 생성한다.
func NewPageHandler(uc *usecase.PageUsecase) *PageHandler { return &PageHandler{uc: uc} }

// List는 GET /api/pages: 본인이 access 가능한 페이지 목록을 반환한다.
func (h *PageHandler) List(c echo.Context) error {
	uid := mw.UserIDFrom(c)
	pages, err := h.uc.ListAccessible(uid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, pages)
}

// Get은 GET /api/pages/:id: read 권한을 검증하고 페이지 상세를 반환한다.
func (h *PageHandler) Get(c echo.Context) error {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return err
	}
	uid := mw.UserIDFrom(c)
	page, err := h.uc.Get(id, uid)
	return respondOrError(c, page, err)
}

type updatePageRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// Update는 PUT /api/pages/:id: edit 권한을 검증하고 title/content를 갱신한다.
func (h *PageHandler) Update(c echo.Context) error {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return err
	}
	uid := mw.UserIDFrom(c)
	var req updatePageRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}
	page, err := h.uc.Update(id, uid, req.Title, req.Content)
	return respondOrError(c, page, err)
}

// parseUintParam은 :param을 uint로 파싱한다. 잘못된 값은 400을 반환한다.
func parseUintParam(c echo.Context, key string) (uint, error) {
	raw := c.Param(key)
	v, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}
	return uint(v), nil
}

// respondOrError는 page+err 조합을 적절한 HTTP 응답으로 변환한다.
// usecase.ErrForbidden → 403, domain.ErrNotFound → 404, 그 외 에러 → 500.
func respondOrError(c echo.Context, page *domain.Page, err error) error {
	if err == nil {
		return c.JSON(http.StatusOK, page)
	}
	if errors.Is(err, usecase.ErrForbidden) {
		return echo.NewHTTPError(http.StatusForbidden, "forbidden")
	}
	var nf domain.ErrNotFound
	if errors.As(err, &nf) {
		return echo.NewHTTPError(http.StatusNotFound, nf.Error())
	}
	return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
}
