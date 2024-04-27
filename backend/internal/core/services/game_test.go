package services_test

import (
	"context"
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	testcontainers "github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/taldoflemis/brain.test/internal/adapters/driven/postgres"
	game "github.com/taldoflemis/brain.test/internal/core/domain/game_aggregate"
	"github.com/taldoflemis/brain.test/internal/core/services"
	testshelpers "github.com/taldoflemis/brain.test/test/helpers"
)

var (
	testGameID = uuid.New()
)

type GameServiceTestSuite struct {
	suite.Suite
	pgContainer *testcontainers.PostgresContainer
	ctx         context.Context
	pool        *pgxpool.Pool
	svc         *services.GameService
}

func (s *GameServiceTestSuite) SetupSuite() {
	s.ctx = context.Background()

	pgContainer, pool, err := testshelpers.CreatePostgresContainerAndMigrate(
		s.ctx,
		"../../adapters/driven/postgres/migrations/",
	)

	if err != nil {
		log.Fatal(err)
	}

	gameStorer := postgres.NewPostgresGameStorer(pool)
	s.pgContainer = pgContainer
	s.pool = pool
	logger := testshelpers.NewDummyLogger(log.Writer())
	s.svc = services.NewGameService(logger, services.NewValidationService(), gameStorer)
}

func (s *GameServiceTestSuite) TearDownTest() {
	_, err := s.pool.Exec(s.ctx, "TRUNCATE TABLE games CASCADE")
	if err != nil {
		log.Fatalf("error truncating games table: %s", err)
	}
}

func (s *GameServiceTestSuite) TearDownSuite() {
	if err := s.pgContainer.Terminate(s.ctx); err != nil {
		log.Fatalf("error terminating postgres container: %s", err)
	}
}

func TestGameServiceIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("too slow for testing.Short")
	}

	suite.Run(t, new(GameServiceTestSuite))
}

func (s *GameServiceTestSuite) generateMockedGame(
	title, description, ownerId string,
	questions []game.Question,
) *game.Game {
	return &game.Game{
		Title:       title,
		Description: description,
		OwnerId:     ownerId,
		Questions:   questions,
	}
}

func (s *GameServiceTestSuite) generateMockedQuizAlternatives() *[]game.Alternative {
	return &[]game.Alternative{
		{
			Data:      "alternativeFirst 15",
			IsCorrect: true,
		},
		{
			Data:      "alternativeSecond 15",
			IsCorrect: false,
		},
		{
			Data:      "alternativeThird 15",
			IsCorrect: true,
		},
		{
			Data:      "alternativeFourth 15",
			IsCorrect: false,
		},
	}
}

func (s *GameServiceTestSuite) TestCreateNewGame() {
	// Arrange
	t := s.T()
	questions := []game.Question{
		&game.TrueFalseQuestion{
			Title:            "testQuestion",
			Points:           2,
			TimeLimit:        30,
			TrueAlternative:  "testTrueAlternative",
			FalseAlternative: "testFalseAlternative",
		},
	}
	mockedGame := s.generateMockedGame("dumb title", "desc here", uuid.NewString(), questions)

	// Act
	err := s.svc.CreateNewGame(s.ctx, testGameID.String(), mockedGame)

	// Assert
	assert.NoError(t, err)

	var count int
	err = s.pool.QueryRow(s.ctx, "SELECT COUNT(*) FROM games").Scan(&count)

	assert.NoError(t, err)
	assert.Equal(t, len(mockedGame.Questions), count)
}

func (s *GameServiceTestSuite) TestCreateNewInvalidGame() {
	t := s.T()
	userID := uuid.New().String()
	table := []struct {
		testDescription string
		title           string
		desc            string
		ownerID         string
		questions       []game.Question
	}{
		{
			testDescription: "game with empty title",
			title:           "",
			desc:            "testDescription 1",
			ownerID:         userID,
			questions: []game.Question{
				&game.TrueFalseQuestion{
					Title:            "testQuestion 1",
					Points:           1,
					TimeLimit:        30,
					TrueAlternative:  "testTrueAlternative 1",
					FalseAlternative: "testFalseAlternative 1",
				},
			},
		},
		{
			testDescription: "game with empty owner_id",
			title:           "title 2",
			desc:            "testDescription 2",
			ownerID:         "",
			questions: []game.Question{
				&game.TrueFalseQuestion{
					Title:            "testQuestion 2",
					Points:           1,
					TimeLimit:        30,
					TrueAlternative:  "testTrueAlternative 2",
					FalseAlternative: "testFalseAlternative 2",
				},
			},
		},
		{
			testDescription: "game with empty questions",
			title:           "title 3",
			desc:            "testDescription 3",
			ownerID:         userID,
			questions:       []game.Question{},
		},
		{
			testDescription: "game with true or false question without title",
			title:           "title 4",
			desc:            "testDescription 4",
			ownerID:         userID,
			questions: []game.Question{
				&game.TrueFalseQuestion{
					Points:           1,
					TimeLimit:        30,
					TrueAlternative:  "testTrueAlternative 4",
					FalseAlternative: "testFalseAlternative 4",
				},
			},
		},
		{
			testDescription: "game with true or false question empty title",
			title:           "title 5",
			desc:            "testDescription 5",
			ownerID:         userID,
			questions: []game.Question{
				&game.TrueFalseQuestion{
					Title:            "",
					Points:           1,
					TimeLimit:        30,
					TrueAlternative:  "testTrueAlternative 5",
					FalseAlternative: "testFalseAlternative 5",
				},
			},
		},
		{
			testDescription: "game with true or false question without points",
			title:           "title 6",
			desc:            "testDescription 6",
			ownerID:         userID,
			questions: []game.Question{
				&game.TrueFalseQuestion{
					Title:            "title 6",
					TimeLimit:        30,
					TrueAlternative:  "testTrueAlternative 6",
					FalseAlternative: "testFalseAlternative 6",
				},
			},
		},
		{
			testDescription: "game with true or false question with invalid points",
			title:           "title 7",
			desc:            "testDescription 7",
			ownerID:         userID,
			questions: []game.Question{
				&game.TrueFalseQuestion{
					Title:            "title 7",
					Points:           3,
					TimeLimit:        30,
					TrueAlternative:  "testTrueAlternative 7",
					FalseAlternative: "testFalseAlternative 7",
				},
			},
		},
		{
			testDescription: "game with true or false question with blank true alternative",
			title:           "title 8",
			desc:            "testDescription 8",
			ownerID:         userID,
			questions: []game.Question{
				&game.TrueFalseQuestion{
					Title:            "title 8",
					Points:           1,
					TimeLimit:        30,
					TrueAlternative:  "",
					FalseAlternative: "testFalseAlternative 8",
				},
			},
		},
		{
			testDescription: "game with true or false question with blank false alternative",
			title:           "title 9",
			desc:            "testDescription 9",
			ownerID:         userID,
			questions: []game.Question{
				&game.TrueFalseQuestion{
					Title:            "title 9",
					Points:           1,
					TimeLimit:        30,
					TrueAlternative:  "testTrueAlternative 9",
					FalseAlternative: "",
				},
			},
		},
		{
			testDescription: "game with true or false question without false alternative",
			title:           "title 10",
			desc:            "testDescription 10",
			ownerID:         userID,
			questions: []game.Question{
				&game.TrueFalseQuestion{
					Title:           "title 10",
					Points:          1,
					TimeLimit:       30,
					TrueAlternative: "testTrueAlternative 10",
				},
			},
		},
		{
			testDescription: "game with true or false question without true alternative",
			title:           "title 11",
			desc:            "testDescription 11",
			ownerID:         userID,
			questions: []game.Question{
				&game.TrueFalseQuestion{
					Title:            "title 11",
					Points:           1,
					TimeLimit:        30,
					FalseAlternative: "false alternative 11",
				},
			},
		},
		{
			testDescription: "game with true or false question without time limit",
			title:           "title 12",
			desc:            "testDescription 12",
			ownerID:         userID,
			questions: []game.Question{
				&game.TrueFalseQuestion{
					Title:            "title 12",
					Points:           1,
					FalseAlternative: "false alternative 12",
					TrueAlternative:  "testTrueAlternative 12",
				},
			},
		},
		{
			testDescription: "game with true or false question with too low time limit",
			title:           "title 13",
			desc:            "testDescription 13",
			ownerID:         userID,
			questions: []game.Question{
				&game.TrueFalseQuestion{
					Title:            "title 13",
					Points:           1,
					TimeLimit:        1,
					FalseAlternative: "false alternative 13",
					TrueAlternative:  "testTrueAlternative 13",
				},
			},
		},
		{
			testDescription: "game with true or false question with too high time limit",
			title:           "title 14",
			desc:            "testDescription 14",
			ownerID:         userID,
			questions: []game.Question{
				&game.TrueFalseQuestion{
					Title:            "title 14",
					Points:           1,
					TimeLimit:        800,
					FalseAlternative: "false alternative 14",
					TrueAlternative:  "testTrueAlternative 14",
				},
			},
		},
		{
			testDescription: "game with quiz question without title",
			title:           "title 15",
			desc:            "testDescription 15",
			ownerID:         userID,
			questions: []game.Question{
				&game.QuizQuestion{
					Points:       1,
					TimeLimit:    30,
					Alternatives: *s.generateMockedQuizAlternatives(),
				},
			},
		},
		{
			testDescription: "game with quiz question with blank title",
			title:           "title 15",
			desc:            "testDescription 15",
			ownerID:         userID,
			questions: []game.Question{
				&game.QuizQuestion{
					Title:        "",
					Points:       1,
					TimeLimit:    30,
					Alternatives: *s.generateMockedQuizAlternatives(),
				},
			},
		},
		{
			testDescription: "game with quiz question with too low points",
			title:           "title 16",
			desc:            "testDescription 16",
			ownerID:         userID,
			questions: []game.Question{
				&game.QuizQuestion{
					Title:        "title 16",
					Points:       0,
					TimeLimit:    30,
					Alternatives: *s.generateMockedQuizAlternatives(),
				},
			},
		},
		{
			testDescription: "game with quiz question with too high points",
			title:           "title 17",
			desc:            "testDescription 17",
			ownerID:         userID,
			questions: []game.Question{
				&game.QuizQuestion{
					Title:        "title 17",
					Points:       5,
					TimeLimit:    30,
					Alternatives: *s.generateMockedQuizAlternatives(),
				},
			},
		},
		{
			testDescription: "game with quiz question with too low time limit",
			title:           "title 18",
			desc:            "testDescription 18",
			ownerID:         userID,
			questions: []game.Question{
				&game.QuizQuestion{
					Title:        "title 18",
					Points:       1,
					TimeLimit:    4,
					Alternatives: *s.generateMockedQuizAlternatives(),
				},
			},
		},
		{
			testDescription: "game with quiz question with too high time limit",
			title:           "title 19",
			desc:            "testDescription 19",
			ownerID:         userID,
			questions: []game.Question{
				&game.QuizQuestion{
					Title:        "title 19",
					Points:       1,
					TimeLimit:    300,
					Alternatives: *s.generateMockedQuizAlternatives(),
				},
			},
		},
		{
			testDescription: "game with quiz question with below minimum alternatives",
			title:           "title 20",
			desc:            "testDescription 20",
			ownerID:         userID,
			questions: []game.Question{
				&game.QuizQuestion{
					Title:     "title 20",
					Points:    1,
					TimeLimit: 30,
					Alternatives: []game.Alternative{
						{
							Data:      "alternativeFirst 20",
							IsCorrect: true,
						},
					},
				},
			},
		},
	}

	validatorError := &services.ValidationError{}
	for _, tt := range table {
		mockedGame := s.generateMockedGame(tt.title, tt.desc, tt.ownerID, tt.questions)

		t.Run(tt.testDescription, func(t *testing.T) {
			// Act
			err := s.svc.CreateNewGame(s.ctx, tt.ownerID, mockedGame)

			// Assert
			assert.ErrorContains(t, err, validatorError.Error())
		})
	}
}
