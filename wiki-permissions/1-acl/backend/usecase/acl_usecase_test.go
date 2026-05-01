package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/tutorials-go/wiki-permissions/1-acl/backend/config"
	"github.com/kenshin579/tutorials-go/wiki-permissions/1-acl/backend/domain"
	"github.com/kenshin579/tutorials-go/wiki-permissions/1-acl/backend/repository"
)

func TestACLUsecase_Grant_OwnerOnly(t *testing.T) {
	db, _ := config.OpenDB(":memory:")
	config.Seed(db)
	uc := NewACLUsecase(repository.NewPageRepository(db), repository.NewACLRepository(db))

	users := repository.NewUserRepository(db)
	alice, _ := users.FindByEmail("alice@example.com")
	bob, _ := users.FindByEmail("bob@example.com")
	carol, _ := users.FindByEmail("carol@example.com")

	var p domain.Page
	require.NoError(t, db.Where("title = ?", "Engineering Roadmap").First(&p).Error)

	// Alice가 owner — OK
	require.NoError(t, uc.Grant(p.ID, alice.ID, carol.ID, domain.ActionEdit))

	// Bob은 owner 아님 — Forbidden
	err := uc.Grant(p.ID, bob.ID, carol.ID, domain.ActionEdit)
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestACLUsecase_Revoke_OwnerOnly(t *testing.T) {
	db, _ := config.OpenDB(":memory:")
	config.Seed(db)
	uc := NewACLUsecase(repository.NewPageRepository(db), repository.NewACLRepository(db))

	users := repository.NewUserRepository(db)
	alice, _ := users.FindByEmail("alice@example.com")
	bob, _ := users.FindByEmail("bob@example.com")

	var p domain.Page
	db.Where("title = ?", "Engineering Roadmap").First(&p)

	require.NoError(t, uc.Revoke(p.ID, alice.ID, bob.ID))

	// Bob (non-owner) tries
	err := uc.Revoke(p.ID, bob.ID, alice.ID)
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestACLUsecase_List_OwnerOnly(t *testing.T) {
	db, _ := config.OpenDB(":memory:")
	config.Seed(db)
	uc := NewACLUsecase(repository.NewPageRepository(db), repository.NewACLRepository(db))

	users := repository.NewUserRepository(db)
	alice, _ := users.FindByEmail("alice@example.com")
	dave, _ := users.FindByEmail("dave@example.com")

	var p domain.Page
	db.Where("title = ?", "Engineering Roadmap").First(&p)

	entries, err := uc.List(p.ID, alice.ID)
	require.NoError(t, err)
	assert.NotEmpty(t, entries)

	_, err = uc.List(p.ID, dave.ID)
	assert.ErrorIs(t, err, ErrForbidden)
}
