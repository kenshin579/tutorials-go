package repository

import (
	"github.com/kenshin579/tutorials-go/rbac/backend/domain"
	"gorm.io/gorm"
)

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) domain.PermissionRepository {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) FindAll() ([]domain.Permission, error) {
	var permissions []domain.Permission
	err := r.db.Order("resource, action").Find(&permissions).Error
	return permissions, err
}

func (r *permissionRepository) FindByUserID(userID uint) ([]domain.Permission, error) {
	var permissions []domain.Permission
	err := r.db.Raw(`
		SELECT DISTINCT p.*
		FROM permissions p
		JOIN role_permissions rp ON rp.permission_id = p.id
		JOIN user_roles ur ON ur.role_id = rp.role_id
		WHERE ur.user_id = ?
		ORDER BY p.resource, p.action
	`, userID).Scan(&permissions).Error
	return permissions, err
}
