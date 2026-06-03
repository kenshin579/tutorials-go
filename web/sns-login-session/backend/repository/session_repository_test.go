package repository

import (
	"testing"
	"time"

	"github.com/kenshin579/tutorials-go/web/sns-login-session/backend/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("DB 연결 실패: %v", err)
	}
	if err := db.AutoMigrate(&model.User{}, &model.Session{}); err != nil {
		t.Fatalf("마이그레이션 실패: %v", err)
	}
	return db
}

func TestSessionRepository_CreateFindDelete(t *testing.T) {
	repo := NewSessionRepository(newTestDB(t))

	sess := &model.Session{ID: "tok-1", UserID: 5, ExpiresAt: time.Now().Add(time.Hour)}
	if err := repo.Create(sess); err != nil {
		t.Fatalf("생성 실패: %v", err)
	}

	found, err := repo.FindByID("tok-1")
	if err != nil {
		t.Fatalf("조회 실패: %v", err)
	}
	if found.UserID != 5 {
		t.Errorf("UserID 기대 5, 실제 %d", found.UserID)
	}

	if err := repo.Delete("tok-1"); err != nil {
		t.Fatalf("삭제 실패: %v", err)
	}
	if _, err := repo.FindByID("tok-1"); err == nil {
		t.Error("삭제 후에도 조회됨")
	}
}
