package go_testcontainers

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/go-connections/nat"
	_ "github.com/go-sql-driver/mysql"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	mysqlDriver "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMysqlDB() *gorm.DB {
	password, database := "root", "test_db"
	host, port, err := startMysqlContainer(password, database)
	if err != nil {
		panic(err)
	}

	dsn := fmt.Sprintf("root:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true", password, host, port, database)
	db, err := gorm.Open(mysqlDriver.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}

func startMysqlContainer(password, database string) (string, string, error) {
	ctx := context.Background()
	port := "3306"

	req := testcontainers.ContainerRequest{
		Image:        "mysql:8",
		ExposedPorts: []string{port + "/tcp"},
		Env: map[string]string{
			"MYSQL_DATABASE":      database,
			"MYSQL_ROOT_PASSWORD": password,
		},
		WaitingFor: wait.ForSQL(nat.Port(port), "mysql", func(port nat.Port) string {
			return fmt.Sprintf("root:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true", password, "localhost", port.Port(), database)
		}).WithPollInterval(3 * time.Second),
	}

	mysqlC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return "", "", err
	}

	host, err := mysqlC.Host(ctx)
	if err != nil {
		return "", "", err
	}

	mappedPort, err := mysqlC.MappedPort(ctx, nat.Port(port))
	if err != nil {
		return "", "", err
	}

	return host, mappedPort.Port(), nil
}
