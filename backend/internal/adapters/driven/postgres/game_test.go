package postgres

import (
	"context"
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	testcontainers "github.com/testcontainers/testcontainers-go/modules/postgres"

	game "github.com/taldoflemis/brain.test/internal/core/domain/game_aggregate"
	"github.com/taldoflemis/brain.test/test/helpers"
)

var (
	testGameTitle       = "testGameTitle"
	testGameId          = uuid.New()
	testGameDescription = "testGameDescription"
	testGameID          = uuid.New()
)

type PostgresGameStorerTestSuite struct {
	suite.Suite
	pgContainer *testcontainers.PostgresContainer
	ctx         context.Context
	repo        *PostgresGameStorer
	pool        *pgxpool.Pool
}

func (suite *PostgresGameStorerTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	pgContainer, pool, err := testshelpers.CreatePostgresContainerAndMigrate(
		suite.ctx,
		"./migrations/",
	)
	if err != nil {
		log.Fatal(err)
	}
	suite.pgContainer = pgContainer
	repo := NewPostgresGameStorer(pool)

	suite.pool = pool
	suite.repo = repo
}

func (suite *PostgresGameStorerTestSuite) TearDownTest() {
	_, err := suite.pool.Exec(suite.ctx, "TRUNCATE TABLE games CASCADE")
	if err != nil {
		log.Fatalf("error truncating users table: %s", err)
	}
}

func (suite *PostgresGameStorerTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating postgres container: %s", err)
	}
}

func TestPostgreGameStorerTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("too slow for testing.Short")
	}

	suite.Run(t, new(PostgresGameStorerTestSuite))
}

func (suite *PostgresGameStorerTestSuite) generateMockedGame() *game.Game {
	return &game.Game{
		Id:          testGameId,
		Title:       testGameTitle,
		Description: testGameDescription,
		OwnerId:     testGameID.String(),
	}
}

func (suite *PostgresGameStorerTestSuite) assertQuestions(
	t *testing.T,
	expected game.Question,
	actual game.Question,
) {
	assert.Equal(t, expected.GetTitle(), actual.GetTitle())
	assert.Equal(t, expected.GetPoints(), actual.GetPoints())
	assert.Equal(t, expected.GetTimeLimit(), actual.GetTimeLimit())
}

func (suite *PostgresGameStorerTestSuite) TestStoreGameWithoutQuestions() {
	// Arrange
	t := suite.T()
	mockedGame := suite.generateMockedGame()

	// Act
	err := suite.repo.StoreGame(suite.ctx, mockedGame)

	// Assert
	assert.NoError(t, err)

	var insertedGame game.Game
	err = suite.pool.QueryRow(
		suite.ctx,
		`SELECT id, title, description, owner_id FROM games WHERE id = $1`,
		testGameId,
	).Scan(&insertedGame.Id, &insertedGame.Title, &insertedGame.Description, &insertedGame.OwnerId)
	assert.NoError(t, err)

	assert.Equal(t, mockedGame.Id, insertedGame.Id)
	assert.Equal(t, mockedGame.Title, insertedGame.Title)
	assert.Equal(t, mockedGame.Description, insertedGame.Description)
	assert.Equal(t, mockedGame.OwnerId, insertedGame.OwnerId)
}

func (suite *PostgresGameStorerTestSuite) TestStoreGameWithTrueFalseQuestions() {
	// Arrange
	t := suite.T()
	mockedGame := suite.generateMockedGame()

	table := []struct {
		gameId      uuid.UUID
		description string
		questions   []game.Question
	}{
		{
			gameId:      uuid.New(),
			description: "single true or false question",
			questions: []game.Question{
				&game.TrueFalseQuestion{
					Title:            "testQuestion",
					Points:           1,
					TimeLimit:        30,
					TrueAlternative:  "testTrueAlternative",
					FalseAlternative: "testFalseAlternative",
				},
			},
		},
		{
			gameId:      uuid.New(),
			description: "multiple true or false questions",
			questions: []game.Question{
				&game.TrueFalseQuestion{
					Title:            "testQuestion",
					Points:           1,
					TimeLimit:        30,
					TrueAlternative:  "testTrueAlternative",
					FalseAlternative: "testFalseAlternative",
				},
				&game.TrueFalseQuestion{
					Title:            "tutupa",
					Points:           1,
					TimeLimit:        30,
					TrueAlternative:  "me lhamo lg",
					FalseAlternative: "tubias here",
				},
			},
		},
	}

	for _, tt := range table {
		mockedGame.Questions = tt.questions
		mockedGame.Id = tt.gameId
		t.Run(tt.description, func(t *testing.T) {
			// Act
			err := suite.repo.StoreGame(suite.ctx, mockedGame)

			// Assert
			assert.NoError(t, err)

			rows, err := suite.pool.Query(
				suite.ctx,
				`SELECT id, title, time_limit, points FROM questions WHERE game_id = $1`,
				tt.gameId,
			)
			assert.NoError(t, err)

			for _, mockedQuestion := range mockedGame.Questions {
				rows.Next()
				var actualQ game.TrueFalseQuestion
				err = rows.Scan(
					&actualQ.Id,
					&actualQ.Title,
					&actualQ.TimeLimit,
					&actualQ.Points,
				)
				assert.NoError(t, err)
				suite.assertQuestions(t, mockedQuestion, &actualQ)

				err := suite.pool.QueryRow(
					suite.ctx,
					`SELECT true_alternative, false_alternative FROM true_false_questions WHERE question_id = $1`,
					actualQ.Id,
				).Scan(&actualQ.TrueAlternative, &actualQ.FalseAlternative)
				assert.NoError(t, err)
				assert.Equal(
					t,
					mockedQuestion.(*game.TrueFalseQuestion).TrueAlternative,
					actualQ.TrueAlternative,
				)
				assert.Equal(
					t,
					mockedQuestion.(*game.TrueFalseQuestion).FalseAlternative,
					actualQ.FalseAlternative,
				)
			}

			assert.False(t, rows.Next())
		})
	}
}

func (suite *PostgresGameStorerTestSuite) TestStoreGameWithQuizQuestions() {
	// Arrange
	t := suite.T()
	mockedGame := suite.generateMockedGame()
	mockedGame.Questions = []game.Question{
		&game.QuizQuestion{
			Id:        uuid.New(),
			Title:     "quiz 1",
			Points:    1,
			TimeLimit: 30,
			Alternatives: []game.Alternative{
				{
					Data:      "testAlternative 3",
					IsCorrect: true,
				},
				{
					Data:      "testAlternative 4",
					IsCorrect: true,
				},
				{
					Data:      "testAlternative 5",
					IsCorrect: true,
				},
			},
		},
		&game.QuizQuestion{
			Id:        uuid.New(),
			Title:     "quiz 2",
			Points:    1,
			TimeLimit: 30,
			Alternatives: []game.Alternative{
				{
					Data:      "testAlternative 7",
					IsCorrect: true,
				},
				{
					Data:      "testAlternative 8",
					IsCorrect: true,
				},
			},
		},
	}

	// Act
	err := suite.repo.StoreGame(suite.ctx, mockedGame)

	// Assert
	assert.NoError(t, err)

	rows, err := suite.pool.Query(
		suite.ctx,
		`SELECT id FROM questions WHERE game_id = $1`,
		mockedGame.Id)
	assert.NoError(t, err)

	for _, mockedQuestion := range mockedGame.Questions {
		rows.Next()
		var id uuid.UUID
		err = rows.Scan(&id)
		assert.NoError(t, err)
		alternatives, err := suite.pool.Query(
			suite.ctx,
			`SELECT data, correct FROM quiz_questions WHERE question_id = $1`,
			id,
		)
		mockedQuestion.GetTitle()
		assert.NoError(t, err)
		for _, mockedAlternative := range mockedQuestion.(*game.QuizQuestion).Alternatives {
			alternatives.Next()
			var actualA game.Alternative

			err = alternatives.Scan(&actualA.Data, &actualA.IsCorrect)
			assert.NoError(t, err)
			assert.Equal(t, mockedAlternative, actualA)
		}
	}
	assert.False(t, rows.Next())
}

func (suite *PostgresGameStorerTestSuite) TestStoreGameWithDiverseQuestions() {
	t := suite.T()
	mockedGame := suite.generateMockedGame()

	trueFalseQuestions := []game.Question{
		&game.TrueFalseQuestion{
			Title:            "testQuestion",
			Points:           1,
			TimeLimit:        30,
			TrueAlternative:  "testTrueAlternative",
			FalseAlternative: "testFalseAlternative",
		},
		&game.TrueFalseQuestion{
			Title:            "tutupa",
			Points:           1,
			TimeLimit:        30,
			TrueAlternative:  "me lhamo lg",
			FalseAlternative: "tubias here",
		},
	}

	quizQuestions := []game.Question{
		&game.QuizQuestion{
			Id:        uuid.New(),
			Title:     "quiz 1",
			Points:    1,
			TimeLimit: 30,
			Alternatives: []game.Alternative{
				{
					Data:      "testAlternative 3",
					IsCorrect: true,
				},
				{
					Data:      "testAlternative 4",
					IsCorrect: true,
				},
				{
					Data:      "testAlternative 5",
					IsCorrect: true,
				},
			},
		},
		&game.QuizQuestion{
			Id:        uuid.New(),
			Title:     "quiz 2",
			Points:    1,
			TimeLimit: 30,
			Alternatives: []game.Alternative{
				{
					Data:      "testAlternative 7",
					IsCorrect: true,
				},
				{
					Data:      "testAlternative 8",
					IsCorrect: true,
				},
			},
		},
	}

	var q []game.Question
	q = append(q, trueFalseQuestions...)
	q = append(q, quizQuestions...)

	mockedGame.Questions = q

	// Act
	err := suite.repo.StoreGame(suite.ctx, mockedGame)

	// Assert
	assert.NoError(t, err)
	var amount int

	err = suite.pool.QueryRow(suite.ctx, `SELECT COUNT(*) FROM questions WHERE game_id = $1`, mockedGame.Id).
		Scan(&amount)
	assert.NoError(t, err)
	assert.Equal(t, len(q), amount)

	err = suite.pool.QueryRow(suite.ctx, `SELECT COUNT(DISTINCT question_id) FROM quiz_questions`).
		Scan(&amount)
	assert.NoError(t, err)
	assert.Equal(t, len(quizQuestions), amount)

	err = suite.pool.QueryRow(suite.ctx, `SELECT COUNT(*) FROM true_false_questions`).Scan(&amount)
	assert.NoError(t, err)
	assert.Equal(t, len(trueFalseQuestions), amount)
}
