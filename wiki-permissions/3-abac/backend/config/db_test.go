package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenDB_AutoMigratesAllTables(t *testing.T) {
	db, err := OpenDB(":memory:")
	require.NoError(t, err)

	for _, table := range []string{"departments", "users", "pages"} {
		assert.True(t, db.Migrator().HasTable(table), "table %s should exist", table)
	}
}
