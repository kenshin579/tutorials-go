package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"

	"github.com/kenshin579/tutorials-go/web/sns-login/backend/model"
	"github.com/kenshin579/tutorials-go/web/sns-login/backend/provider"
	"github.com/kenshin579/tutorials-go/web/sns-login/backend/repository"
	"gorm.io/gorm"
)

type AuthService struct {
	providers    map[string]provider.OAuthProvider
	userRepo     *repository.UserRepository
	tokenService *TokenService
	states       sync.Map // state 파라미터 저장 (CSRF 방지)
}

func NewAuthService(
	providers map[string]provider.OAuthProvider,
	userRepo *repository.UserRepository,
	tokenService *TokenService,
) *AuthService {
	return &AuthService{
		providers:    providers,
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

// GetAuthURL은 OAuth 인증 URL과 state를 반환한다
func (s *AuthService) GetAuthURL(providerName string) (string, error) {
	p, ok := s.providers[providerName]
	if !ok {
		return "", errors.New("지원하지 않는 provider: " + providerName)
	}

	state := generateState()
	s.states.Store(state, true)

	return p.GetAuthURL(state), nil
}

// HandleCallback은 OAuth 콜백을 처리하고 JWT 토큰을 반환한다
func (s *AuthService) HandleCallback(ctx context.Context, providerName, code, state string) (*TokenPair, *model.User, error) {
	// state 검증 (CSRF 방지)
	if _, ok := s.states.LoadAndDelete(state); !ok {
		return nil, nil, errors.New("유효하지 않은 state")
	}

	p, ok := s.providers[providerName]
	if !ok {
		return nil, nil, errors.New("지원하지 않는 provider: " + providerName)
	}

	// Authorization Code → 사용자 정보 교환
	userInfo, err := p.ExchangeCode(ctx, code)
	if err != nil {
		return nil, nil, err
	}

	// 사용자 조회 또는 생성
	user, err := s.findOrCreateUser(userInfo)
	if err != nil {
		return nil, nil, err
	}

	// JWT 토큰 발급
	tokens, err := s.tokenService.GenerateTokenPair(user.ID)
	if err != nil {
		return nil, nil, err
	}

	return tokens, user, nil
}

// RefreshToken은 Refresh Token으로 새 토큰 쌍을 발급한다
func (s *AuthService) RefreshToken(refreshToken string) (*TokenPair, error) {
	claims, err := s.tokenService.ValidateToken(refreshToken)
	if err != nil {
		return nil, errors.New("유효하지 않은 refresh token")
	}

	return s.tokenService.GenerateTokenPair(claims.UserID)
}

// GetUser는 사용자 ID로 사용자를 조회한다
func (s *AuthService) GetUser(userID uint) (*model.User, error) {
	return s.userRepo.FindByID(userID)
}

func (s *AuthService) findOrCreateUser(info *provider.UserInfo) (*model.User, error) {
	// 기존 사용자 조회
	user, err := s.userRepo.FindByProviderID(info.Provider, info.ProviderID)
	if err == nil {
		return user, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 새 사용자 생성
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
	rand.Read(b)
	return hex.EncodeToString(b)
}
