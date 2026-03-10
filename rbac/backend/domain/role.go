package domain

import "time"

type Role struct {
	ID          uint         `gorm:"primaryKey" json:"id"`
	Name        string       `gorm:"uniqueIndex;not null" json:"name"`
	Description string       `json:"description"`
	Permissions []Permission `gorm:"many2many:role_permissions" json:"permissions"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type RoleRepository interface {
	Create(role *Role) error
	FindByID(id uint) (*Role, error)
	FindByName(name string) (*Role, error)
	FindAll() ([]Role, error)
	Update(role *Role) error
	Delete(id uint) error
	AssignPermission(roleID, permissionID uint) error
	RemovePermission(roleID, permissionID uint) error
}
