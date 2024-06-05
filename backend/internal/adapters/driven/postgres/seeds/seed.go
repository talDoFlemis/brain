package seed

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/taldoflemis/brain.test/internal/adapters/driven/postgres"
	game "github.com/taldoflemis/brain.test/internal/core/domain/game_aggregate"
	"github.com/taldoflemis/brain.test/internal/ports"
)

func seedAuth(pool *pgxpool.Pool) (*ports.LocalIDPUserEntity, error) {
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
	user, err := idp.StoreUser(
		context.Background(),
		users[0].username,
		users[0].email,
		users[0].hashedPassword,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func seedGames(pool *pgxpool.Pool, userId string) error {
	games := []game.Game{
		{
			Id:          uuid.New(),
			Title:       "Test game 01",
			Description: "Game generated from the seed file",
			OwnerId:     userId,
			Questions:   nil,
		},
		{
			Id:          uuid.New(),
			Title:       "Test game 02",
			Description: "Game generated from the seed file",
			OwnerId:     userId,
			Questions:   nil,
		},
		{
			Id:          uuid.New(),
			Title:       "Test game 03",
			Description: "Game generated from the seed file",
			OwnerId:     userId,
			Questions:   nil,
		},
	}

	game_storer := postgres.NewPostgresGameStorer(pool)

	for _, game := range games {
		err := game_storer.StoreGame(context.Background(), &game)

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

	user, err := seedAuth(pool)
	if err != nil {
		return err
	}

	err = seedGames(pool, user.ID)
	if err != nil {
		return err
	}

	return nil
}
