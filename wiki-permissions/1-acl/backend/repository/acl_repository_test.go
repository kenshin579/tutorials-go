package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/tutorials-go/wiki-permissions/1-acl/backend/config"
	"github.com/kenshin579/tutorials-go/wiki-permissions/1-acl/backend/domain"
)

func TestACLRepository_GrantAndFind(t *testing.T) {
	db, err := config.OpenDB(":memory:")
	require.NoError(t, err)

	repo := NewACLRepository(db)
	require.NoError(t, repo.Grant(1, 2, domain.ActionRead))
	require.NoError(t, repo.Grant(1, 2, domain.ActionEdit))

	entries, err := repo.FindByPageAndUser(1, 2)
	require.NoError(t, err)
	assert.Len(t, entries, 2)
}

func TestACLRepository_Grant_Idempotent(t *testing.T) {
	db, _ := config.OpenDB(":memory:")
	repo := NewACLRepository(db)
	require.NoError(t, repo.Grant(1, 2, domain.ActionRead))
	require.NoError(t, repo.Grant(1, 2, domain.ActionRead)) // duplicate must not error
	entries, _ := repo.FindByPageAndUser(1, 2)
	assert.Len(t, entries, 1)
}

func TestACLRepository_Revoke_RemovesAllActions(t *testing.T) {
	db, _ := config.OpenDB(":memory:")
	repo := NewACLRepository(db)
	repo.Grant(1, 2, domain.ActionRead)
	repo.Grant(1, 2, domain.ActionEdit)

	require.NoError(t, repo.Revoke(1, 2))
	entries, _ := repo.FindByPageAndUser(1, 2)
	assert.Empty(t, entries)
}
