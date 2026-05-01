package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/tutorials-go/wiki-permissions/1-acl/backend/domain"
)

func TestSeed_PopulatesUsersPagesACLs(t *testing.T) {
	db, err := OpenDB(":memory:")
	require.NoError(t, err)
	require.NoError(t, Seed(db))

	var userCount, pageCount, aclCount int64
	db.Model(&domain.User{}).Count(&userCount)
	db.Model(&domain.Page{}).Count(&pageCount)
	db.Model(&domain.ACLEntry{}).Count(&aclCount)

	assert.Equal(t, int64(4), userCount)
	assert.Equal(t, int64(3), pageCount)
	assert.Equal(t, int64(7), aclCount) // EngRoadmap 2 + Q4Marketing 2 + OnboardingGuide 3
}

func TestSeed_Idempotent(t *testing.T) {
	db, _ := OpenDB(":memory:")
	require.NoError(t, Seed(db))
	require.NoError(t, Seed(db)) // 두 번째 호출도 에러 없어야 함

	var userCount int64
	db.Model(&domain.User{}).Count(&userCount)
	assert.Equal(t, int64(4), userCount)
}
