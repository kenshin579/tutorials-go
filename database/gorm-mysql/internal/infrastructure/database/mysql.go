package database

import (
	"fmt"
	"strings"
	"time"

	"github.com/kenshin579/tutorials-go/database/gorm-mysql/config"
	"github.com/kenshin579/tutorials-go/database/gorm-mysql/internal/domain"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewMySQLDB(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.MySQL.URL), &gorm.Config{
		Logger:      logger.Default.LogMode(parseLogLevel(cfg.MySQL.LogLevel)),
		PrepareStmt: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.User{},
		&domain.Profile{},
		&domain.Post{},
		&domain.Tag{},
	)
}

func parseLogLevel(level string) logger.LogLevel {
	switch strings.ToUpper(level) {
	case "WARN":
		return logger.Warn
	case "ERROR":
		return logger.Error
	case "SILENT":
		return logger.Silent
	default:
		return logger.Info
	}
}
