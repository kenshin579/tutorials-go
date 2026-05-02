package usecase

import (
	"github.com/kenshin579/tutorials-go/wiki-permissions/2-rbac/backend/domain"
)

// RoleUsecase는 사용자 role 관리(목록 조회, 부여, 회수)를 담당한다.
// 모든 메서드는 users:manage 권한이 있는 사용자(=admin)만 호출 가능하다.
type RoleUsecase struct {
	users domain.UserRepository
	roles domain.RoleRepository
	perms domain.PermissionRepository
}

// NewRoleUsecase는 user/role/permission 저장소를 주입받아 RoleUsecase를 생성한다.
func NewRoleUsecase(users domain.UserRepository, roles domain.RoleRepository, perms domain.PermissionRepository) *RoleUsecase {
	return &RoleUsecase{users: users, roles: roles, perms: perms}
}

// requireAdmin은 요청자가 users:manage 권한을 갖는지 확인한다.
func (u *RoleUsecase) requireAdmin(requesterID uint) error {
	ps, err := u.perms.FindByUserID(requesterID)
	if err != nil {
		return err
	}
	if !HasPermission(ps, "users:manage") {
		return ErrForbidden
	}
	return nil
}

// ListUsers는 admin에게 모든 사용자(+ 각자의 roles)를 반환한다.
func (u *RoleUsecase) ListUsers(requesterID uint) ([]domain.User, error) {
	if err := u.requireAdmin(requesterID); err != nil {
		return nil, err
	}
	return u.users.List()
}

// ListRoles는 admin에게 모든 role(+ 각 role의 permissions)을 반환한다.
func (u *RoleUsecase) ListRoles(requesterID uint) ([]domain.Role, error) {
	if err := u.requireAdmin(requesterID); err != nil {
		return nil, err
	}
	return u.roles.List()
}

// AssignRole은 admin이 대상 사용자에게 role을 부여한다.
func (u *RoleUsecase) AssignRole(requesterID, targetUserID, roleID uint) error {
	if err := u.requireAdmin(requesterID); err != nil {
		return err
	}
	return u.users.AssignRole(targetUserID, roleID)
}

// RevokeRole은 admin이 대상 사용자에게서 role을 회수한다.
func (u *RoleUsecase) RevokeRole(requesterID, targetUserID, roleID uint) error {
	if err := u.requireAdmin(requesterID); err != nil {
		return err
	}
	return u.users.RevokeRole(targetUserID, roleID)
}
