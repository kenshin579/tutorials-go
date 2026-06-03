package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/kenshin579/tutorials-go/web/sns-login-session/backend/model"
	"github.com/kenshin579/tutorials-go/web/sns-login-session/backend/repository"
)

type SessionService struct {
	repo   *repository.SessionRepository
	expiry time.Duration
}

func NewSessionService(repo *repository.SessionRepository, expiry time.Duration) *SessionService {
	return &SessionService{repo: repo, expiry: expiry}
}

func (s *SessionService) Create(userID uint) (*model.Session, error) {
	sess := &model.Session{
		ID:        generateSessionID(),
		UserID:    userID,
		ExpiresAt: time.Now().Add(s.expiry),
	}
	if err := s.repo.Create(sess); err != nil {
		return nil, err
	}
	return sess, nil
}

// Validate는 세션 ID로 사용자 ID를 반환한다. 만료/없음이면 에러.
func (s *SessionService) Validate(id string) (uint, error) {
	sess, err := s.repo.FindByID(id)
	if err != nil {
		return 0, err
	}
	if time.Now().After(sess.ExpiresAt) {
		_ = s.repo.Delete(id) // 만료 세션 정리
		return 0, errors.New("만료된 세션")
	}
	return sess.UserID, nil
}

func (s *SessionService) Delete(id string) error {
	return s.repo.Delete(id)
}

func generateSessionID() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
