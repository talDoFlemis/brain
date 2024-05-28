package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	game "github.com/taldoflemis/brain.test/internal/core/domain/game_aggregate"
	"github.com/taldoflemis/brain.test/internal/ports"
)

type QuestionKind string

const (
	QuizQuestionKind      QuestionKind = "quiz"
	TrueFalseQuestionKind QuestionKind = "true_false"
)

type PostgresGameStorer struct {
	pool *pgxpool.Pool
}

func NewPostgresGameStorer(pool *pgxpool.Pool) *PostgresGameStorer {
	return &PostgresGameStorer{
		pool: pool,
	}
}

func (p *PostgresGameStorer) StoreGame(
	ctx context.Context,
	game *game.Game,
) error {
	tx, err := p.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	args := pgx.NamedArgs{
		"id":          game.Id,
		"title":       game.Title,
		"description": game.Description,
		"owner_id":    game.OwnerId,
	}

	insert := `INSERT INTO games (id, title, description, owner_id) VALUES (@id, @title, @description, @owner_id)`
	_, err = tx.Exec(ctx, insert, args)

	if err != nil {
		return err
	}

	for i, question := range game.Questions {
		err := p.storeQuestion(ctx, tx.Conn(), game.Id, i, question)
		if err != nil {
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return tx.Rollback(ctx)
	}
	return nil
}

func (p *PostgresGameStorer) UpdateGameInfo(ctx context.Context) error {
	panic("not implemented") // TODO: Implement
}

func (p *PostgresGameStorer) UpdateGameQuestions(ctx context.Context) error {
	panic("not implemented") // TODO: Implement
}

func (p *PostgresGameStorer) DeleteGame(ctx context.Context, id uuid.UUID) error {
	args := pgx.NamedArgs{
		"gameId": id,
	}

	query := "DELETE FROM games WHERE id = @gameId"

	_, err := p.pool.Exec(ctx, query, args)

	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresGameStorer) FindGameById(
	ctx context.Context,
	id uuid.UUID,
) (*game.Game, error) {
	args := pgx.NamedArgs{
		"gameId": id,
	}

	query := `SELECT id, owner_id, title, description FROM games WHERE id = @gameId;`

	row := p.pool.QueryRow(ctx, query, args)

	game := game.Game{Questions: nil}

	err := row.Scan(&game.Id, &game.Title, &game.Description, &game.OwnerId)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ports.ErrGameNotFound
		}

		return nil, err
	}

	return &game, nil
}

func (p *PostgresGameStorer) FindAllGamesByUserId(
	ctx context.Context,
	userId string,
) ([]*game.Game, error) {
	args := pgx.NamedArgs{
		"userId": userId,
	}

	query := `SELECT id, owner_id, title, description FROM games WHERE owner_id = @userId;`

	rows, err := p.pool.Query(ctx, query, args)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	games := []*game.Game{}

	for rows.Next() {
		game := game.Game{Questions: nil}

		err := rows.Scan(&game.Id, &game.Title, &game.Description, &game.OwnerId)

		if err != nil {
			return nil, err
		}
		games = append(games, &game)
	}

	return games, nil
}

func (p *PostgresGameStorer) storeQuestion(
	ctx context.Context,
	conn *pgx.Conn,
	gameId uuid.UUID,
	order int,
	question game.Question,
) error {
	var err error
	var kind QuestionKind
	id := uuid.New()

	switch q := question.(type) {
	case *game.QuizQuestion:
		kind = QuizQuestionKind
		err = p.storeQuizQuestion(ctx, conn, id, q)
	case *game.TrueFalseQuestion:
		kind = TrueFalseQuestionKind
		err = p.storeTrueFalseQuestion(ctx, conn, id, q)
	default:
		err = ports.ErrUnknownQuestionKind
	}

	if err != nil {
		return err
	}

	args := pgx.NamedArgs{
		"id":         id,
		"game_id":    gameId,
		"order":      order,
		"title":      question.GetTitle(),
		"time_limit": question.GetTimeLimit(),
		"points":     question.GetPoints(),
		"kind":       kind,
	}

	insert := `INSERT INTO questions (id, game_id, "order", kind, title, time_limit, points) VALUES (@id, @game_id, @order, @kind, @title, @time_limit, @points)`
	_, err = conn.Exec(ctx, insert, args)

	return err
}

func (p *PostgresGameStorer) storeQuizQuestion(
	ctx context.Context,
	conn *pgx.Conn,
	questionId uuid.UUID,
	question *game.QuizQuestion,
) error {
	INSERT := `INSERT INTO quiz_questions (id, question_id, "order", data, correct) VALUES (@id, @question_id, @order, @data, @correct)`

	for i, alternative := range question.Alternatives {
		args := pgx.NamedArgs{
			"id":          uuid.New(),
			"question_id": questionId,
			"order":       i,
			"data":        alternative.Data,
			"correct":     alternative.IsCorrect,
		}

		_, err := conn.Exec(ctx, INSERT, args)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *PostgresGameStorer) storeTrueFalseQuestion(
	ctx context.Context,
	conn *pgx.Conn,
	questionId uuid.UUID,
	question *game.TrueFalseQuestion,
) error {
	args := pgx.NamedArgs{
		"question_id":       questionId,
		"true_alternative":  question.TrueAlternative,
		"false_alternative": question.FalseAlternative,
	}

	insert := `INSERT INTO true_false_questions (question_id, true_alternative, false_alternative) VALUES (@question_id, @true_alternative, @false_alternative)`

	_, err := conn.Exec(ctx, insert, args)
	return err
}
