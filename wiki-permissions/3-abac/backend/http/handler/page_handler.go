package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/kenshin579/tutorials-go/wiki-permissions/3-abac/backend/domain"
	mw "github.com/kenshin579/tutorials-go/wiki-permissions/3-abac/backend/http/middleware"
	"github.com/kenshin579/tutorials-go/wiki-permissions/3-abac/backend/usecase"
)

// PageHandler는 페이지 List/Get/Update 엔드포인트를 처리한다.
// Get 응답에는 ABAC Decision (read/edit 가능 여부와 reason)을 함께 내려보내 클라이언트가
// 사용자에게 "왜 허용/거부됐는지"를 자연스럽게 표시할 수 있도록 한다.
type PageHandler struct{ uc *usecase.PageUsecase }

func NewPageHandler(uc *usecase.PageUsecase) *PageHandler { return &PageHandler{uc: uc} }

// List는 GET /api/pages: 호출자가 read 가능한 페이지만 반환한다 (정책 평가 통과).
func (h *PageHandler) List(c echo.Context) error {
	uid := mw.UserIDFrom(c)
	pages, err := h.uc.List(uid)
	if err != nil {
		return mapPageError(err)
	}
	return c.JSON(http.StatusOK, pages)
}

// Get은 GET /api/pages/:id: 페이지 상세 + Decision(can_read, can_edit)을 반환한다.
// read 거부 시 페이지 자체를 노출하지 않고 403만 반환한다.
func (h *PageHandler) Get(c echo.Context) error {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return err
	}
	uid := mw.UserIDFrom(c)
	res, err := h.uc.Get(id, uid)
	if err != nil {
		return mapPageError(err)
	}
	return c.JSON(http.StatusOK, res)
}

type updatePageRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// Update는 PUT /api/pages/:id: edit 정책 평가 후 title/content를 갱신한다.
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

func parseUintParam(c echo.Context, key string) (uint, error) {
	raw := c.Param(key)
	v, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}
	return uint(v), nil
}

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
