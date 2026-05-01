package repository

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/tutorials-go/wiki-permissions/1-acl/backend/config"
	"github.com/kenshin579/tutorials-go/wiki-permissions/1-acl/backend/domain"
)

func TestUserRepository_CreateAndFindByEmail(t *testing.T) {
	db, err := config.OpenDB(":memory:")
	require.NoError(t, err)
	repo := NewUserRepository(db)

	u := &domain.User{Email: "alice@example.com", Name: "Alice", PasswordHash: "x"}
	require.NoError(t, repo.Create(u))
	assert.NotZero(t, u.ID)

	found, err := repo.FindByEmail("alice@example.com")
	require.NoError(t, err)
	assert.Equal(t, u.ID, found.ID)
}

func TestUserRepository_FindByID_NotFound(t *testing.T) {
	db, err := config.OpenDB(":memory:")
	require.NoError(t, err)
	repo := NewUserRepository(db)

	_, err = repo.FindByID(999)
	var nf domain.ErrNotFound
	require.True(t, errors.As(err, &nf))
	assert.Equal(t, "user", nf.Resource)
}

func TestUserRepository_FindByEmail_NotFound(t *testing.T) {
	db, err := config.OpenDB(":memory:")
	require.NoError(t, err)
	repo := NewUserRepository(db)

	_, err = repo.FindByEmail("missing@example.com")
	var nf domain.ErrNotFound
	require.True(t, errors.As(err, &nf))
	assert.Equal(t, "user", nf.Resource)
}
