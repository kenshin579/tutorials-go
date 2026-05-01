package repository

import (
	"errors"

	"gorm.io/gorm"

	"github.com/kenshin579/tutorials-go/wiki-permissions/1-acl/backend/domain"
)

// UserRepository는 GORM 기반 domain.UserRepository 구현체다.
type UserRepository struct{ db *gorm.DB }

var _ domain.UserRepository = (*UserRepository)(nil)

// NewUserRepository는 *gorm.DB에서 동작하는 UserRepository를 생성한다.
func NewUserRepository(db *gorm.DB) *UserRepository { return &UserRepository{db: db} }

// Create는 새 사용자를 저장하고 자동 부여된 ID를 u.ID에 채운다.
func (r *UserRepository) Create(u *domain.User) error {
	return r.db.Create(u).Error
}

// FindByEmail은 email로 사용자를 조회한다. 없으면 domain.ErrNotFound를 반환한다.
func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	var u domain.User
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound{Resource: "user"}
		}
		return nil, err
	}
	return &u, nil
}

// FindByID는 id로 사용자를 조회한다. 없으면 domain.ErrNotFound를 반환한다.
func (r *UserRepository) FindByID(id uint) (*domain.User, error) {
	var u domain.User
	if err := r.db.First(&u, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound{Resource: "user"}
		}
		return nil, err
	}
	return &u, nil
}
