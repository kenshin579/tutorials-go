package config

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/kenshin579/tutorials-go/wiki-permissions/2-rbac/backend/domain"
)

// OpenDB는 SQLite를 열고 도메인 엔티티(User, Page, Role, Permission)에 대해 AutoMigrate를 수행한다.
// User.Roles와 Role.Permissions의 many2many 태그가 user_roles, role_permissions 테이블을 자동 생성한다.
func OpenDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(
		&domain.User{},
		&domain.Page{},
		&domain.Role{},
		&domain.Permission{},
	); err != nil {
		return nil, err
	}
	return db, nil
}
