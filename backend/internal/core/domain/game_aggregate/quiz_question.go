package game

import (
	"github.com/google/uuid"
)

type QuizQuestion struct {
	Id           uuid.UUID     `validate:"omitempty"`
	Title        string        `validate:"required,gte=1,lte=120"`
	Points       Points        `validate:"required,gte=0,lte=2"`
	TimeLimit    TimeLimit     `validate:"required,gte=5,lte=180"`
	Alternatives []Alternative `validate:"required,min=3,dive,required"`
}

type Alternative struct {
	Data      string `validate:"required,gte=1,lte=120"`
	IsCorrect bool   `validate:"omitempty"`
}

func (q *QuizQuestion) GetTitle() string {
	return q.Title
}

func (q *QuizQuestion) IsCorrect(answer any) bool {
	answerAlternatives, ok := answer.([]*Alternative)
	if !ok {
		return false
	}

	if len(answerAlternatives) != len(q.Alternatives) {
		return false
	}

	for i, correctAlternative := range q.Alternatives {
		candidateAlternative := answerAlternatives[i]

		if candidateAlternative.IsCorrect != correctAlternative.IsCorrect {
			return false
		}
	}

	return true
}

func (q *QuizQuestion) GetPoints() Points {
	return q.Points
}

func (q *QuizQuestion) GetTimeLimit() TimeLimit {
	return q.TimeLimit
}
