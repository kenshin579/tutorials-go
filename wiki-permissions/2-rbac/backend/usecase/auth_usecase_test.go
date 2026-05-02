package usecase

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/tutorials-go/wiki-permissions/2-rbac/backend/config"
	"github.com/kenshin579/tutorials-go/wiki-permissions/2-rbac/backend/repository"
)

func newAuthUC(t *testing.T) *AuthUsecase {
	t.Helper()
	db, err := config.OpenDB(":memory:")
	require.NoError(t, err)
	require.NoError(t, config.Seed(db))
	users := repository.NewUserRepository(db)
	perms := repository.NewPermissionRepository(db)
	return NewAuthUsecase(users, perms, "secret", time.Hour)
}

func TestAuthUsecase_Login_AdminGetsAllPermissions(t *testing.T) {
	uc := newAuthUC(t)
	res, err := uc.Login("alice@example.com", "password")
	require.NoError(t, err)
	assert.NotEmpty(t, res.Token)
	assert.Equal(t, "alice@example.com", res.User.Email)
	assert.Len(t, res.Permissions, 6) // admin: 6 permissions (pages:* + users:*)
}

func TestAuthUsecase_Login_ViewerGetsOnePermission(t *testing.T) {
	uc := newAuthUC(t)
	res, err := uc.Login("carol@example.com", "password")
	require.NoError(t, err)
	require.Len(t, res.Permissions, 1)
	assert.Equal(t, "pages:read", res.Permissions[0].Key())
}

func TestAuthUsecase_Login_WrongPassword(t *testing.T) {
	uc := newAuthUC(t)
	_, err := uc.Login("alice@example.com", "wrong")
	assert.ErrorIs(t, err, ErrInvalidCredentials)
}

func TestAuthUsecase_Login_UserNotFound(t *testing.T) {
	uc := newAuthUC(t)
	_, err := uc.Login("nobody@example.com", "password")
	assert.ErrorIs(t, err, ErrInvalidCredentials)
}
