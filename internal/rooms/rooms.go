package rooms

import (
	"time"

	"github.com/ServiceWeaver/weaver"
	"github.com/google/uuid"
)

type CreateRoomRequest struct {
	Owner       uuid.UUID `json:"owner"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`

	weaver.AutoMarshal
}

type RoomResponse struct {
	ID          uuid.UUID `json:"id"`
	Owner       uuid.UUID `json:"owner"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Room struct {
	id          uuid.UUID
	owner       uuid.UUID
	name        string
	description string
	createdAt   time.Time
	updatedAt   time.Time
}

func (r *Room) ToRoomResponse() RoomResponse {
	return RoomResponse{
		ID:          r.id,
		CreatedAt:   r.createdAt,
		UpdatedAt:   r.updatedAt,
		Owner:       r.owner,
		Description: r.description,
	}
}
