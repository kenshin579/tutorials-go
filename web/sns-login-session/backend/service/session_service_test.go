package service

import (
	"testing"
	"time"

	"github.com/kenshin579/tutorials-go/web/sns-login-session/backend/model"
	"github.com/kenshin579/tutorials-go/web/sns-login-session/backend/repository"
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

func TestSessionService_CreateAndValidate(t *testing.T) {
	repo := repository.NewSessionRepository(newTestDB(t))
	s := NewSessionService(repo, time.Hour)

	sess, err := s.Create(11)
	if err != nil {
		t.Fatalf("세션 생성 실패: %v", err)
	}
	if sess.ID == "" {
		t.Fatal("세션 ID가 비어있음")
	}

	userID, err := s.Validate(sess.ID)
	if err != nil {
		t.Fatalf("검증 실패: %v", err)
	}
	if userID != 11 {
		t.Errorf("UserID 기대 11, 실제 %d", userID)
	}
}

func TestSessionService_DeleteInvalidates(t *testing.T) {
	repo := repository.NewSessionRepository(newTestDB(t))
	s := NewSessionService(repo, time.Hour)

	sess, _ := s.Create(11)
	if err := s.Delete(sess.ID); err != nil {
		t.Fatalf("삭제 실패: %v", err)
	}
	if _, err := s.Validate(sess.ID); err == nil {
		t.Error("삭제 후에도 세션이 유효함 (서버측 로그아웃 실패)")
	}
}

func TestSessionService_RejectsExpired(t *testing.T) {
	repo := repository.NewSessionRepository(newTestDB(t))
	s := NewSessionService(repo, -time.Minute) // 이미 만료

	sess, _ := s.Create(11)
	if _, err := s.Validate(sess.ID); err == nil {
		t.Error("만료된 세션이 통과됨")
	}
}
