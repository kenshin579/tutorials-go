package domain

import "time"

type Permission struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Resource    string    `gorm:"size:100;not null" json:"resource"`
	Action      string    `gorm:"size:100;not null" json:"action"`
	Description string    `gorm:"size:255" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// Key returns the permission key (e.g. "users:read")
func (p Permission) Key() string {
	return p.Resource + ":" + p.Action
}

type PermissionRepository interface {
	FindAll() ([]Permission, error)
	FindByUserID(userID uint) ([]Permission, error)
}
