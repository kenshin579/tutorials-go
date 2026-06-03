package service

import (
	"testing"

	"github.com/kenshin579/tutorials-go/web/sns-login-session/backend/provider"
	"github.com/kenshin579/tutorials-go/web/sns-login-session/backend/repository"
)

func TestFindOrCreateUser_UpdatesProfileOnRelogin(t *testing.T) {
	db := newTestDB(t) // session_service_test.go 정의 재사용
	repo := repository.NewUserRepository(db)
	s := &AuthService{userRepo: repo}

	u1, err := s.findOrCreateUser(&provider.UserInfo{
		Email: "a@gmail.com", Name: "홍길동", AvatarURL: "old.png",
		Provider: "google", ProviderID: "g-1",
	})
	if err != nil {
		t.Fatalf("최초 생성 실패: %v", err)
	}

	u2, err := s.findOrCreateUser(&provider.UserInfo{
		Email: "a@gmail.com", Name: "홍길동2", AvatarURL: "new.png",
		Provider: "google", ProviderID: "g-1",
	})
	if err != nil {
		t.Fatalf("재로그인 실패: %v", err)
	}
	if u2.ID != u1.ID {
		t.Errorf("동일 사용자 ID 기대 %d, 실제 %d", u1.ID, u2.ID)
	}
	if u2.Name != "홍길동2" || u2.AvatarURL != "new.png" {
		t.Errorf("프로필 갱신 안 됨: name=%q avatar=%q", u2.Name, u2.AvatarURL)
	}
}
