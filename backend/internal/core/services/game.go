package services

import (
	"context"

	"github.com/google/uuid"

	"github.com/taldoflemis/brain.test/internal/core/domain/game_aggregate"
	"github.com/taldoflemis/brain.test/internal/ports"
)

type GameService struct {
	logger            ports.Logger
	validationService *ValidationService
	gameStorer        ports.GameStorer
}

func NewGameService(
	logger ports.Logger,
	validationService *ValidationService,
	gameStorer ports.GameStorer,
) *GameService {
	return &GameService{
		logger:            logger,
		validationService: validationService,
		gameStorer:        gameStorer,
	}
}

func (s *GameService) CreateNewGame(
	ctx context.Context,
	userId string,
	req *game.Game,
) error {
	req.OwnerId = userId
	req.Id = uuid.New()

	err := s.validationService.Validate(req)
	if err != nil {
		s.logger.Errorf("Failed to validate game %v", err)
		return err
	}

	err = s.gameStorer.StoreGame(ctx, req)
	if err != nil {
		s.logger.Errorf("Failed to store game %v", err)
		return err
	}

	return nil
}

func (s *GameService) GetGamesByUserId(
	ctx context.Context,
	userId string,
) ([]*game.Game, error) {
	return s.gameStorer.FindAllGamesByUserId(ctx, userId)
}

func (s *GameService) GetGameById(
	ctx context.Context,
	gameId uuid.UUID,
) (*game.Game, error) {
	return s.gameStorer.FindGameById(ctx, gameId)
}
