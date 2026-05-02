package usecase

import (
	"errors"
	"time"

	"github.com/kenshin579/tutorials-go/wiki-permissions/3-abac/backend/domain"
	jwthelper "github.com/kenshin579/tutorials-go/wiki-permissions/3-abac/backend/pkg/jwt"
	"github.com/kenshin579/tutorials-go/wiki-permissions/3-abac/backend/pkg/passwordhash"
)

// ErrInvalidCredentials는 잘못된 이메일 또는 비밀번호가 입력되었을 때 반환된다.
// 사용자 존재 여부를 응답 본문 차원에서 노출하지 않기 위해 두 경우 모두 동일 에러를 사용한다.
// 단, 응답 시간 차이(timing oracle) 방지는 본 튜토리얼 범위 밖이다.
var ErrInvalidCredentials = errors.New("invalid credentials")

// AuthUsecase는 로그인 흐름을 담당한다 — 자격 증명 검증 후 JWT + 사용자(속성 포함)를 반환한다.
//
// 1편(ACL): 응답에 token + user. 2편(RBAC): + permissions/roles 추가.
// 3편(ABAC): user 자체에 department/employment_type 속성이 들어있으므로 별도 권한 응답 불필요.
// 권한 평가는 매 요청 시 ABAC 정책 평가기가 user 속성과 page 속성을 결합 평가한다.
type AuthUsecase struct {
	users     domain.UserRepository
	jwtSecret string
	tokenTTL  time.Duration
}

func NewAuthUsecase(users domain.UserRepository, secret string, ttl time.Duration) *AuthUsecase {
	return &AuthUsecase{users: users, jwtSecret: secret, tokenTTL: ttl}
}

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
