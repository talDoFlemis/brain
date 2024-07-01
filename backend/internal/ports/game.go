package ports

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/taldoflemis/brain.test/internal/core/domain/game_aggregate"
)

var (
	ErrUnknownQuestionKind = errors.New("Unknown question type")
	ErrGameNotFound        = errors.New("Game not found")
)

type GameStorer interface {
	StoreGame(ctx context.Context, game *game.Game) error
	UpdateGameInfo(ctx context.Context) error
	UpdateGameQuestions(ctx context.Context) error
	DeleteGame(ctx context.Context, id uuid.UUID) error
	FindGameById(ctx context.Context, id uuid.UUID) (*game.Game, error)
	FindAllGamesByUserId(ctx context.Context, userId string) ([]*game.Game, error)
}
