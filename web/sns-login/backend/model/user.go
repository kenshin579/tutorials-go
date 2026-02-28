package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	Email      string         `gorm:"uniqueIndex;not null" json:"email"`
	Name       string         `json:"name"`
	AvatarURL  string         `json:"avatar_url"`
	Provider   string         `gorm:"not null" json:"provider"`
	ProviderID string         `gorm:"not null;index" json:"provider_id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
