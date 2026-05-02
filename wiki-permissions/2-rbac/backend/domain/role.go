package domain

import "time"

// Role은 사용자에게 부여되는 역할이다 (예: admin, editor, viewer).
// 한 사용자가 여러 role을 가질 수 있고(M:N), 한 role은 여러 permission을 가질 수 있다(M:N).
type Role struct {
	ID          uint         `gorm:"primaryKey" json:"id"`
	Name        string       `gorm:"size:100;uniqueIndex;not null" json:"name"`
	Description string       `gorm:"size:255" json:"description"`
	Permissions []Permission `gorm:"many2many:role_permissions" json:"permissions,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
}
