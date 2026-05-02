package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/tutorials-go/wiki-permissions/2-rbac/backend/config"
	"github.com/kenshin579/tutorials-go/wiki-permissions/2-rbac/backend/domain"
	"github.com/kenshin579/tutorials-go/wiki-permissions/2-rbac/backend/repository"
)

func setupPageEnv(t *testing.T) (*PageUsecase, map[string]uint) {
	t.Helper()
	db, err := config.OpenDB(":memory:")
	require.NoError(t, err)
	require.NoError(t, config.Seed(db))

	pages := repository.NewPageRepository(db)
	perms := repository.NewPermissionRepository(db)
	uc := NewPageUsecase(pages, perms)

	users := repository.NewUserRepository(db)
	emailToID := map[string]uint{}
	for _, e := range []string{"alice@example.com", "bob@example.com", "carol@example.com", "dave@example.com"} {
		u, err := users.FindByEmail(e)
		require.NoError(t, err)
		emailToID[e] = u.ID
	}
	return uc, emailToID
}

func TestHasPermission_Lookup(t *testing.T) {
	perms := []domain.Permission{
		{Resource: "pages", Action: "read"},
		{Resource: "pages", Action: "edit"},
	}
	assert.True(t, HasPermission(perms, "pages:read"))
	assert.True(t, HasPermission(perms, "pages:edit"))
	assert.False(t, HasPermission(perms, "pages:delete"))
	assert.False(t, HasPermission(nil, "anything"))
}

func TestPageUsecase_List_AnyAuthenticatedReadsAll(t *testing.T) {
	uc, ids := setupPageEnv(t)
	for _, email := range []string{"alice@example.com", "bob@example.com", "carol@example.com", "dave@example.com"} {
		got, err := uc.List(ids[email])
		require.NoError(t, err, "user %s should have pages:read", email)
		assert.Len(t, got, 3, "user %s should see all 3 pages", email)
	}
}

func TestPageUsecase_Update_RequiresEdit(t *testing.T) {
	uc, ids := setupPageEnv(t)
	// bob (editor) — OK
	p, err := uc.Update(1, ids["bob@example.com"], "new title", "x")
	require.NoError(t, err)
	assert.Equal(t, "new title", p.Title)

	// carol (viewer) — Forbidden
	_, err = uc.Update(1, ids["carol@example.com"], "x", "y")
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestPageUsecase_Delete_OnlyAdmin(t *testing.T) {
	uc, ids := setupPageEnv(t)
	// bob (editor) — Forbidden
	err := uc.Delete(1, ids["bob@example.com"])
	assert.ErrorIs(t, err, ErrForbidden)

	// alice (admin) — OK
	require.NoError(t, uc.Delete(1, ids["alice@example.com"]))
}

func TestPageUsecase_Create_RequiresCreate(t *testing.T) {
	uc, ids := setupPageEnv(t)
	// alice (admin), bob (editor) — OK
	for _, email := range []string{"alice@example.com", "bob@example.com"} {
		p, err := uc.Create(ids[email], "new page", "content")
		require.NoError(t, err, "user %s should be able to create", email)
		assert.NotZero(t, p.ID)
	}
	// carol (viewer) — Forbidden
	_, err := uc.Create(ids["carol@example.com"], "x", "y")
	assert.ErrorIs(t, err, ErrForbidden)
}
