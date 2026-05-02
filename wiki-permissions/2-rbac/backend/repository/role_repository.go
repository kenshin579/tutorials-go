package repository

import (
	"errors"

	"gorm.io/gorm"

	"github.com/kenshin579/tutorials-go/wiki-permissions/2-rbac/backend/domain"
)

// RoleRepository는 GORM 기반 domain.RoleRepository 구현체다.
type RoleRepository struct{ db *gorm.DB }

var _ domain.RoleRepository = (*RoleRepository)(nil)

// NewRoleRepository는 *gorm.DB에서 동작하는 RoleRepository를 생성한다.
func NewRoleRepository(db *gorm.DB) *RoleRepository { return &RoleRepository{db: db} }

// FindByID는 id로 role을 조회하며 Permissions를 함께 로딩한다.
func (r *RoleRepository) FindByID(id uint) (*domain.Role, error) {
	var role domain.Role
	if err := r.db.Preload("Permissions").First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound{Resource: "role"}
		}
		return nil, err
	}
	return &role, nil
}

// FindByName은 name으로 role을 조회하며 Permissions를 함께 로딩한다.
func (r *RoleRepository) FindByName(name string) (*domain.Role, error) {
	var role domain.Role
	if err := r.db.Preload("Permissions").Where("name = ?", name).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound{Resource: "role"}
		}
		return nil, err
	}
	return &role, nil
}

// List는 모든 role을 id 오름차순으로 반환하며 각 role의 Permissions를 함께 로딩한다.
func (r *RoleRepository) List() ([]domain.Role, error) {
	var roles []domain.Role
	err := r.db.Preload("Permissions").Order("id ASC").Find(&roles).Error
	return roles, err
}
