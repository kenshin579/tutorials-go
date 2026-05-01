package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/kenshin579/tutorials-go/wiki-permissions/1-acl/backend/domain"
	mw "github.com/kenshin579/tutorials-go/wiki-permissions/1-acl/backend/http/middleware"
	"github.com/kenshin579/tutorials-go/wiki-permissions/1-acl/backend/usecase"
)

// ACLHandler는 페이지 공유 관리(ACL List/Grant/Revoke) 엔드포인트를 처리한다.
type ACLHandler struct{ uc *usecase.ACLUsecase }

// NewACLHandler는 ACLUsecase를 주입받아 ACLHandler를 생성한다.
func NewACLHandler(uc *usecase.ACLUsecase) *ACLHandler { return &ACLHandler{uc: uc} }

// List는 GET /api/pages/:id/acl: owner에게 해당 페이지의 ACL 목록을 반환한다.
func (h *ACLHandler) List(c echo.Context) error {
	pageID, err := parseUintParam(c, "id")
	if err != nil {
		return err
	}
	requester := mw.UserIDFrom(c)
	entries, err := h.uc.List(pageID, requester)
	if err != nil {
		return mapACLError(err)
	}
	return c.JSON(http.StatusOK, entries)
}

type grantRequest struct {
	UserID uint          `json:"user_id"`
	Action domain.Action `json:"action"`
}

// Grant는 POST /api/pages/:id/acl: owner가 (user_id, action) ACL을 추가한다.
func (h *ACLHandler) Grant(c echo.Context) error {
	pageID, err := parseUintParam(c, "id")
	if err != nil {
		return err
	}
	requester := mw.UserIDFrom(c)
	var req grantRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}
	if !req.Action.Valid() {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid action")
	}
	if err := h.uc.Grant(pageID, requester, req.UserID, req.Action); err != nil {
		return mapACLError(err)
	}
	return c.NoContent(http.StatusNoContent)
}

// Revoke는 DELETE /api/pages/:id/acl/:userId: owner가 대상 사용자의 모든 action ACL을 회수한다.
func (h *ACLHandler) Revoke(c echo.Context) error {
	pageID, err := parseUintParam(c, "id")
	if err != nil {
		return err
	}
	userID, err := parseUintParam(c, "userId")
	if err != nil {
		return err
	}
	requester := mw.UserIDFrom(c)
	if err := h.uc.Revoke(pageID, requester, userID); err != nil {
		return mapACLError(err)
	}
	return c.NoContent(http.StatusNoContent)
}

// mapACLError는 ACL usecase 에러를 HTTP 상태로 매핑한다.
// ErrInvalidAction → 400, ErrForbidden → 403, ErrNotFound → 404, 그 외 → 500.
func mapACLError(err error) error {
	if errors.Is(err, usecase.ErrInvalidAction) {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid action")
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
