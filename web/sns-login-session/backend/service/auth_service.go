package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"

	"github.com/kenshin579/tutorials-go/web/sns-login-session/backend/model"
	"github.com/kenshin579/tutorials-go/web/sns-login-session/backend/provider"
	"github.com/kenshin579/tutorials-go/web/sns-login-session/backend/repository"
	"gorm.io/gorm"
)

type AuthService struct {
	providers      map[string]provider.OAuthProvider
	userRepo       *repository.UserRepository
	sessionService *SessionService
	states         sync.Map
}

func NewAuthService(
	providers map[string]provider.OAuthProvider,
	userRepo *repository.UserRepository,
	sessionService *SessionService,
) *AuthService {
	return &AuthService{
		providers:      providers,
		userRepo:       userRepo,
		sessionService: sessionService,
	}
}

func (s *AuthService) GetAuthURL(providerName string) (string, error) {
	p, ok := s.providers[providerName]
	if !ok {
		return "", errors.New("지원하지 않는 provider: " + providerName)
	}
	state := generateState()
	s.states.Store(state, true)
	return p.GetAuthURL(state), nil
}

// HandleCallback은 OAuth 콜백을 처리하고 세션을 생성, 세션 ID를 반환한다.
func (s *AuthService) HandleCallback(ctx context.Context, providerName, code, state string) (*model.Session, *model.User, error) {
	if _, ok := s.states.LoadAndDelete(state); !ok {
		return nil, nil, errors.New("유효하지 않은 state")
	}
	p, ok := s.providers[providerName]
	if !ok {
		return nil, nil, errors.New("지원하지 않는 provider: " + providerName)
	}
	userInfo, err := p.ExchangeCode(ctx, code)
	if err != nil {
		return nil, nil, err
	}
	user, err := s.findOrCreateUser(userInfo)
	if err != nil {
		return nil, nil, err
	}
	sess, err := s.sessionService.Create(user.ID)
	if err != nil {
		return nil, nil, err
	}
	return sess, user, nil
}

func (s *AuthService) Logout(sessionID string) error {
	return s.sessionService.Delete(sessionID)
}

func (s *AuthService) GetUser(userID uint) (*model.User, error) {
	return s.userRepo.FindByID(userID)
}

func (s *AuthService) findOrCreateUser(info *provider.UserInfo) (*model.User, error) {
	user, err := s.userRepo.FindByProviderID(info.Provider, info.ProviderID)
	if err == nil {
		user.Name = info.Name
		user.AvatarURL = info.AvatarURL
		if err := s.userRepo.Update(user); err != nil {
			return nil, err
		}
		return user, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	newUser := &model.User{
		Email:      info.Email,
		Name:       info.Name,
		AvatarURL:  info.AvatarURL,
		Provider:   info.Provider,
		ProviderID: info.ProviderID,
	}
	if err := s.userRepo.Create(newUser); err != nil {
		return nil, err
	}
	return newUser, nil
}

func generateState() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
