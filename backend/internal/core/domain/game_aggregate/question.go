package game

type Points int
type TimeLimit int64

type Question interface {
	GetTitle() string
	IsCorrect(answer any) bool
	GetPoints() Points
	GetTimeLimit() TimeLimit
}
