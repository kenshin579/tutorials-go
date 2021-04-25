package database

import (
	"database/sql"
	"fmt"
	"net/url"

	"github.com/spf13/viper"
)

func New(v *viper.Viper) (*sql.DB, error) {
	fmt.Println("db config")
	dbHost := v.GetString(`database.host`)
	dbPort := v.GetString(`database.port`)
	dbUser := v.GetString(`database.user`)
	dbPass := v.GetString(`database.pass`)
	dbName := v.GetString(`database.name`)
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Seoul")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())

	return sql.Open(`mysql`, dsn)

}
