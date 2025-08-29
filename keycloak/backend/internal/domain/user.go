package domain

import "context"

type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type UserRepository interface {
	GetUserByID(ctx context.Context, userID string) (*User, error)
	GetUserInfo(ctx context.Context, token string) (*User, error)
	ValidateToken(ctx context.Context, token string) (bool, error)
}

type UserUseCase interface {
	GetUserInfo(ctx context.Context, token string) (*User, error)
	ValidateToken(ctx context.Context, token string) (bool, error)
}
