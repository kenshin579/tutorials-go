package domain

import "gorm.io/gorm"

type Tag struct {
	gorm.Model
	Name  string `gorm:"type:varchar(50);uniqueIndex;not null"`
	Posts []Post `gorm:"many2many:post_tags"`
}
