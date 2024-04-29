package game

import (
	"github.com/google/uuid"
)

type TrueFalseQuestion struct {
	Id               uuid.UUID `validate:"omitempty"`
	Title            string    `validate:"required,gte=1,lte=120"`
	Points           Points    `validate:"required,gte=0,lte=2"`
	TimeLimit        TimeLimit `validate:"required,gte=5,lte=180"`
	TrueAlternative  string    `validate:"required,gte=1,lte=120"`
	FalseAlternative string    `validate:"required,gte=1,lte=120"`
}

func (t *TrueFalseQuestion) GetTitle() string {
	return t.Title
}

func (t *TrueFalseQuestion) IsCorrect(answer any) bool {
	answerString, ok := answer.(string)
	if !ok {
		return false
	}

	return answerString == t.TrueAlternative
}

func (t *TrueFalseQuestion) GetPoints() Points {
	return t.Points
}

func (t *TrueFalseQuestion) GetTimeLimit() TimeLimit {
	return t.TimeLimit
}
