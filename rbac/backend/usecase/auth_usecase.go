package usecase

import (
	"errors"

	"github.com/kenshin579/tutorials-go/rbac/backend/domain"
	"github.com/kenshin579/tutorials-go/rbac/backend/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidToken       = errors.New("invalid token")
)

// AuthUsecase defines authentication operations.
type AuthUsecase interface {
	Register(email, password, name string) (*domain.User, error)
	Login(email, password, jwtSecret string) (*jwt.TokenPair, *domain.User, error)
	Refresh(refreshToken, jwtSecret string) (*jwt.TokenPair, error)
}

type authUsecase struct {
	userRepo domain.UserRepository
	roleRepo domain.RoleRepository
}

// NewAuthUsecase creates a new AuthUsecase.
func NewAuthUsecase(userRepo domain.UserRepository, roleRepo domain.RoleRepository) AuthUsecase {
	return &authUsecase{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

func (u *authUsecase) Register(email, password, name string) (*domain.User, error) {
	// Check if email already exists
	existing, _ := u.userRepo.FindByEmail(email)
	if existing != nil {
		return nil, ErrEmailAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:        email,
		PasswordHash: string(hash),
		Name:         name,
	}

	if err := u.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Assign default "user" role
	role, err := u.roleRepo.FindByName("user")
	if err == nil && role != nil {
		_ = u.userRepo.AssignRole(user.ID, role.ID)
	}

	// Reload user with roles
	user, err = u.userRepo.FindByID(user.ID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *authUsecase) Login(email, password, jwtSecret string) (*jwt.TokenPair, *domain.User, error) {
	user, err := u.userRepo.FindByEmail(email)
	if err != nil {
		return nil, nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, nil, ErrInvalidCredentials
	}

	// Extract role names
	roles := make([]string, len(user.Roles))
	for i, r := range user.Roles {
		roles[i] = r.Name
	}

	pair, err := jwt.GenerateTokenPair(user.ID, roles, jwtSecret)
	if err != nil {
		return nil, nil, err
	}

	return pair, user, nil
}

func (u *authUsecase) Refresh(refreshToken, jwtSecret string) (*jwt.TokenPair, error) {
	claims, err := jwt.ParseToken(refreshToken, jwtSecret)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// Reload user to get current roles
	user, err := u.userRepo.FindByID(claims.UserID)
	if err != nil {
		return nil, err
	}

	roles := make([]string, len(user.Roles))
	for i, r := range user.Roles {
		roles[i] = r.Name
	}

	pair, err := jwt.GenerateTokenPair(user.ID, roles, jwtSecret)
	if err != nil {
		return nil, err
	}

	return pair, nil
}
