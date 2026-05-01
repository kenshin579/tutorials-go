package config

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/kenshin579/tutorials-go/wiki-permissions/1-acl/backend/domain"
)

func OpenDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&domain.User{}, &domain.Page{}, &domain.ACLEntry{}); err != nil {
		return nil, err
	}
	return db, nil
}
