package usecase

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/tutorials-go/wiki-permissions/3-abac/backend/config"
	"github.com/kenshin579/tutorials-go/wiki-permissions/3-abac/backend/domain"
	"github.com/kenshin579/tutorials-go/wiki-permissions/3-abac/backend/repository"
)

func newAuthUC(t *testing.T) *AuthUsecase {
	t.Helper()
	db, err := config.OpenDB(":memory:")
	require.NoError(t, err)
	require.NoError(t, config.Seed(db))
	return NewAuthUsecase(repository.NewUserRepository(db), "secret", time.Hour)
}

func TestAuthUsecase_Login_ReturnsUserWithAttributes(t *testing.T) {
	uc := newAuthUC(t)
	token, user, err := uc.Login("dave@example.com", "password")
	require.NoError(t, err)
	assert.NotEmpty(t, token)
	require.NotNil(t, user.Department)
	assert.Equal(t, "Marketing", user.Department.Name)
	assert.Equal(t, domain.EmploymentContract, user.EmploymentType)
}

func TestAuthUsecase_Login_WrongPassword(t *testing.T) {
	uc := newAuthUC(t)
	_, _, err := uc.Login("alice@example.com", "wrong")
	assert.ErrorIs(t, err, ErrInvalidCredentials)
}

func TestAuthUsecase_Login_UserNotFound(t *testing.T) {
	uc := newAuthUC(t)
	_, _, err := uc.Login("nobody@example.com", "password")
	assert.ErrorIs(t, err, ErrInvalidCredentials)
}
