package repository

import (
	"github.com/kenshin579/tutorials-go/web/sns-login-jwt/backend/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// FindByProviderID는 provider와 providerID로 사용자를 조회한다
func (r *UserRepository) FindByProviderID(provider, providerID string) (*model.User, error) {
	var user model.User
	err := r.db.Where("provider = ? AND provider_id = ?", provider, providerID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create는 새 사용자를 생성한다
func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// FindByID는 ID로 사용자를 조회한다
func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update는 사용자 정보를 저장한다
func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}
