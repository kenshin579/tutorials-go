package database

import (
	"database/sql"
	"fmt"
	"net/url"

	"github.com/kenshin579/tutorials-go/project-layout/go-clean-arch-v2/pkg/config"
)

func New(cfg *config.Config) (*sql.DB, error) {
	fmt.Println("db config")
	d := cfg.Database
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", d.User, d.Pass, d.Host, d.Port, d.Name)
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Seoul")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())

	return sql.Open("mysql", dsn)
}
