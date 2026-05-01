package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/tutorials-go/wiki-permissions/1-acl/backend/config"
	"github.com/kenshin579/tutorials-go/wiki-permissions/1-acl/backend/domain"
)

func TestPageRepository_ListAccessibleBy(t *testing.T) {
	db, err := config.OpenDB(":memory:")
	require.NoError(t, err)

	users := NewUserRepository(db)
	pages := NewPageRepository(db)
	acls := NewACLRepository(db)

	alice := &domain.User{Email: "a@x", Name: "Alice", PasswordHash: "x"}
	bob := &domain.User{Email: "b@x", Name: "Bob", PasswordHash: "x"}
	require.NoError(t, users.Create(alice))
	require.NoError(t, users.Create(bob))

	p1 := &domain.Page{Title: "Owned by Alice", OwnerID: alice.ID}
	p2 := &domain.Page{Title: "Owned by Bob, shared with Alice", OwnerID: bob.ID}
	p3 := &domain.Page{Title: "Owned by Bob, no share", OwnerID: bob.ID}
	require.NoError(t, pages.Create(p1))
	require.NoError(t, pages.Create(p2))
	require.NoError(t, pages.Create(p3))

	require.NoError(t, acls.Grant(p2.ID, alice.ID, domain.ActionRead))

	got, err := pages.ListAccessibleBy(alice.ID)
	require.NoError(t, err)
	assert.Len(t, got, 2) // p1 (owner) + p2 (shared)
}
