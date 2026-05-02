package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/tutorials-go/wiki-permissions/3-abac/backend/domain"
)

func TestSeed_PopulatesAll(t *testing.T) {
	db, err := OpenDB(":memory:")
	require.NoError(t, err)
	require.NoError(t, Seed(db))

	var deptCount, userCount, pageCount int64
	db.Model(&domain.Department{}).Count(&deptCount)
	db.Model(&domain.User{}).Count(&userCount)
	db.Model(&domain.Page{}).Count(&pageCount)

	assert.Equal(t, int64(2), deptCount)
	assert.Equal(t, int64(4), userCount)
	assert.Equal(t, int64(3), pageCount)
}

func TestSeed_Idempotent(t *testing.T) {
	db, err := OpenDB(":memory:")
	require.NoError(t, err)
	require.NoError(t, Seed(db))
	require.NoError(t, Seed(db))

	var deptCount, userCount, pageCount int64
	db.Model(&domain.Department{}).Count(&deptCount)
	db.Model(&domain.User{}).Count(&userCount)
	db.Model(&domain.Page{}).Count(&pageCount)

	assert.Equal(t, int64(2), deptCount)
	assert.Equal(t, int64(4), userCount)
	assert.Equal(t, int64(3), pageCount)
}

// 시드 사용자의 속성이 의도대로 저장되는지 확인.
func TestSeed_UserAttributes(t *testing.T) {
	db, err := OpenDB(":memory:")
	require.NoError(t, err)
	require.NoError(t, Seed(db))

	var dave domain.User
	require.NoError(t, db.Preload("Department").Where("email = ?", "dave@example.com").First(&dave).Error)
	assert.Equal(t, "Marketing", dave.Department.Name)
	assert.Equal(t, domain.EmploymentContract, dave.EmploymentType)

	var alice domain.User
	require.NoError(t, db.Preload("Department").Where("email = ?", "alice@example.com").First(&alice).Error)
	assert.Equal(t, "Engineering", alice.Department.Name)
	assert.Equal(t, domain.EmploymentFulltime, alice.EmploymentType)
}
