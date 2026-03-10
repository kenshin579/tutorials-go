package usecase

import "github.com/kenshin579/tutorials-go/rbac/backend/domain"

// RbacUsecase defines RBAC management operations.
type RbacUsecase interface {
	// Roles
	CreateRole(role *domain.Role) error
	GetRoleByID(id uint) (*domain.Role, error)
	GetAllRoles() ([]domain.Role, error)
	UpdateRole(role *domain.Role) error
	DeleteRole(id uint) error
	AssignPermission(roleID, permissionID uint) error
	RemovePermission(roleID, permissionID uint) error

	// Permissions
	GetAllPermissions() ([]domain.Permission, error)
}

type rbacUsecase struct {
	roleRepo domain.RoleRepository
	permRepo domain.PermissionRepository
}

// NewRbacUsecase creates a new RbacUsecase.
func NewRbacUsecase(roleRepo domain.RoleRepository, permRepo domain.PermissionRepository) RbacUsecase {
	return &rbacUsecase{
		roleRepo: roleRepo,
		permRepo: permRepo,
	}
}

func (u *rbacUsecase) CreateRole(role *domain.Role) error {
	return u.roleRepo.Create(role)
}

func (u *rbacUsecase) GetRoleByID(id uint) (*domain.Role, error) {
	return u.roleRepo.FindByID(id)
}

func (u *rbacUsecase) GetAllRoles() ([]domain.Role, error) {
	return u.roleRepo.FindAll()
}

func (u *rbacUsecase) UpdateRole(role *domain.Role) error {
	return u.roleRepo.Update(role)
}

func (u *rbacUsecase) DeleteRole(id uint) error {
	return u.roleRepo.Delete(id)
}

func (u *rbacUsecase) AssignPermission(roleID, permissionID uint) error {
	return u.roleRepo.AssignPermission(roleID, permissionID)
}

func (u *rbacUsecase) RemovePermission(roleID, permissionID uint) error {
	return u.roleRepo.RemovePermission(roleID, permissionID)
}

func (u *rbacUsecase) GetAllPermissions() ([]domain.Permission, error) {
	return u.permRepo.FindAll()
}
