package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/tutorials-go/wiki-permissions/3-abac/backend/config"
	"github.com/kenshin579/tutorials-go/wiki-permissions/3-abac/backend/domain"
)

func TestPageRepository_List_ReturnsAllPagesWithDepartment(t *testing.T) {
	db, err := config.OpenDB(":memory:")
	require.NoError(t, err)
	users := NewUserRepository(db)
	pages := NewPageRepository(db)

	dept := &domain.Department{Name: "Engineering"}
	require.NoError(t, db.Create(dept).Error)

	alice := &domain.User{
		Email: "a@x", Name: "A", PasswordHash: "x",
		DepartmentID: dept.ID, EmploymentType: domain.EmploymentFulltime,
	}
	require.NoError(t, users.Create(alice))

	deptID := dept.ID
	require.NoError(t, pages.Create(&domain.Page{
		Title: "Internal", OwnerID: alice.ID,
		Confidentiality: domain.ConfidentialityInternal, DepartmentID: &deptID,
	}))
	require.NoError(t, pages.Create(&domain.Page{
		Title: "Public", OwnerID: alice.ID,
		Confidentiality: domain.ConfidentialityPublic, DepartmentID: nil,
	}))

	got, err := pages.List()
	require.NoError(t, err)
	require.Len(t, got, 2)
	// internal 페이지에는 Department가 preload되어야 함
	assert.NotNil(t, got[0].Department)
	assert.Equal(t, "Engineering", got[0].Department.Name)
	// public 페이지에는 Department가 nil이어야 함
	assert.Nil(t, got[1].Department)
}
