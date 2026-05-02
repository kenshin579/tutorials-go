package repository

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/tutorials-go/wiki-permissions/2-rbac/backend/config"
	"github.com/kenshin579/tutorials-go/wiki-permissions/2-rbac/backend/domain"
)

func TestUserRepository_CreateAndFindByEmail(t *testing.T) {
	db, err := config.OpenDB(":memory:")
	require.NoError(t, err)
	repo := NewUserRepository(db)

	u := &domain.User{Email: "alice@example.com", Name: "Alice", PasswordHash: "x"}
	require.NoError(t, repo.Create(u))
	assert.NotZero(t, u.ID)

	found, err := repo.FindByEmail("alice@example.com")
	require.NoError(t, err)
	assert.Equal(t, u.ID, found.ID)
}

func TestUserRepository_FindByID_NotFound(t *testing.T) {
	db, err := config.OpenDB(":memory:")
	require.NoError(t, err)
	repo := NewUserRepository(db)

	_, err = repo.FindByID(999)
	var nf domain.ErrNotFound
	require.True(t, errors.As(err, &nf))
	assert.Equal(t, "user", nf.Resource)
}

func TestUserRepository_AssignAndRevokeRole(t *testing.T) {
	db, err := config.OpenDB(":memory:")
	require.NoError(t, err)
	users := NewUserRepository(db)

	u := &domain.User{Email: "a@x", Name: "A", PasswordHash: "x"}
	require.NoError(t, users.Create(u))

	role := &domain.Role{Name: "admin"}
	require.NoError(t, db.Create(role).Error)

	require.NoError(t, users.AssignRole(u.ID, role.ID))

	got, err := users.FindByID(u.ID)
	require.NoError(t, err)
	require.Len(t, got.Roles, 1)
	assert.Equal(t, "admin", got.Roles[0].Name)

	require.NoError(t, users.RevokeRole(u.ID, role.ID))
	got2, err := users.FindByID(u.ID)
	require.NoError(t, err)
	assert.Empty(t, got2.Roles)
}
