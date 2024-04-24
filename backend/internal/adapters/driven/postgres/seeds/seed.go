package seed

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/taldoflemis/brain.test/internal/adapters/driven/postgres"
)

func seedAuth(pool *pgxpool.Pool) error {
	users := []struct {
		username       string
		email          string
		password       string
		hashedPassword string
	}{
		{
			username: "marcelomx30",
			email:    "marcelomx30@gmail.com",
			// mypassword
			hashedPassword: "$2a$12$TSjLw2cqeD5bcjPUgOWaaew3xP88soPytNTnMi27vxcNMCDaLFkBa",
		},
	}

	idp := postgres.NewLocalIDPPostgresStorer(pool)
	for _, user := range users {
		_, err := idp.StoreUser(
			context.Background(),
			user.username,
			user.email,
			user.hashedPassword,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func Seed(connStr, basePath string) error {
	pool, err := postgres.NewPool(connStr)
	if err != nil {
		return err
	}

	err = seedAuth(pool)
	if err != nil {
		return err
	}

	return nil
}
