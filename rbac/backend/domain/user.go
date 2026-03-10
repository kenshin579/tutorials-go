package domain

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Email        string    `gorm:"size:255;uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	Name         string    `gorm:"size:100;not null" json:"name"`
	Roles        []Role    `gorm:"many2many:user_roles" json:"roles"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserRepository interface {
	Create(user *User) error
	FindByID(id uint) (*User, error)
	FindByEmail(email string) (*User, error)
	FindAll() ([]User, error)
	Update(user *User) error
	Delete(id uint) error
	AssignRole(userID, roleID uint) error
	RemoveRole(userID, roleID uint) error
}
