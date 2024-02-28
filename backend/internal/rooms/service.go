package rooms

import (
	"context"

	"github.com/ServiceWeaver/weaver"
	"github.com/google/uuid"
)

type Roomer interface {
	CreateRoom(ctx context.Context, req CreateRoomRequest) (*uuid.UUID, error)
	RemoveRoom(ctx context.Context, id uuid.UUID) error
}

type service struct {
	weaver.Implements[Roomer]
	repo weaver.Ref[Repository]
}

func (r *service) CreateRoom(ctx context.Context, room CreateRoomRequest) (*uuid.UUID, error) {
	log := r.Logger(ctx)
	log.Debug("Creating room", "room", room)

	id, err := r.repo.Get().Store(ctx, room)
	if err != nil {
		log.Error("Failed to create room", "error", err)
		return nil, err
	}
	log.Info("Room created", "id", id)
	return id, nil
}

func (r *service) RemoveRoom(ctx context.Context, id uuid.UUID) error {
	log := r.Logger(ctx)
	log.Debug("Removing room", "id", id)

	err := r.repo.Get().Remove(ctx, id)
	if err != nil {
		log.Error("Failed to remove room", "error", err)
		return err
	}
	log.Info("Room removed", "id", id)
	return nil
}
