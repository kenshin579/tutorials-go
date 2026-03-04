package repository

import (
	"github.com/kenshin579/tutorials-go/database/gorm-mysql/internal/domain"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) CreateBatch(users []domain.User) error {
	return r.db.Create(&users).Error
}

func (r *userRepository) FindByID(id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.Preload("Profile").Preload("Posts").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindAll(offset, limit int) ([]domain.User, error) {
	var users []domain.User
	err := r.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error
	return users, err
}

func (r *userRepository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&domain.User{}, id).Error
}

func (r *userRepository) HardDelete(id uint) error {
	return r.db.Unscoped().Delete(&domain.User{}, id).Error
}

func (r *userRepository) CreateWithProfile(user *domain.User) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(user).Error
	})
}
