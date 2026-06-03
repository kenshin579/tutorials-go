package repository

import (
	"github.com/kenshin579/tutorials-go/web/sns-login-session/backend/model"
	"gorm.io/gorm"
)

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(s *model.Session) error {
	return r.db.Create(s).Error
}

func (r *SessionRepository) FindByID(id string) (*model.Session, error) {
	var s model.Session
	if err := r.db.First(&s, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SessionRepository) Delete(id string) error {
	return r.db.Delete(&model.Session{}, "id = ?", id).Error
}
