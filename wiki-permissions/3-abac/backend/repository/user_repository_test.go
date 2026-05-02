package repository

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/tutorials-go/wiki-permissions/3-abac/backend/config"
	"github.com/kenshin579/tutorials-go/wiki-permissions/3-abac/backend/domain"
)

func TestUserRepository_CreateAndFindByEmail_PreloadsDepartment(t *testing.T) {
	db, err := config.OpenDB(":memory:")
	require.NoError(t, err)
	users := NewUserRepository(db)

	dept := &domain.Department{Name: "Engineering"}
	require.NoError(t, db.Create(dept).Error)

	u := &domain.User{
		Email:          "alice@example.com",
		Name:           "Alice",
		PasswordHash:   "x",
		DepartmentID:   dept.ID,
		EmploymentType: domain.EmploymentFulltime,
	}
	require.NoError(t, users.Create(u))

	got, err := users.FindByEmail("alice@example.com")
	require.NoError(t, err)
	require.NotNil(t, got.Department)
	assert.Equal(t, "Engineering", got.Department.Name)
	assert.Equal(t, domain.EmploymentFulltime, got.EmploymentType)
}

func TestUserRepository_FindByID_NotFound(t *testing.T) {
	db, err := config.OpenDB(":memory:")
	require.NoError(t, err)
	users := NewUserRepository(db)

	_, err = users.FindByID(999)
	var nf domain.ErrNotFound
	require.True(t, errors.As(err, &nf))
	assert.Equal(t, "user", nf.Resource)
}
