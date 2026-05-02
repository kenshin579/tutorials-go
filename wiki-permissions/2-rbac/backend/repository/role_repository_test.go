package repository

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/tutorials-go/wiki-permissions/2-rbac/backend/config"
	"github.com/kenshin579/tutorials-go/wiki-permissions/2-rbac/backend/domain"
)

func TestRoleRepository_FindByName_PreloadsPermissions(t *testing.T) {
	db, err := config.OpenDB(":memory:")
	require.NoError(t, err)
	repo := NewRoleRepository(db)

	role := &domain.Role{Name: "admin"}
	require.NoError(t, db.Create(role).Error)
	perm := &domain.Permission{Resource: "pages", Action: "read"}
	require.NoError(t, db.Create(perm).Error)
	require.NoError(t, db.Model(role).Association("Permissions").Append(perm))

	got, err := repo.FindByName("admin")
	require.NoError(t, err)
	require.Len(t, got.Permissions, 1)
	assert.Equal(t, "pages:read", got.Permissions[0].Key())
}

func TestRoleRepository_FindByID_NotFound(t *testing.T) {
	db, err := config.OpenDB(":memory:")
	require.NoError(t, err)
	repo := NewRoleRepository(db)

	_, err = repo.FindByID(999)
	var nf domain.ErrNotFound
	require.True(t, errors.As(err, &nf))
	assert.Equal(t, "role", nf.Resource)
}
