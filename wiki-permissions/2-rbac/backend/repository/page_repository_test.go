package repository

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/tutorials-go/wiki-permissions/2-rbac/backend/config"
	"github.com/kenshin579/tutorials-go/wiki-permissions/2-rbac/backend/domain"
)

func TestPageRepository_List_ReturnsAllPages(t *testing.T) {
	db, err := config.OpenDB(":memory:")
	require.NoError(t, err)

	users := NewUserRepository(db)
	pages := NewPageRepository(db)

	alice := &domain.User{Email: "a@x", Name: "Alice", PasswordHash: "x"}
	require.NoError(t, users.Create(alice))

	require.NoError(t, pages.Create(&domain.Page{Title: "P1", OwnerID: alice.ID}))
	require.NoError(t, pages.Create(&domain.Page{Title: "P2", OwnerID: alice.ID}))

	got, err := pages.List()
	require.NoError(t, err)
	assert.Len(t, got, 2)
}

func TestPageRepository_FindByID_NotFound(t *testing.T) {
	db, err := config.OpenDB(":memory:")
	require.NoError(t, err)
	repo := NewPageRepository(db)

	_, err = repo.FindByID(999)
	var nf domain.ErrNotFound
	require.True(t, errors.As(err, &nf))
	assert.Equal(t, "page", nf.Resource)
}

func TestPageRepository_Delete(t *testing.T) {
	db, err := config.OpenDB(":memory:")
	require.NoError(t, err)
	users := NewUserRepository(db)
	pages := NewPageRepository(db)

	alice := &domain.User{Email: "a@x", Name: "Alice", PasswordHash: "x"}
	require.NoError(t, users.Create(alice))

	p := &domain.Page{Title: "tmp", OwnerID: alice.ID}
	require.NoError(t, pages.Create(p))

	require.NoError(t, pages.Delete(p.ID))

	_, err = pages.FindByID(p.ID)
	var nf domain.ErrNotFound
	require.True(t, errors.As(err, &nf))
}
