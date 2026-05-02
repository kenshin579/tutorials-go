package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/tutorials-go/wiki-permissions/2-rbac/backend/domain"
)

func TestSeed_PopulatesAll(t *testing.T) {
	db, err := OpenDB(":memory:")
	require.NoError(t, err)
	require.NoError(t, Seed(db))

	var userCount, pageCount, permCount, roleCount, userRoleCount, rolePermCount int64
	db.Model(&domain.User{}).Count(&userCount)
	db.Model(&domain.Page{}).Count(&pageCount)
	db.Model(&domain.Permission{}).Count(&permCount)
	db.Model(&domain.Role{}).Count(&roleCount)
	db.Table("user_roles").Count(&userRoleCount)
	db.Table("role_permissions").Count(&rolePermCount)

	assert.Equal(t, int64(4), userCount)
	assert.Equal(t, int64(3), pageCount)
	assert.Equal(t, int64(6), permCount)
	assert.Equal(t, int64(3), roleCount)
	assert.Equal(t, int64(4), userRoleCount)
	// admin 6 + editor 3 + viewer 1 = 10
	assert.Equal(t, int64(10), rolePermCount)
}

func TestSeed_Idempotent(t *testing.T) {
	db, err := OpenDB(":memory:")
	require.NoError(t, err)
	require.NoError(t, Seed(db))
	require.NoError(t, Seed(db))

	var userCount, pageCount, permCount, roleCount, userRoleCount, rolePermCount int64
	db.Model(&domain.User{}).Count(&userCount)
	db.Model(&domain.Page{}).Count(&pageCount)
	db.Model(&domain.Permission{}).Count(&permCount)
	db.Model(&domain.Role{}).Count(&roleCount)
	db.Table("user_roles").Count(&userRoleCount)
	db.Table("role_permissions").Count(&rolePermCount)

	assert.Equal(t, int64(4), userCount)
	assert.Equal(t, int64(3), pageCount)
	assert.Equal(t, int64(6), permCount)
	assert.Equal(t, int64(3), roleCount)
	assert.Equal(t, int64(4), userRoleCount)
	assert.Equal(t, int64(10), rolePermCount)
}
