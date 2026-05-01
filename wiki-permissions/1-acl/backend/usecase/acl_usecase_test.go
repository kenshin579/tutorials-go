package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/tutorials-go/wiki-permissions/1-acl/backend/config"
	"github.com/kenshin579/tutorials-go/wiki-permissions/1-acl/backend/domain"
	"github.com/kenshin579/tutorials-go/wiki-permissions/1-acl/backend/repository"
)

func setupACLEnv(t *testing.T) (uc *ACLUsecase, alice, bob, carol, dave *domain.User, page *domain.Page) {
	t.Helper()
	db, err := config.OpenDB(":memory:")
	require.NoError(t, err)
	require.NoError(t, config.Seed(db))

	uc = NewACLUsecase(repository.NewPageRepository(db), repository.NewACLRepository(db))
	users := repository.NewUserRepository(db)

	alice, err = users.FindByEmail("alice@example.com")
	require.NoError(t, err)
	bob, err = users.FindByEmail("bob@example.com")
	require.NoError(t, err)
	carol, err = users.FindByEmail("carol@example.com")
	require.NoError(t, err)
	dave, err = users.FindByEmail("dave@example.com")
	require.NoError(t, err)

	page = &domain.Page{}
	require.NoError(t, db.Where("title = ?", "Engineering Roadmap").First(page).Error)
	return
}

func TestACLUsecase_Grant_OwnerOnly(t *testing.T) {
	uc, alice, bob, carol, _, p := setupACLEnv(t)

	// Alice가 owner — OK
	require.NoError(t, uc.Grant(p.ID, alice.ID, carol.ID, domain.ActionEdit))

	// 부여가 실제로 영속됐는지 List로 재확인
	entries, err := uc.List(p.ID, alice.ID)
	require.NoError(t, err)
	var found bool
	for _, e := range entries {
		if e.UserID == carol.ID && e.Action == domain.ActionEdit {
			found = true
			break
		}
	}
	assert.True(t, found, "carol should now have edit on the page")

	// Bob은 owner 아님 — Forbidden
	err = uc.Grant(p.ID, bob.ID, carol.ID, domain.ActionEdit)
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestACLUsecase_Grant_InvalidAction(t *testing.T) {
	uc, alice, _, carol, _, p := setupACLEnv(t)
	err := uc.Grant(p.ID, alice.ID, carol.ID, domain.Action("delete"))
	assert.ErrorIs(t, err, ErrInvalidAction)
}

func TestACLUsecase_Grant_PageNotFound(t *testing.T) {
	uc, alice, _, carol, _, _ := setupACLEnv(t)
	err := uc.Grant(99999, alice.ID, carol.ID, domain.ActionRead)
	var nf domain.ErrNotFound
	require.ErrorAs(t, err, &nf)
	assert.Equal(t, "page", nf.Resource)
}

func TestACLUsecase_Revoke_OwnerOnly(t *testing.T) {
	uc, alice, bob, _, _, p := setupACLEnv(t)

	require.NoError(t, uc.Revoke(p.ID, alice.ID, bob.ID))

	// bob의 entry가 실제로 사라졌는지 List로 재확인
	entries, err := uc.List(p.ID, alice.ID)
	require.NoError(t, err)
	for _, e := range entries {
		assert.NotEqual(t, bob.ID, e.UserID, "bob's ACL should be gone")
	}

	// Bob (non-owner) tries
	err = uc.Revoke(p.ID, bob.ID, alice.ID)
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestACLUsecase_List_OwnerOnly(t *testing.T) {
	uc, alice, _, _, dave, p := setupACLEnv(t)

	entries, err := uc.List(p.ID, alice.ID)
	require.NoError(t, err)
	assert.NotEmpty(t, entries)

	_, err = uc.List(p.ID, dave.ID)
	assert.ErrorIs(t, err, ErrForbidden)
}
