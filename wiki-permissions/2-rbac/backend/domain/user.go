package domain

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Email        string    `gorm:"size:255;uniqueIndex;not null" json:"email"`
	Name         string    `gorm:"size:100;not null" json:"name"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	Roles        []Role    `gorm:"many2many:user_roles" json:"roles,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
