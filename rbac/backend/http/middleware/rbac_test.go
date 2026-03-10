package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kenshin579/tutorials-go/rbac/backend/domain"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// mockPermissionRepository is a hand-rolled mock that satisfies domain.PermissionRepository.
type mockPermissionRepository struct {
	permissions []domain.Permission
	err         error
}

func (m *mockPermissionRepository) FindAll() ([]domain.Permission, error) {
	return m.permissions, m.err
}

func (m *mockPermissionRepository) FindByUserID(_ uint) ([]domain.Permission, error) {
	return m.permissions, m.err
}

// okHandler is a trivial next handler that always responds 200 OK.
var okHandler echo.HandlerFunc = func(c echo.Context) error {
	return c.String(http.StatusOK, "ok")
}

// newContext creates an Echo context with the given user_id already set,
// simulating the output of a preceding JWT-auth middleware.
func newContext(e *echo.Echo, userID uint) echo.Context {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", userID)
	return c
}

func TestRequirePermission_UserHasPermission(t *testing.T) {
	e := echo.New()

	repo := &mockPermissionRepository{
		permissions: []domain.Permission{
			{ID: 1, Resource: "product", Action: "read"},
			{ID: 2, Resource: "order", Action: "write"},
		},
	}

	mw := RequirePermission("product:read", repo)
	handler := mw(okHandler)

	c := newContext(e, 42)
	err := handler(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, c.Response().Status)
}

func TestRequirePermission_UserLacksPermission(t *testing.T) {
	e := echo.New()

	repo := &mockPermissionRepository{
		permissions: []domain.Permission{
			{ID: 1, Resource: "order", Action: "write"},
		},
	}

	mw := RequirePermission("product:read", repo)
	handler := mw(okHandler)

	c := newContext(e, 42)
	err := handler(c)

	assert.Error(t, err)
	var httpErr *echo.HTTPError
	assert.True(t, errors.As(err, &httpErr))
	assert.Equal(t, http.StatusForbidden, httpErr.Code)
}

func TestRequirePermission_RepositoryError(t *testing.T) {
	e := echo.New()

	repo := &mockPermissionRepository{
		err: errors.New("db connection lost"),
	}

	mw := RequirePermission("product:read", repo)
	handler := mw(okHandler)

	c := newContext(e, 42)
	err := handler(c)

	assert.Error(t, err)
	var httpErr *echo.HTTPError
	assert.True(t, errors.As(err, &httpErr))
	assert.Equal(t, http.StatusInternalServerError, httpErr.Code)
}
