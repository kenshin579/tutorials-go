package db

import (
	"fmt"
	"strings"

	"github.com/kenshin579/tutorials-go/go-mysql/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewMysqlDB(cfg *config.Config) *gorm.DB {
	url := cfg.MysqlConfig.Url
	fmt.Println(url)
	db, err := gorm.Open(mysql.Open(url), &gorm.Config{
		Logger:      logger.Default.LogMode(parseLevel(cfg.MysqlConfig.LogLevel)),
		PrepareStmt: true,
	})
	if err != nil {
		panic(err)
	}

	return db
}

func parseLevel(level string) logger.LogLevel {
	switch strings.ToUpper(level) {
	case "WARN":
		return logger.Warn
	case "ERROR":
		return logger.Error
	default:
		return logger.Info
	}
}
