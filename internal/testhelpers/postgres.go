package testhelpers

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	_pgImage    = "docker.io/postgres:16.4"
	_pgDatabase = "test"
	_pgUsername = "test"
	_pgPassword = "test"
)

type PostgresContainer struct {
	*postgres.PostgresContainer
	ConnStr string
	migrate *migrate.Migrate
}

func BuildContainer(ctx context.Context, migrationsSourceUrl string) (*PostgresContainer, error) {
	pgC, err := postgres.Run(
		ctx,
		_pgImage,
		postgres.WithDatabase(_pgDatabase),
		postgres.WithUsername(_pgUsername),
		postgres.WithPassword(_pgPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(5*time.Second),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("testhelpers: postgres: BuildContainer: Run: %w", err)
	}

	connStr, err := pgC.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("testhelpers: postgres: BuildContainer: ConnectionString: %w", err)
	}

	m, err := migrate.New(migrationsSourceUrl, connStr)
	if err != nil {
		return nil, fmt.Errorf("testhelpers: postgres: BuildContainer: New: %w", err)
	}

	return &PostgresContainer{
		PostgresContainer: pgC,
		ConnStr:           connStr,
		migrate:           m,
	}, nil
}

func (c *PostgresContainer) MigrateUp() error {
	if err := c.migrate.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("testhelpers: postgres: MigrateUp: Up: %w", err)
	}

	return nil
}
