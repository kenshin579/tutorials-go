package repository

import (
	"gorm.io/gorm"

	"github.com/kenshin579/tutorials-go/wiki-permissions/2-rbac/backend/domain"
)

// PermissionRepository는 사용자별 효과적 권한(effective permissions)을 조회하는 핵심 저장소다.
// RBAC 평가의 데이터 측면 entry point.
type PermissionRepository struct{ db *gorm.DB }

var _ domain.PermissionRepository = (*PermissionRepository)(nil)

// NewPermissionRepository는 *gorm.DB에서 동작하는 PermissionRepository를 생성한다.
func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
	return &PermissionRepository{db: db}
}

// FindByUserID는 user → role → permission 3-hop JOIN으로 사용자의 모든 효과적 권한을 모아 반환한다.
// 같은 permission이 여러 role에 묶여 중복될 수 있어 Distinct로 정리하고, resource/action 순으로 정렬한다.
//
//   SELECT DISTINCT permissions.* FROM permissions
//     JOIN role_permissions ON role_permissions.permission_id = permissions.id
//     JOIN user_roles       ON user_roles.role_id = role_permissions.role_id
//     WHERE user_roles.user_id = ?
//     ORDER BY permissions.resource, permissions.action
func (r *PermissionRepository) FindByUserID(userID uint) ([]domain.Permission, error) {
	var perms []domain.Permission
	err := r.db.
		Distinct("permissions.*").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ?", userID).
		Order("permissions.resource, permissions.action").
		Find(&perms).Error
	return perms, err
}
