package game

import (
	"github.com/google/uuid"
)

type Game struct {
	Id          uuid.UUID  `json:"id"          validate:"required,uuid4"`
	Title       string     `json:"title"       validate:"required,gte=1,lte=120"`
	Description string     `json:"description" validate:"min=0,max=200,omitempty"`
	OwnerId     string     `json:"owner_id"    validate:"required,min=1"`
	Questions   []Question `json:"questions"   validate:"required,min=1,dive"`
}
