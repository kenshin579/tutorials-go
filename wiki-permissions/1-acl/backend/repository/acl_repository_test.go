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

// ListByPage는 user_id, action 순으로 정렬해 owner UI에서 그룹화하기 쉽게 한다.
// 'edit' < 'read' (알파벳순)이므로 한 사용자가 둘 다 가지면 edit이 먼저 온다.
func TestACLRepository_ListByPage_OrdersByUserThenAction(t *testing.T) {
	db, err := config.OpenDB(":memory:")
	require.NoError(t, err)
	repo := NewACLRepository(db)

	require.NoError(t, repo.Grant(1, 3, domain.ActionRead))
	require.NoError(t, repo.Grant(1, 2, domain.ActionRead))
	require.NoError(t, repo.Grant(1, 2, domain.ActionEdit))

	entries, err := repo.ListByPage(1)
	require.NoError(t, err)
	require.Len(t, entries, 3)

	assert.Equal(t, uint(2), entries[0].UserID)
	assert.Equal(t, domain.ActionEdit, entries[0].Action)
	assert.Equal(t, uint(2), entries[1].UserID)
	assert.Equal(t, domain.ActionRead, entries[1].Action)
	assert.Equal(t, uint(3), entries[2].UserID)
}
