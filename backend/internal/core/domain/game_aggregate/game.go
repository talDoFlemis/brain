package game

import (
	"github.com/google/uuid"
)

type Game struct {
	Id          uuid.UUID  `validate:"required,uuid4"`
	Title       string     `validate:"required,gte=1,lte=120"`
	Description string     `validate:"min=1,max=200"`
	OwnerId     string     `validate:"required,min=1"`
	Questions   []Question `validate:"required,min=1,dive"`
}
