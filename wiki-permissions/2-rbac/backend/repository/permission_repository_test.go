package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/tutorials-go/wiki-permissions/2-rbac/backend/config"
	"github.com/kenshin579/tutorials-go/wiki-permissions/2-rbac/backend/domain"
)

// 두 role(admin: read+edit, viewer: read)을 모두 가진 사용자에 대해
// pages:read 권한이 한 번만 노출되는지 검증한다 (Distinct).
func TestPermissionRepository_FindByUserID_DistinctOnMultipleRoles(t *testing.T) {
	db, err := config.OpenDB(":memory:")
	require.NoError(t, err)

	users := NewUserRepository(db)
	repo := NewPermissionRepository(db)

	u := &domain.User{Email: "a@x", Name: "A", PasswordHash: "x"}
	require.NoError(t, users.Create(u))

	pRead := &domain.Permission{Resource: "pages", Action: "read"}
	pEdit := &domain.Permission{Resource: "pages", Action: "edit"}
	require.NoError(t, db.Create(pRead).Error)
	require.NoError(t, db.Create(pEdit).Error)

	admin := &domain.Role{Name: "admin"}
	viewer := &domain.Role{Name: "viewer"}
	require.NoError(t, db.Create(admin).Error)
	require.NoError(t, db.Create(viewer).Error)
	require.NoError(t, db.Model(admin).Association("Permissions").Append(pRead, pEdit))
	require.NoError(t, db.Model(viewer).Association("Permissions").Append(pRead))

	require.NoError(t, users.AssignRole(u.ID, admin.ID))
	require.NoError(t, users.AssignRole(u.ID, viewer.ID))

	got, err := repo.FindByUserID(u.ID)
	require.NoError(t, err)
	require.Len(t, got, 2)
	assert.Equal(t, "pages:edit", got[0].Key()) // 알파벳순: edit < read
	assert.Equal(t, "pages:read", got[1].Key())
}

func TestPermissionRepository_FindByUserID_NoRoles(t *testing.T) {
	db, err := config.OpenDB(":memory:")
	require.NoError(t, err)

	users := NewUserRepository(db)
	repo := NewPermissionRepository(db)

	u := &domain.User{Email: "a@x", Name: "A", PasswordHash: "x"}
	require.NoError(t, users.Create(u))

	got, err := repo.FindByUserID(u.ID)
	require.NoError(t, err)
	assert.Empty(t, got)
}
