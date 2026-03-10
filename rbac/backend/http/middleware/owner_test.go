package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// newOwnerContext creates an Echo context with user_id, roles, and an ":id" path param pre-set.
func newOwnerContext(e *echo.Echo, userID uint, roles []string, resourceID string) echo.Context {
	req := httptest.NewRequest(http.MethodGet, "/resources/"+resourceID, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(resourceID)
	c.Set("user_id", userID)
	c.Set("roles", roles)
	return c
}

// newNilDB returns a zero-value *gorm.DB whose underlying sql.DB is nil.
// Calling any real query method on it will panic inside gorm, which is the
// expected behaviour for "reached DB path" assertions (used with assert.Panics).
func newNilDB() *gorm.DB {
	return &gorm.DB{}
}

// ── Bypass-role tests (DB is never queried) ───────────────────────────────────

func TestRequireOwner_AdminBypasses(t *testing.T) {
	e := echo.New()

	cfg := OwnerConfig{
		ResourceTable: "products",
		OwnerField:    "created_by",
		BypassRoles:   []string{"admin"},
	}

	mw := RequireOwner(cfg, newNilDB())
	handler := mw(okHandler)

	c := newOwnerContext(e, 42, []string{"admin"}, "99")
	err := handler(c)

	// No DB call is made; next handler executes and returns 200.
	assert.NoError(t, err)
}

func TestRequireOwner_SuperuserBypasses(t *testing.T) {
	e := echo.New()

	cfg := OwnerConfig{
		ResourceTable: "products",
		OwnerField:    "created_by",
		BypassRoles:   []string{"admin", "superuser"},
	}

	mw := RequireOwner(cfg, newNilDB())
	handler := mw(okHandler)

	c := newOwnerContext(e, 42, []string{"superuser"}, "99")
	err := handler(c)

	assert.NoError(t, err)
}

func TestRequireOwner_MultipleRoles_BypassMatchesOne(t *testing.T) {
	e := echo.New()

	cfg := OwnerConfig{
		ResourceTable: "products",
		OwnerField:    "created_by",
		BypassRoles:   []string{"admin"},
	}

	mw := RequireOwner(cfg, newNilDB())
	handler := mw(okHandler)

	// User holds "user" and "admin"; "admin" triggers bypass before DB is reached.
	c := newOwnerContext(e, 42, []string{"user", "admin"}, "99")
	err := handler(c)

	assert.NoError(t, err)
}

// ── Tests that prove the bypass loop is skipped and the DB path is reached ────
//
// When no bypass role matches, the middleware advances to the gorm query.
// A nil-backed *gorm.DB panics inside gorm itself (nil pointer dereference in
// Statement.QuoteTo). We treat that panic as proof that the code reached the DB
// path, i.e. the bypass did NOT fire.

func TestRequireOwner_NoMatchingBypassRole_ReachesDBPath(t *testing.T) {
	e := echo.New()

	cfg := OwnerConfig{
		ResourceTable: "products",
		OwnerField:    "created_by",
		BypassRoles:   []string{"admin"},
	}

	mw := RequireOwner(cfg, newNilDB())
	handler := mw(okHandler)

	c := newOwnerContext(e, 42, []string{"user", "moderator"}, "99")

	// "user" and "moderator" are not bypass roles → middleware proceeds to DB query
	// → nil DB panics, confirming the bypass was not triggered.
	assert.Panics(t, func() {
		_ = handler(c) //nolint:errcheck
	})
}

func TestRequireOwner_EmptyRoles_ReachesDBPath(t *testing.T) {
	e := echo.New()

	cfg := OwnerConfig{
		ResourceTable: "products",
		OwnerField:    "created_by",
		BypassRoles:   []string{"admin"},
	}

	mw := RequireOwner(cfg, newNilDB())
	handler := mw(okHandler)

	// No roles at all → bypass loop body never runs → middleware proceeds to DB.
	c := newOwnerContext(e, 42, []string{}, "99")
	assert.Panics(t, func() {
		_ = handler(c) //nolint:errcheck
	})
}

func TestRequireOwner_EmptyBypassRoles_ReachesDBPath(t *testing.T) {
	e := echo.New()

	cfg := OwnerConfig{
		ResourceTable: "products",
		OwnerField:    "created_by",
		BypassRoles:   []string{}, // bypass list is empty; nothing can ever bypass
	}

	mw := RequireOwner(cfg, newNilDB())
	handler := mw(okHandler)

	// Even "admin" won't bypass when BypassRoles is empty.
	c := newOwnerContext(e, 42, []string{"admin"}, "99")
	assert.Panics(t, func() {
		_ = handler(c) //nolint:errcheck
	})
}

// ── Notes on DB-dependent tests ───────────────────────────────────────────────
//
// The following scenarios require a real database connection and are covered by
// the integration test suite (step 9 of the RBAC implementation plan), which
// spins up a MySQL container via testcontainers:
//
//   - DB query succeeds and ownerID == userID  → next handler called (200)
//   - DB query succeeds and ownerID != userID  → 403 "not the owner of this resource"
//   - DB query fails / row not found           → 404 "resource not found"
//
// None of these are exercised here because:
//   (a) gorm.Open with a bad DSN returns a *gorm.DB whose Statement is nil,
//       and subsequent calls to db.Table() panic inside gorm before returning
//       any error — so there is no safe way to drive the 404/403 paths from a
//       unit test without an actual driver + connection or a SQL mock library.
//   (b) Adding github.com/DATA-DOG/go-sqlmock would be a new go.mod dependency
//       outside the scope of this unit-test file.
//
// All three branches above ARE exercised in the integration tests.
