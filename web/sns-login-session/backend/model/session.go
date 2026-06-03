package model

import "time"

// Session은 서버측 세션 (SQLite 저장)
type Session struct {
	ID        string    `gorm:"primarykey" json:"id"` // 랜덤 세션 토큰
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}
