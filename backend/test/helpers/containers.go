package testshelpers

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	pgxUUID "github.com/vgarvardt/pgx-google-uuid/v5"

	"github.com/taldoflemis/brain.test/internal/adapters/driven/postgres/migrations"
)

const (
	databaseName = "test-db"
	username     = "postgres"
	password     = "postgres"
)

type PostgresContainer struct {
	*postgres.PostgresContainer
	ConnStr string
}

func CreatePostgresContainer(ctx context.Context) (*postgres.PostgresContainer, string, error) {
	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16.2-alpine"),
		postgres.WithDatabase(databaseName),
		postgres.WithUsername(username),
		postgres.WithPassword(password),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, "", err
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, "", err
	}

	return pgContainer, connStr, nil
}

func NewPool(connStr string) (*pgxpool.Pool, error) {
	pgxConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	pgxConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxUUID.Register(conn.TypeMap())
		return nil
	}

	pool, err := pgxpool.NewWithConfig(context.TODO(), pgxConfig)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func CreatePostgresContainerAndMigrate(
	ctx context.Context,
	migrationsPath string,
) (*postgres.PostgresContainer, *pgxpool.Pool, error) {
	pgContainer, connStr, err := CreatePostgresContainer(ctx)
	if err != nil {
		return nil, nil, err
	}

	migrate.Migrate(connStr, migrationsPath)
	if err != nil {
		return nil, nil, err
	}

	pool, err := NewPool(connStr)
	if err != nil {
		return nil, nil, err
	}

	return pgContainer, pool, nil
}
