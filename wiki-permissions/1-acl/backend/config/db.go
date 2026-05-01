package config

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/kenshin579/tutorials-go/wiki-permissions/1-acl/backend/domain"
)

// OpenDB는 주어진 DSN으로 SQLite를 열고 도메인 엔티티(User, Page, ACLEntry)에 대한
// AutoMigrate를 수행한 *gorm.DB를 반환한다. 단위 테스트에서는 ":memory:" DSN을 사용한다.
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
