package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenDB_AutoMigratesAllTables(t *testing.T) {
	db, err := OpenDB(":memory:")
	require.NoError(t, err)

	tables := []string{"users", "pages", "roles", "permissions", "user_roles", "role_permissions"}
	for _, table := range tables {
		assert.True(t, db.Migrator().HasTable(table), "table %s should exist", table)
	}
}
