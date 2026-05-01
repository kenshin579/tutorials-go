package usecase

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/tutorials-go/wiki-permissions/1-acl/backend/config"
	"github.com/kenshin579/tutorials-go/wiki-permissions/1-acl/backend/domain"
	"github.com/kenshin579/tutorials-go/wiki-permissions/1-acl/backend/pkg/passwordhash"
	"github.com/kenshin579/tutorials-go/wiki-permissions/1-acl/backend/repository"
)

func TestAuthUsecase_Login_Success(t *testing.T) {
	db, _ := config.OpenDB(":memory:")
	users := repository.NewUserRepository(db)
	hash, _ := passwordhash.Hash("password")
	require.NoError(t, users.Create(&domain.User{Email: "a@x", Name: "A", PasswordHash: hash}))

	uc := NewAuthUsecase(users, "secret", time.Hour)
	token, _, err := uc.Login("a@x", "password")
	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestAuthUsecase_Login_WrongPassword(t *testing.T) {
	db, _ := config.OpenDB(":memory:")
	users := repository.NewUserRepository(db)
	hash, _ := passwordhash.Hash("password")
	users.Create(&domain.User{Email: "a@x", Name: "A", PasswordHash: hash})

	uc := NewAuthUsecase(users, "secret", time.Hour)
	_, _, err := uc.Login("a@x", "wrong")
	assert.ErrorIs(t, err, ErrInvalidCredentials)
}

func TestAuthUsecase_Login_UserNotFound(t *testing.T) {
	db, _ := config.OpenDB(":memory:")
	users := repository.NewUserRepository(db)
	uc := NewAuthUsecase(users, "secret", time.Hour)

	_, _, err := uc.Login("nope@x", "password")
	assert.ErrorIs(t, err, ErrInvalidCredentials)
}
