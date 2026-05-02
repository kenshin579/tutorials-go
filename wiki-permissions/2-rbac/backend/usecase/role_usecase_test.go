package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/tutorials-go/wiki-permissions/2-rbac/backend/config"
	"github.com/kenshin579/tutorials-go/wiki-permissions/2-rbac/backend/domain"
	"github.com/kenshin579/tutorials-go/wiki-permissions/2-rbac/backend/repository"
)

func setupRoleEnv(t *testing.T) (uc *RoleUsecase, alice, bob, dave *domain.User, viewerRoleID, editorRoleID uint) {
	t.Helper()
	db, err := config.OpenDB(":memory:")
	require.NoError(t, err)
	require.NoError(t, config.Seed(db))

	users := repository.NewUserRepository(db)
	roles := repository.NewRoleRepository(db)
	perms := repository.NewPermissionRepository(db)
	uc = NewRoleUsecase(users, roles, perms)

	alice, err = users.FindByEmail("alice@example.com")
	require.NoError(t, err)
	bob, err = users.FindByEmail("bob@example.com")
	require.NoError(t, err)
	dave, err = users.FindByEmail("dave@example.com")
	require.NoError(t, err)

	viewer, err := roles.FindByName("viewer")
	require.NoError(t, err)
	viewerRoleID = viewer.ID
	editor, err := roles.FindByName("editor")
	require.NoError(t, err)
	editorRoleID = editor.ID
	return
}

func TestRoleUsecase_ListUsers_AdminOnly(t *testing.T) {
	uc, alice, bob, _, _, _ := setupRoleEnv(t)

	users, err := uc.ListUsers(alice.ID)
	require.NoError(t, err)
	assert.Len(t, users, 4)

	_, err = uc.ListUsers(bob.ID)
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestRoleUsecase_AssignRole_AdminOK_OthersForbidden(t *testing.T) {
	uc, alice, bob, dave, _, editorRoleID := setupRoleEnv(t)

	// admin (alice)이 dave에게 editor role 부여 → OK
	require.NoError(t, uc.AssignRole(alice.ID, dave.ID, editorRoleID))

	// 부여가 영속됐는지 ListUsers로 재확인
	users, err := uc.ListUsers(alice.ID)
	require.NoError(t, err)
	var daveRoles []string
	for _, u := range users {
		if u.ID == dave.ID {
			for _, r := range u.Roles {
				daveRoles = append(daveRoles, r.Name)
			}
			break
		}
	}
	assert.Contains(t, daveRoles, "editor")

	// non-admin (bob, editor) → Forbidden
	err = uc.AssignRole(bob.ID, dave.ID, editorRoleID)
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestRoleUsecase_RevokeRole_AdminOK_OthersForbidden(t *testing.T) {
	uc, alice, bob, _, viewerRoleID, _ := setupRoleEnv(t)

	// admin이 bob의 viewer role 회수 (bob은 editor만 갖고 있어 효과 없지만 호출은 성공)
	require.NoError(t, uc.RevokeRole(alice.ID, bob.ID, viewerRoleID))

	// non-admin → Forbidden
	err := uc.RevokeRole(bob.ID, alice.ID, viewerRoleID)
	assert.ErrorIs(t, err, ErrForbidden)
}
