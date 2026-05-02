package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/kenshin579/tutorials-go/wiki-permissions/2-rbac/backend/domain"
	mw "github.com/kenshin579/tutorials-go/wiki-permissions/2-rbac/backend/http/middleware"
	"github.com/kenshin579/tutorials-go/wiki-permissions/2-rbac/backend/usecase"
)

// PageHandler는 페이지 CRUD 엔드포인트를 처리한다.
// 1편(ACL)과 차이: ACL grant/revoke 라우트 대신 Create/Delete 라우트가 있다.
type PageHandler struct{ uc *usecase.PageUsecase }

// NewPageHandler는 PageUsecase를 주입받아 PageHandler를 생성한다.
func NewPageHandler(uc *usecase.PageUsecase) *PageHandler { return &PageHandler{uc: uc} }

// List는 GET /api/pages: pages:read 권한이 있으면 모든 페이지를 반환한다.
func (h *PageHandler) List(c echo.Context) error {
	uid := mw.UserIDFrom(c)
	pages, err := h.uc.List(uid)
	if err != nil {
		return mapPageError(err)
	}
	return c.JSON(http.StatusOK, pages)
}

// Get은 GET /api/pages/:id: pages:read 권한 검증 후 페이지 상세를 반환한다.
func (h *PageHandler) Get(c echo.Context) error {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return err
	}
	uid := mw.UserIDFrom(c)
	page, err := h.uc.Get(id, uid)
	if err != nil {
		return mapPageError(err)
	}
	return c.JSON(http.StatusOK, page)
}

type createPageRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// Create는 POST /api/pages: pages:create 권한 검증 후 새 페이지를 생성한다 (owner = 요청자).
func (h *PageHandler) Create(c echo.Context) error {
	uid := mw.UserIDFrom(c)
	var req createPageRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}
	page, err := h.uc.Create(uid, req.Title, req.Content)
	if err != nil {
		return mapPageError(err)
	}
	return c.JSON(http.StatusCreated, page)
}

type updatePageRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// Update는 PUT /api/pages/:id: pages:edit 권한 검증 후 title/content를 갱신한다.
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
	if err != nil {
		return mapPageError(err)
	}
	return c.JSON(http.StatusOK, page)
}

// Delete는 DELETE /api/pages/:id: pages:delete 권한(admin)이 있을 때 페이지를 삭제한다.
func (h *PageHandler) Delete(c echo.Context) error {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return err
	}
	uid := mw.UserIDFrom(c)
	if err := h.uc.Delete(id, uid); err != nil {
		return mapPageError(err)
	}
	return c.NoContent(http.StatusNoContent)
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

// mapPageError는 usecase 에러를 HTTP 상태로 매핑한다.
// usecase.ErrForbidden → 403, domain.ErrNotFound → 404, 그 외 → 500.
func mapPageError(err error) error {
	if errors.Is(err, usecase.ErrForbidden) {
		return echo.NewHTTPError(http.StatusForbidden, "forbidden")
	}
	var nf domain.ErrNotFound
	if errors.As(err, &nf) {
		return echo.NewHTTPError(http.StatusNotFound, nf.Error())
	}
	return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
}
