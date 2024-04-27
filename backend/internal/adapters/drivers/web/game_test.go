package web_test

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	testcontainers "github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/taldoflemis/brain.test/internal/adapters/driven/auth"
	"github.com/taldoflemis/brain.test/internal/adapters/driven/postgres"
	"github.com/taldoflemis/brain.test/internal/adapters/drivers/web"
	"github.com/taldoflemis/brain.test/internal/core/services"
	"github.com/taldoflemis/brain.test/internal/ports"
	testshelpers "github.com/taldoflemis/brain.test/test/helpers"
)

const (
	route = "/game/"
)

func newJWTMiddleware(
	logger ports.Logger,
	pool *pgxpool.Pool,
) (fiber.Handler, ports.AuthenticationManager) {
	cfg := auth.NewLocalIdpConfig(
		seed,
		"issuer",
		"audience",
		accessMaxAgeInMin,
		refreshMaxAgeInHours,
	)
	authManager := auth.NewLocalIdp(
		*cfg,
		logger,
		postgres.NewLocalIDPPostgresStorer(pool),
	)
	return web.NewJWTMiddleware(authManager), authManager
}

type GameHandlerTestSuite struct {
	app *fiber.App
	suite.Suite
	pgContainer *testcontainers.PostgresContainer
	ctx         context.Context
	pool        *pgxpool.Pool
	svc         *services.GameService
	idp         ports.AuthenticationManager
}

func (suite *GameHandlerTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	logger := testshelpers.NewDummyLogger(log.Writer())
	pgContainer, pool, err := testshelpers.CreatePostgresContainerAndMigrate(
		suite.ctx,
		"../../driven/postgres/migrations/",
	)
	if err != nil {
		log.Fatal("failed to create container and migrate", err)
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: web.ErrorHandlerMiddleware,
	})

	gameStorer := postgres.NewPostgresGameStorer(pool)
	validationService := services.NewValidationService()
	gameService := services.NewGameService(logger, validationService, gameStorer)

	jwtMiddleware, idp := newJWTMiddleware(logger, pool)
	gameHandler := web.NewGameHandler(jwtMiddleware, validationService, gameService)
	gameHandler.RegisterRoutes(app)

	suite.app = app
	suite.pgContainer = pgContainer
	suite.pool = pool
	suite.svc = gameService
	suite.idp = idp

	_, err = suite.pool.Exec(
		suite.ctx,
		"INSERT INTO users (id, username, email, password) VALUES ($1, $2, $3, $4)",
		testUserId,
		testUsername,
		testEmail,
		testHashedPassword,
	)
	if err != nil {
		log.Fatalf("error inserting user: %s", err)
	}
}

func (suite *GameHandlerTestSuite) TearDownTest() {
	_, err := suite.pool.Exec(suite.ctx, "TRUNCATE TABLE games CASCADE")
	if err != nil {
		log.Fatalf("error truncating games table: %s", err)
	}
}

func (suite *GameHandlerTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating postgres container: %s", err)
	}
}

func TestGameHandlerTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("too slow for testing.Short")
	}
	suite.Run(t, new(GameHandlerTestSuite))
}

func (s *GameHandlerTestSuite) TestCreateGame() {
	// Arrange
	t := s.T()
	tok, err := s.idp.CreateToken(s.ctx, testUserId)
	assert.NoError(t, err)

	table := []struct {
		desc string
		req  map[string]interface{}
	}{
		{
			desc: "create game with true false question",
			req: map[string]any{
				"title":       "title",
				"description": "description",
				"questions": []map[string]any{
					{
						"kind": "true_false",
						"data": map[string]any{
							"title":             "title true false",
							"points":            1,
							"time_limit":        30,
							"true_alternative":  "true here",
							"false_alternative": "false here",
						},
					},
				},
			},
		},
		{
			desc: "create game with quiz question",
			req: map[string]any{
				"title":       "title",
				"description": "description",
				"questions": []map[string]any{
					{
						"kind": "quiz",
						"data": map[string]any{
							"title":      "title quiz",
							"points":     1,
							"time_limit": 30,
							"alternatives": []map[string]any{
								{
									"data":       "alternative 1",
									"is_correct": true,
								},
								{
									"data":       "alternative 2",
									"is_correct": true,
								},
								{
									"data":       "alternative 3",
									"is_correct": true,
								},
							},
						},
					},
				},
			},
		},
	}

	headers := map[string]string{
		"Authorization": authHeaderPrefix + tok.AccessToken,
	}

	server := httptest.NewServer(adaptor.FiberApp(s.app))
	e := httpexpect.Default(t, server.URL)

	for _, tt := range table {
		t.Run(tt.desc, func(t *testing.T) {
			// Act
			resp := e.POST(route).WithHeaders(headers).WithJSON(tt.req).Expect()

			// Assert
			resp.Status(http.StatusCreated)
		})
	}

}

func (s *GameHandlerTestSuite) TestCreateGameWithUnknownQuestionKind() {
	// Arrange
	t := s.T()
	tok, err := s.idp.CreateToken(s.ctx, testUserId)
	assert.NoError(t, err)
	req := map[string]any{
		"title":       "wrong kind  title",
		"description": "wrong kind  description",
		"questions": []map[string]any{
			{
				"kind": "tubias",
				"data": map[string]any{
					"title":             "wrong kind title true false",
					"points":            1,
					"time_limit":        30,
					"true_alternative":  "wrong kind true here",
					"false_alternative": "wrong kind false here",
				},
			},
		},
	}

	headers := map[string]string{
		"Authorization": authHeaderPrefix + tok.AccessToken,
	}

	server := httptest.NewServer(adaptor.FiberApp(s.app))
	e := httpexpect.Default(t, server.URL)

	// Act
	resp := e.POST(route).WithHeaders(headers).WithJSON(req).Expect()

	// Assert
	resp.Status(http.StatusBadRequest)
}

func (s *GameHandlerTestSuite) TestCreateGameWithInvalidRequest() {
	// Arrange
	t := s.T()
	tok, err := s.idp.CreateToken(s.ctx, testUserId)
	assert.NoError(t, err)

	table := []struct {
		desc string
		req  map[string]interface{}
	}{
		{
			desc: "missing title",
			req: map[string]interface{}{
				"description": "missing tittle  description",
				"questions": []map[string]any{
					{
						"kind": "true_false",
						"data": map[string]any{
							"title":             "missing tittle title true false",
							"points":            1,
							"time_limit":        30,
							"true_alternative":  "missing tittle true here",
							"false_alternative": "missing tittle false here",
						},
					},
				},
			},
		},
		{
			desc: "missing questions",
			req: map[string]interface{}{
				"title":       "missing description  title",
				"description": "missing description  description",
			},
		},
	}

	headers := map[string]string{
		"Authorization": authHeaderPrefix + tok.AccessToken,
	}

	server := httptest.NewServer(adaptor.FiberApp(s.app))
	e := httpexpect.Default(t, server.URL)

	for _, tt := range table {
		t.Run(tt.desc, func(t *testing.T) {
			// Act
			resp := e.POST(route).WithHeaders(headers).WithJSON(tt.req).Expect()

			// Assert
			resp.Status(http.StatusUnprocessableEntity)
		})
	}
}
