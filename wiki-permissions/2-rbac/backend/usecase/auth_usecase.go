package usecase

import (
	"errors"
	"time"

	"github.com/kenshin579/tutorials-go/wiki-permissions/2-rbac/backend/domain"
	jwthelper "github.com/kenshin579/tutorials-go/wiki-permissions/2-rbac/backend/pkg/jwt"
	"github.com/kenshin579/tutorials-go/wiki-permissions/2-rbac/backend/pkg/passwordhash"
)

// ErrInvalidCredentials는 잘못된 이메일 또는 비밀번호가 입력되었을 때 반환된다.
// 사용자 존재 여부를 응답 본문(payload) 차원에서 노출하지 않기 위해 두 경우 모두 동일 에러를 사용한다.
// 단, 응답 시간 차이(timing oracle)로 인한 사용자 존재 추론 방지는 본 튜토리얼 범위 밖이다.
var ErrInvalidCredentials = errors.New("invalid credentials")

// LoginResult는 로그인 응답에 필요한 토큰 + 사용자 + 효과적 권한 묶음이다.
// 1편(ACL)의 Login은 (token, *User, error)로 끝났지만, RBAC에서는 클라이언트가
// 사전 게이팅(PermissionGate)에 사용할 permissions를 함께 받아야 하므로 추가했다.
type LoginResult struct {
	Token       string
	User        *domain.User
	Permissions []domain.Permission
}

// AuthUsecase는 로그인 흐름을 담당한다 — 자격 증명 검증 후 JWT + 효과적 권한을 반환한다.
type AuthUsecase struct {
	users     domain.UserRepository
	perms     domain.PermissionRepository
	jwtSecret string
	tokenTTL  time.Duration
}

// NewAuthUsecase는 사용자/권한 저장소, JWT 시크릿, 토큰 만료 기간을 받아 AuthUsecase를 생성한다.
func NewAuthUsecase(users domain.UserRepository, perms domain.PermissionRepository, secret string, ttl time.Duration) *AuthUsecase {
	return &AuthUsecase{users: users, perms: perms, jwtSecret: secret, tokenTTL: ttl}
}

// Login은 email/plainPassword를 검증하고 LoginResult를 반환한다.
func (u *AuthUsecase) Login(email, plainPassword string) (*LoginResult, error) {
	found, err := u.users.FindByEmail(email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	if !passwordhash.Verify(plainPassword, found.PasswordHash) {
		return nil, ErrInvalidCredentials
	}
	tok, err := jwthelper.Issue(found.ID, u.jwtSecret, u.tokenTTL)
	if err != nil {
		return nil, err
	}
	perms, err := u.perms.FindByUserID(found.ID)
	if err != nil {
		return nil, err
	}
	return &LoginResult{Token: tok, User: found, Permissions: perms}, nil
}
