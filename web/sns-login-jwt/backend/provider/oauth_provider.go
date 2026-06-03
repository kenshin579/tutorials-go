package provider

import "context"

// OAuthProvider는 SNS 로그인 제공자 인터페이스
type OAuthProvider interface {
	// GetAuthURL은 OAuth 인증 URL을 반환한다
	GetAuthURL(state string) string
	// ExchangeCode는 Authorization Code를 사용자 정보로 교환한다
	ExchangeCode(ctx context.Context, code string) (*UserInfo, error)
	// Name은 provider 이름을 반환한다 (google, github 등)
	Name() string
}

// UserInfo는 OAuth provider에서 가져온 사용자 정보
type UserInfo struct {
	Email      string
	Name       string
	AvatarURL  string
	Provider   string
	ProviderID string
}
