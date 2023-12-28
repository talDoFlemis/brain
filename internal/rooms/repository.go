package rooms

import (
	"context"
	"errors"
	"time"

	"github.com/ServiceWeaver/weaver"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/taldoflemis/gahoot/internal/database"
)

var (
	ErrRoomNotFound = errors.New("room not found")
)

type Repository interface {
	Store(ctx context.Context, req CreateRoomRequest) (*uuid.UUID, error)
	Remove(ctx context.Context, id uuid.UUID) error
}

type repository struct {
	weaver.Implements[Repository]
	weaver.WithConfig[database.Config]

	conn *pgx.Conn
}

func (r *repository) Init(_ context.Context) error {
	conn, err := database.NewConn(*r.Config())
	if err != nil {
		return err
	}
	r.conn = conn
	return nil
}

func (p *repository) Store(
	ctx context.Context,
	req CreateRoomRequest,
) (*uuid.UUID, error) {
	id := uuid.New()
	args := pgx.NamedArgs{
		"id":          id,
		"owner":       req.Owner,
		"name":        req.Name,
		"description": req.Description,
		"created_at":  time.Now(),
		"updated_at":  time.Now(),
	}
	insert := `INSERT INTO rooms (id, owner_id, name, description, created_at, updated_at)
    VALUES (@id, @owner, @name, @description, @created_at, @updated_at);`

	_, err := p.conn.Exec(ctx, insert, args)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (p *repository) Remove(ctx context.Context, id uuid.UUID) error {
	delete := `DELETE FROM rooms WHERE id = $1`
	tag, err := p.conn.Exec(ctx, delete, id)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return ErrRoomNotFound
	}

	return nil
}
