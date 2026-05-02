package repository

import (
	"errors"

	"gorm.io/gorm"

	"github.com/kenshin579/tutorials-go/wiki-permissions/3-abac/backend/domain"
)

// UserRepository는 GORM 기반 domain.UserRepository 구현체다.
type UserRepository struct{ db *gorm.DB }

var _ domain.UserRepository = (*UserRepository)(nil)

func NewUserRepository(db *gorm.DB) *UserRepository { return &UserRepository{db: db} }

func (r *UserRepository) Create(u *domain.User) error { return r.db.Create(u).Error }

// FindByEmail은 email로 사용자를 조회하며 Department를 함께 로딩한다.
func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	var u domain.User
	if err := r.db.Preload("Department").Where("email = ?", email).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound{Resource: "user"}
		}
		return nil, err
	}
	return &u, nil
}

// FindByID는 id로 사용자를 조회하며 Department를 함께 로딩한다.
func (r *UserRepository) FindByID(id uint) (*domain.User, error) {
	var u domain.User
	if err := r.db.Preload("Department").First(&u, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound{Resource: "user"}
		}
		return nil, err
	}
	return &u, nil
}
