package config

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/kenshin579/tutorials-go/wiki-permissions/3-abac/backend/domain"
)

// OpenDB는 SQLite를 열고 도메인 엔티티(User, Page, Department)에 대해 AutoMigrate를 수행한다.
func OpenDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&domain.Department{}, &domain.User{}, &domain.Page{}); err != nil {
		return nil, err
	}
	return db, nil
}
