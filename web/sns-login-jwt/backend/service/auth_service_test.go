package service

import (
	"testing"

	"github.com/kenshin579/tutorials-go/web/sns-login-jwt/backend/model"
	"github.com/kenshin579/tutorials-go/web/sns-login-jwt/backend/provider"
	"github.com/kenshin579/tutorials-go/web/sns-login-jwt/backend/repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("DB 연결 실패: %v", err)
	}
	if err := db.AutoMigrate(&model.User{}); err != nil {
		t.Fatalf("마이그레이션 실패: %v", err)
	}
	return db
}

func TestFindOrCreateUser_UpdatesProfileOnRelogin(t *testing.T) {
	db := newTestDB(t)
	repo := repository.NewUserRepository(db)
	s := &AuthService{userRepo: repo}

	// 최초 로그인 → 생성
	u1, err := s.findOrCreateUser(&provider.UserInfo{
		Email: "a@gmail.com", Name: "홍길동", AvatarURL: "old.png",
		Provider: "google", ProviderID: "g-1",
	})
	if err != nil {
		t.Fatalf("최초 생성 실패: %v", err)
	}

	// 재로그인 → 이름/아바타 변경 반영, 같은 ID 유지
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
