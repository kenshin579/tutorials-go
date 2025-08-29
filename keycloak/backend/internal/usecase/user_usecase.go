package usecase

import (
	"context"

	"github.com/kenshin579/tutorials-go/keycloak/backend/internal/domain"
)

type UserUseCaseImpl struct {
	userRepo domain.UserRepository
}

func NewUserUseCase(userRepo domain.UserRepository) domain.UserUseCase {
	return &UserUseCaseImpl{
		userRepo: userRepo,
	}
}

func (u *UserUseCaseImpl) GetUserInfo(ctx context.Context, token string) (*domain.User, error) {
	return u.userRepo.GetUserInfo(ctx, token)
}

func (u *UserUseCaseImpl) ValidateToken(ctx context.Context, token string) (bool, error) {
	return u.userRepo.ValidateToken(ctx, token)
}
