package repository

import (
	"errors"

	"gorm.io/gorm"

	"github.com/kenshin579/tutorials-go/wiki-permissions/2-rbac/backend/domain"
)

// UserRepository는 GORM 기반 domain.UserRepository 구현체다.
type UserRepository struct{ db *gorm.DB }

var _ domain.UserRepository = (*UserRepository)(nil)

// NewUserRepository는 *gorm.DB에서 동작하는 UserRepository를 생성한다.
func NewUserRepository(db *gorm.DB) *UserRepository { return &UserRepository{db: db} }

// Create는 새 사용자를 저장한다 (저장 후 GORM이 u.ID를 자동 채움).
func (r *UserRepository) Create(u *domain.User) error {
	return r.db.Create(u).Error
}

// FindByEmail은 email로 사용자를 조회하며 Roles를 함께 로딩한다. 없으면 domain.ErrNotFound를 반환한다.
func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	var u domain.User
	if err := r.db.Preload("Roles").Where("email = ?", email).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound{Resource: "user"}
		}
		return nil, err
	}
	return &u, nil
}

// FindByID는 id로 사용자를 조회하며 Roles를 함께 로딩한다. 없으면 domain.ErrNotFound를 반환한다.
func (r *UserRepository) FindByID(id uint) (*domain.User, error) {
	var u domain.User
	if err := r.db.Preload("Roles").First(&u, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound{Resource: "user"}
		}
		return nil, err
	}
	return &u, nil
}

// List는 모든 사용자를 id 오름차순으로 반환하며 각자의 Roles를 함께 로딩한다.
func (r *UserRepository) List() ([]domain.User, error) {
	var users []domain.User
	err := r.db.Preload("Roles").Order("id ASC").Find(&users).Error
	return users, err
}

// AssignRole은 사용자에게 role을 부여한다. 이미 부여된 경우 중복 추가는 GORM이 dedup한다.
func (r *UserRepository) AssignRole(userID, roleID uint) error {
	user := &domain.User{ID: userID}
	role := &domain.Role{ID: roleID}
	return r.db.Model(user).Association("Roles").Append(role)
}

// RevokeRole은 사용자의 role 매핑을 제거한다 (Role 자체는 삭제하지 않음).
func (r *UserRepository) RevokeRole(userID, roleID uint) error {
	user := &domain.User{ID: userID}
	role := &domain.Role{ID: roleID}
	return r.db.Model(user).Association("Roles").Delete(role)
}
