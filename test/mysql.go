package test

import (
	"context"
	"testing"

	"github.com/docker/go-connections/nat"
	_ "github.com/go-sql-driver/mysql"
	"github.com/testcontainers/testcontainers-go"
)

func NewMysqlDB() {

}

func createMysqlContainer(ctx context.Context, tb testing.TB, password, database string) (string, string, error) {
	port := "3306"
	req := testcontainers.GenericContainerRequest{
		Logger: testcontainers.TestLogger(tb),
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "mysql:8",
			ExposedPorts: []string{port + "/tcp"},
			// HostConfigModifier: func(config *container.HostConfig) {
			// 	config.AutoRemove = true
			// },
			Env: map[string]string{
				"MYSQL_DATABASE":      database,
				"MYSQL_ROOT_PASSWORD": password,
			},
			AlwaysPullImage: false,
			// WaitingFor: wait.ForSQL(nat.Port(port), "mysql", func(host string, port nat.Port) string {
			// 	return fmt.Sprintf("root:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=true", password, host, port.Port(), database)
			// }).WithPollInterval(3 * time.Second),
		},
		Started: true,
	}

	genericContainer, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return "", "", err
	}

	host, err := genericContainer.Host(ctx)
	if err != nil {
		return "", "", err
	}

	mappedPort, err := genericContainer.MappedPort(ctx, nat.Port(port))
	if err != nil {
		return "", "", err
	}

	return host, mappedPort.Port(), nil
}
