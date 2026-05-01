package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/tutorials-go/wiki-permissions/1-acl/backend/domain"
)

func TestOpenDB_AutoMigratesAllTables(t *testing.T) {
	db, err := OpenDB(":memory:")
	require.NoError(t, err)

	tables := []string{"users", "pages", "acl_entries"}
	for _, table := range tables {
		assert.True(t, db.Migrator().HasTable(table), "table %s should exist", table)
	}

	// FK indexes
	assert.True(t, db.Migrator().HasIndex(&domain.Page{}, "owner_id"))
	assert.True(t, db.Migrator().HasIndex(&domain.ACLEntry{}, "idx_page_user_action"))

	sqlDB, err := db.DB()
	require.NoError(t, err)
	assert.NoError(t, sqlDB.Close())
}
