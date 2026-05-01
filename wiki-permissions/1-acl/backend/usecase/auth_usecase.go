package usecase

import (
	"errors"
	"time"

	"github.com/kenshin579/tutorials-go/wiki-permissions/1-acl/backend/domain"
	jwthelper "github.com/kenshin579/tutorials-go/wiki-permissions/1-acl/backend/pkg/jwt"
	"github.com/kenshin579/tutorials-go/wiki-permissions/1-acl/backend/pkg/passwordhash"
)

// ErrInvalidCredentials는 잘못된 이메일 또는 비밀번호가 입력되었을 때 반환된다.
// 사용자 존재 여부를 응답 본문(payload) 차원에서 노출하지 않기 위해
// 두 경우(이메일 미존재 / 비번 불일치) 모두 동일 에러를 사용한다.
// 단, 응답 시간 차이(timing oracle)로 인한 사용자 존재 추론 방지는 본 튜토리얼 범위 밖이다.
var ErrInvalidCredentials = errors.New("invalid credentials")

// AuthUsecase는 로그인 흐름을 담당한다. 이메일+비번 검증 후 JWT access token을 발급한다.
type AuthUsecase struct {
	users     domain.UserRepository
	jwtSecret string
	tokenTTL  time.Duration
}

// NewAuthUsecase는 사용자 저장소, JWT 시크릿, 토큰 만료 기간을 받아 AuthUsecase를 생성한다.
func NewAuthUsecase(users domain.UserRepository, secret string, ttl time.Duration) *AuthUsecase {
	return &AuthUsecase{users: users, jwtSecret: secret, tokenTTL: ttl}
}

// Login은 email/plainPassword를 검증하고 access token + 사용자 정보를 반환한다.
// 자격 증명이 잘못된 경우 ErrInvalidCredentials를 반환한다 (사용자 존재 여부 노출 방지).
func (u *AuthUsecase) Login(email, plainPassword string) (token string, user *domain.User, err error) {
	found, err := u.users.FindByEmail(email)
	if err != nil {
		return "", nil, ErrInvalidCredentials
	}
	if !passwordhash.Verify(plainPassword, found.PasswordHash) {
		return "", nil, ErrInvalidCredentials
	}
	tok, err := jwthelper.Issue(found.ID, u.jwtSecret, u.tokenTTL)
	if err != nil {
		return "", nil, err
	}
	return tok, found, nil
}
