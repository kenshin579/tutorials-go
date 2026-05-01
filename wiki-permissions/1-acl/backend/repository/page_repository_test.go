package repository

import (
	"errors"
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
	require.Len(t, got, 2) // p1 (owner) + p2 (shared)
	// 결과 페이지 식별 + Order("pages.id ASC") 보장
	assert.Equal(t, p1.ID, got[0].ID)
	assert.Equal(t, p2.ID, got[1].ID)
}

// 같은 페이지에 read+edit 두 ACL이 부여돼도 Distinct로 한 번만 노출돼야 한다.
func TestPageRepository_ListAccessibleBy_DistinctOnMultipleACLs(t *testing.T) {
	db, err := config.OpenDB(":memory:")
	require.NoError(t, err)

	users := NewUserRepository(db)
	pages := NewPageRepository(db)
	acls := NewACLRepository(db)

	alice := &domain.User{Email: "a@x", Name: "Alice", PasswordHash: "x"}
	bob := &domain.User{Email: "b@x", Name: "Bob", PasswordHash: "x"}
	require.NoError(t, users.Create(alice))
	require.NoError(t, users.Create(bob))

	page := &domain.Page{Title: "Bob's page", OwnerID: bob.ID}
	require.NoError(t, pages.Create(page))

	require.NoError(t, acls.Grant(page.ID, alice.ID, domain.ActionRead))
	require.NoError(t, acls.Grant(page.ID, alice.ID, domain.ActionEdit))

	got, err := pages.ListAccessibleBy(alice.ID)
	require.NoError(t, err)
	assert.Len(t, got, 1)
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
