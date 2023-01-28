package containers

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresContainer struct {
	container testcontainers.Container
	Db        *gorm.DB
}

func NewPostgresContainer(database string, ctx context.Context) (*PostgresContainer, error) {
	user := "user"
	password := "Password1!"
	port, err := nat.NewPort("tcp", "5432")

	if err != nil {
		return nil, err
	}

	r := testcontainers.ContainerRequest{
		Image: "postgres:14.1",
		Env: map[string]string{
			"POSTGRES_USER":     user,
			"POSTGRES_PASSWORD": password,
			"POSTGRES_DB":       database,
		},
		ExposedPorts: []string{port.Port()},
		Cmd:          []string{"postgres", "-c", "fsync=off"},
		WaitingFor:   wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(5 * time.Second),
	}

	gr := testcontainers.GenericContainerRequest{
		Started:          true,
		ContainerRequest: r,
	}

	container, err := testcontainers.GenericContainer(ctx, gr)

	if err != nil {
		return nil, err
	}

	host, err := container.Host(ctx)

	if err != nil {
		return nil, err
	}

	containerPort, err := container.MappedPort(ctx, port)

	if err != nil {
		return nil, err
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, containerPort.Port(), user, password, "postgres")

	db, err := gorm.Open(postgres.Open(connStr))

	if err != nil {
		return nil, err
	}

	tx := db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", database))

	if tx.Error != nil {
		panic(tx.Error)
	}
	tx = db.Exec(fmt.Sprintf("CREATE DATABASE %s", database))

	if tx.Error != nil {
		panic(tx.Error)
	}

	return &PostgresContainer{
		container: container,
		Db:        db,
	}, nil

}

func (c *PostgresContainer) Terminate(ctx context.Context) {
	c.container.Terminate(ctx)
}
