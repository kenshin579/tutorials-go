package usecase

import "github.com/kenshin579/tutorials-go/rbac/backend/domain"

// UserUsecase defines user management operations.
type UserUsecase interface {
	GetByID(id uint) (*domain.User, error)
	GetAll() ([]domain.User, error)
	Update(user *domain.User) error
	Delete(id uint) error
	AssignRole(userID, roleID uint) error
	RemoveRole(userID, roleID uint) error
}

type userUsecase struct {
	userRepo domain.UserRepository
}

// NewUserUsecase creates a new UserUsecase.
func NewUserUsecase(userRepo domain.UserRepository) UserUsecase {
	return &userUsecase{userRepo: userRepo}
}

func (u *userUsecase) GetByID(id uint) (*domain.User, error) {
	return u.userRepo.FindByID(id)
}

func (u *userUsecase) GetAll() ([]domain.User, error) {
	return u.userRepo.FindAll()
}

func (u *userUsecase) Update(user *domain.User) error {
	return u.userRepo.Update(user)
}

func (u *userUsecase) Delete(id uint) error {
	return u.userRepo.Delete(id)
}

func (u *userUsecase) AssignRole(userID, roleID uint) error {
	return u.userRepo.AssignRole(userID, roleID)
}

func (u *userUsecase) RemoveRole(userID, roleID uint) error {
	return u.userRepo.RemoveRole(userID, roleID)
}
