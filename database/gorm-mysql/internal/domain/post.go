package domain

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Title   string `gorm:"type:varchar(200);not null"`
	Content string `gorm:"type:text"`
	UserID  uint   `gorm:"index;not null"`
	User    User
	Tags    []Tag `gorm:"many2many:post_tags"`
}

type PostRepository interface {
	Create(post *Post) error
	FindByID(id uint) (*Post, error)
	FindByUserID(userID uint) ([]Post, error)
	FindWithTags(id uint) (*Post, error)
	Update(post *Post) error
	Delete(id uint) error
	AddTag(postID uint, tag *Tag) error
	RemoveTags(postID uint) error
}
