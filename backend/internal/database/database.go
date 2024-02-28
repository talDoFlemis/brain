package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
	pgxUUID "github.com/vgarvardt/pgx-google-uuid/v5"
)

type Config struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
}

func generateConnStr(cfg Config) string {
	username := cfg.Username
	password := cfg.Password
	host := cfg.Host
	port := cfg.Port
	database := cfg.Database
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		username,
		password,
		host,
		port,
		database,
	)
	return connStr
}

func NewPool(cfg Config) (*pgxpool.Pool, error) {
	str := generateConnStr(cfg)
	pgxConfig, err := pgxpool.ParseConfig(str)
	if err != nil {
		return nil, err
	}
	pgxConfig.AfterConnect = func(_ context.Context, conn *pgx.Conn) error {
		pgxUUID.Register(conn.TypeMap())
		return nil
	}
	db, err := pgxpool.NewWithConfig(context.Background(), pgxConfig)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func NewConn(cfg Config) (*pgx.Conn, error) {
	str := generateConnStr(cfg)
	pgxConfig, err := pgx.ParseConfig(str)
	if err != nil {
		return nil, err
	}
	conn, err := pgx.ConnectConfig(context.Background(), pgxConfig)
	if err != nil {
		return nil, err
	}
	pgxUUID.Register(conn.TypeMap())
	return conn, nil
}
