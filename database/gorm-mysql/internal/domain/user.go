package domain

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name    string  `gorm:"type:varchar(100);not null"`
	Email   string  `gorm:"type:varchar(200);uniqueIndex;not null"`
	Profile Profile // Has One
	Posts   []Post  // Has Many
}

type Profile struct {
	gorm.Model
	UserID uint   `gorm:"uniqueIndex;not null"`
	Bio    string `gorm:"type:text"`
	Avatar string `gorm:"type:varchar(500)"`
}

type UserRepository interface {
	Create(user *User) error
	CreateBatch(users []User) error
	FindByID(id uint) (*User, error)
	FindByEmail(email string) (*User, error)
	FindAll(offset, limit int) ([]User, error)
	Update(user *User) error
	Delete(id uint) error
	HardDelete(id uint) error
	CreateWithProfile(user *User) error
}
