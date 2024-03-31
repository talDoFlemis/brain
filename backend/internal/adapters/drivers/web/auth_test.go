package web

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/taldoflemis/brain.test/internal/adapters/driven/auth"
	"github.com/taldoflemis/brain.test/internal/adapters/driven/postgres"
	"github.com/taldoflemis/brain.test/internal/core/services"
	testshelpers "github.com/taldoflemis/brain.test/test/helpers"
)

const (
	seed             = "SzceVsT4GdFOlrZn60XMgrFcvMNUMuuJ"
	authHeaderPrefix = "Bearer "
)

var (
	testUsername         = "gepeto"
	testEmail            = "gepeto@gmail.com"
	testPassword         = "mypassword"
	testHashedPassword   = "$2a$12$TSjLw2cqeD5bcjPUgOWaaew3xP88soPytNTnMi27vxcNMCDaLFkBa"
	testUserId           = uuid.New().String()
	accessMaxAgeInMin    = 1
	refreshMaxAgeInHours = 1
)

type AuthHandlerTestSuite struct {
	app         *fiber.App
	pgContainer *testshelpers.PostgresContainer
	suite.Suite
	ctx  context.Context
	pool *pgxpool.Pool
	svc  *services.AuthenticationService
}

func (suite *AuthHandlerTestSuite) SetupSuite() {
	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandlerMiddleware,
	})
	suite.ctx = context.Background()

	logger := testshelpers.NewDummyLogger(log.Writer())
	pgContainer, err := testshelpers.CreatePostgresContainer(suite.ctx)
	if err != nil {
		log.Fatal("tubias", err)
	}
	pool, err := postgres.NewPool(pgContainer.ConnStr)
	if err != nil {
		log.Fatal(err)
	}
	postgres.Migrate(pgContainer.ConnStr, "../../driven/postgres/migrations/")

	repository := postgres.NewLocalIDPPostgresStorer(pool)
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
		repository,
	)

	jwtMiddleware := NewJWTMiddleware(authManager)
	validationService := services.NewValidationService()
	authService := services.NewAuthenticationService(
		logger,
		authManager,
		validationService,
	)

	authHandler := NewAuthHandler(jwtMiddleware, authService, validationService)

	authHandler.RegisterRoutes(app)

	suite.app = app
	suite.pgContainer = pgContainer
	suite.pool = pool
	suite.svc = authService
}

func (suite *AuthHandlerTestSuite) SetupTest() {
	_, err := suite.pool.Exec(
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

func (suite *AuthHandlerTestSuite) TearDownTest() {
	_, err := suite.pool.Exec(suite.ctx, "TRUNCATE TABLE users")
	if err != nil {
		log.Fatalf("error truncating users table: %s", err)
	}
}

func (suite *AuthHandlerTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating postgres container: %s", err)
	}
}

func TestLoginTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("too slow for testing.Short")
	}
	suite.Run(t, new(AuthHandlerTestSuite))
}

func (suite *AuthHandlerTestSuite) TestRegisterUser() {
	// Arrange
	t := suite.T()
	username := "tubias"
	email := "tubias@gmail.com"
	password := "hashedpass"

	req := map[string]interface{}{
		"username": username,
		"email":    email,
		"password": password,
	}

	server := httptest.NewServer(adaptor.FiberApp(suite.app))
	e := httpexpect.Default(t, server.URL)

	// Act
	resp := e.POST("/auth/").WithJSON(req).Expect()

	// Assert
	resp.Status(http.StatusCreated)
	obj := resp.JSON().Object()
	obj.ContainsKey("access_token")
	obj.ContainsKey("refresh_token")
	obj.ContainsKey("expire_at")
}

func (suite *AuthHandlerTestSuite) TestRegisterUserWithBadInput() {
	// Arrange
	t := suite.T()
	server := httptest.NewServer(adaptor.FiberApp(suite.app))
	e := httpexpect.Default(t, server.URL)

	table := []struct {
		req         map[string]interface{}
		description string
		errorCount  int
	}{
		{
			req:         map[string]interface{}{},
			description: "empty request",
		},
		{
			req: map[string]interface{}{
				"username": "",
				"email":    testEmail,
				"password": testPassword,
			},
			description: "empty username",
		},
		{
			req: map[string]interface{}{
				"username": "temido",
				"email":    "",
				"password": testPassword,
			},
			description: "empty email",
		},
		{
			req: map[string]interface{}{
				"username": "gabrigas",
				"email":    testEmail,
				"password": "",
			},
			description: "empty password",
		},
		{
			req: map[string]interface{}{
				"username": "mx30",
				"email":    "invalidemail",
				"password": testPassword,
			},
			description: "invalid email",
		},
	}

	for _, tt := range table {
		t.Run(tt.description, func(t *testing.T) {
			// Act
			resp := e.POST("/auth/").WithJSON(tt.req).Expect()

			// Assert
			resp.Status(http.StatusUnprocessableEntity)
			obj := resp.JSON().Object()
			obj.ContainsKey("errors")
			obj.Value("errors").Array().NotEmpty()
		})
	}

}

func (suite *AuthHandlerTestSuite) TestRegisterUserThatAlreadyExists() {
	// Arrange
	t := suite.T()
	req := map[string]interface{}{
		"username": testUsername,
		"email":    testEmail,
		"password": testPassword,
	}

	server := httptest.NewServer(adaptor.FiberApp(suite.app))
	e := httpexpect.Default(t, server.URL)

	// Act
	resp := e.POST("/auth/").WithJSON(req).Expect()

	// Assert
	resp.Status(http.StatusConflict)
}

func (suite *AuthHandlerTestSuite) TestLogin() {
	// Arrange
	t := suite.T()
	req := map[string]interface{}{
		"username": testUsername,
		"password": testPassword,
	}

	server := httptest.NewServer(adaptor.FiberApp(suite.app))
	e := httpexpect.Default(t, server.URL)

	// Act
	resp := e.POST("/auth/login").WithJSON(req).Expect()

	// Assert
	resp.Status(http.StatusOK)
	obj := resp.JSON().Object()
	obj.ContainsKey("access_token")
	obj.ContainsKey("refresh_token")
	obj.ContainsKey("expire_at")
}

func (suite *AuthHandlerTestSuite) TestLoginWithUserThatDontExist() {
	// Arrange
	t := suite.T()
	req := map[string]interface{}{
		"username": "nonexistentuser",
		"password": testPassword,
	}

	server := httptest.NewServer(adaptor.FiberApp(suite.app))
	e := httpexpect.Default(t, server.URL)

	// Act
	resp := e.POST("/auth/login").WithJSON(req).Expect()

	// Assert
	resp.Status(http.StatusUnauthorized)
}

func (suite *AuthHandlerTestSuite) TestLoginWithWrongPassword() {
	// Arrange
	t := suite.T()
	req := map[string]interface{}{
		"username": testUsername,
		"password": "wrongpassword",
	}

	server := httptest.NewServer(adaptor.FiberApp(suite.app))
	e := httpexpect.Default(t, server.URL)

	// Act
	resp := e.POST("/auth/login").WithJSON(req).Expect()

	// Assert
	resp.Status(http.StatusUnauthorized)
}

func (suite *AuthHandlerTestSuite) TestRefreshToken() {
	// Arrange
	t := suite.T()
	tok, err := suite.svc.CreateToken(suite.ctx, testUserId)
	assert.NoError(t, err)
	assert.NotNil(t, tok)

	req := map[string]interface{}{
		"refresh_token": tok.RefreshToken,
	}

	server := httptest.NewServer(adaptor.FiberApp(suite.app))
	e := httpexpect.Default(t, server.URL)

	// Act
	resp := e.POST("/auth/refresh").WithJSON(req).Expect()

	// Assert
	resp.Status(http.StatusOK)
	obj := resp.JSON().Object()
	obj.ContainsKey("access_token")
	obj.ContainsKey("refresh_token")
	obj.ContainsKey("expire_at")
}

func (suite *AuthHandlerTestSuite) TestRefreshTokenWithInvalidRefreshToken() {
	// Arrange
	t := suite.T()
	req := map[string]interface{}{
		"refresh_token": "invalid",
	}

	table := []struct {
		req         map[string]interface{}
		description string
	}{
		{
			req: map[string]interface{}{
				"refresh_token": "invalid",
			},
			description: "invalid refresh token",
		},
		{
			req: map[string]interface{}{
				"refresh_token": "",
			},
			description: "empty refresh token",
		},
	}

	server := httptest.NewServer(adaptor.FiberApp(suite.app))
	e := httpexpect.Default(t, server.URL)

	for _, tt := range table {
		t.Run(tt.description, func(t *testing.T) {
			// Act
			resp := e.POST("/auth/refresh").WithJSON(req).Expect()

			// Assert
			resp.Status(http.StatusBadRequest)
		})
	}
}

func (suite *AuthHandlerTestSuite) TestUpdateUser() {
	// Arrange
	t := suite.T()
	tok, err := suite.svc.CreateToken(suite.ctx, testUserId)
	assert.NoError(t, err)

	req := map[string]interface{}{
		"username": "newusername",
		"email":    "geponto@gmail.com",
		"password": "newpassword",
	}

	headers := map[string]string{
		"Authorization": authHeaderPrefix + tok.AccessToken,
	}

	server := httptest.NewServer(adaptor.FiberApp(suite.app))
	e := httpexpect.Default(t, server.URL)

	// Act
	resp := e.PUT("/auth/").WithHeaders(headers).WithJSON(req).Expect()

	// Assert
	resp.Status(http.StatusNoContent)
}

func (suite *AuthHandlerTestSuite) TestUpdateUserWithInvalidInput() {
	// Arrange
	t := suite.T()
	tok, err := suite.svc.CreateToken(suite.ctx, testUserId)
	assert.NoError(t, err)

	headers := map[string]string{
		"Authorization": authHeaderPrefix + tok.AccessToken,
	}

	table := []struct {
		req         map[string]interface{}
		description string
	}{
		{
			req:         map[string]interface{}{},
			description: "blank input",
		},
		{
			req: map[string]interface{}{
				"username": "",
				"email":    testEmail,
				"password": testPassword,
			},
			description: "blank username",
		},
		{
			req: map[string]interface{}{
				"username": testUsername,
				"email":    "",
				"password": testPassword,
			},
			description: "empty email",
		},
		{
			req: map[string]interface{}{
				"username": testUsername,
				"email":    testEmail,
				"password": "",
			},
			description: "empty password",
		},
		{
			req: map[string]interface{}{
				"username": testUsername,
				"email":    "bademail",
				"password": testPassword,
			},
			description: "empty password",
		},
	}

	server := httptest.NewServer(adaptor.FiberApp(suite.app))
	e := httpexpect.Default(t, server.URL)

	for _, tt := range table {
		t.Run(tt.description, func(t *testing.T) {
			// Act
			resp := e.PUT("/auth/").WithHeaders(headers).WithJSON(tt.req).Expect()

			// Assert
			resp.Status(http.StatusUnprocessableEntity)
			resp.JSON().Object().Value("errors").Array().NotEmpty()
		})
	}
}

func (suite *AuthHandlerTestSuite) TestUpdateUserWithoutAuthorization() {
	// Arrange
	t := suite.T()

	req := map[string]interface{}{
		"username": "newusername",
		"email":    "email@gmail.com",
		"password": "newpassword",
	}
	headers := map[string]string{
		"Authorization": authHeaderPrefix + "invalid_token",
	}
	server := httptest.NewServer(adaptor.FiberApp(suite.app))
	e := httpexpect.Default(t, server.URL)

	// Act
	resp := e.PUT("/auth/").WithHeaders(headers).WithJSON(req).Expect()

	// Assert
	resp.Status(http.StatusUnauthorized)
}

func (suite *AuthHandlerTestSuite) TestUpdateUserToUsernameThatAlreadyExists() {
	// Arrange
	t := suite.T()
	otherUsername := "otherusername"
	_, err := suite.pool.Exec(
		suite.ctx,
		"INSERT INTO users (id, username, email, password) VALUES ($1, $2, $3, $4)",
		uuid.New().String(),
		otherUsername,
		testEmail,
		testHashedPassword,
	)
	assert.NoError(t, err)

	tok, err := suite.svc.CreateToken(suite.ctx, testUserId)
	assert.NoError(t, err)

	req := map[string]interface{}{
		"username": otherUsername,
		"email":    testEmail,
		"password": testPassword,
	}
	headers := map[string]string{
		"Authorization": authHeaderPrefix + tok.AccessToken,
	}

	server := httptest.NewServer(adaptor.FiberApp(suite.app))
	e := httpexpect.Default(t, server.URL)

	// Act
	resp := e.PUT("/auth/").WithHeaders(headers).WithJSON(req).Expect()

	// Assert
	resp.Status(http.StatusConflict)
}

func (suite *AuthHandlerTestSuite) TestDeleteAccount() {
	// Arrange
	t := suite.T()
	tok, err := suite.svc.CreateToken(suite.ctx, testUserId)
	assert.NoError(t, err)

	headers := map[string]string{
		"Authorization": authHeaderPrefix + tok.AccessToken,
	}

	server := httptest.NewServer(adaptor.FiberApp(suite.app))
	e := httpexpect.Default(t, server.URL)

	// Act
	resp := e.DELETE("/auth/").WithHeaders(headers).Expect()

	// Assert
	resp.Status(http.StatusNoContent)
}

func (suite *AuthHandlerTestSuite) TestDeleteAccountWithoutAuthorization() {
	// Arrange
	t := suite.T()

	headers := map[string]string{
		"Authorization": authHeaderPrefix + "invalid_token",
	}

	server := httptest.NewServer(adaptor.FiberApp(suite.app))
	e := httpexpect.Default(t, server.URL)

	// Act
	resp := e.DELETE("/auth/").WithHeaders(headers).Expect()

	// Assert
	resp.Status(http.StatusUnauthorized)
}

func (suite *AuthHandlerTestSuite) TestUserInfo() {
	// Arrange
	t := suite.T()
	tok, err := suite.svc.CreateToken(suite.ctx, testUserId)
	assert.NoError(t, err)

	headers := map[string]string{
		"Authorization": authHeaderPrefix + tok.AccessToken,
	}

	server := httptest.NewServer(adaptor.FiberApp(suite.app))
	e := httpexpect.Default(t, server.URL)

	// Act
	resp := e.GET("/auth/userinfo").WithHeaders(headers).Expect()

	// Assert
	resp.Status(http.StatusOK)
	json := resp.JSON().Object()
	json.IsEqual(map[string]interface{}{
		"username": testUsername,
		"email":    testEmail,
		"id":       testUserId,
	})
}

func (suite *AuthHandlerTestSuite) TestUserInfoWithBadInput() {
	// Arrange
	t := suite.T()
	invalidToken, err := suite.svc.CreateToken(suite.ctx, uuid.New().String())
	assert.NoError(t, err)

	table := []struct {
		headers     map[string]string
		status      int
		description string
	}{
		{
			headers:     make(map[string]string),
			status:      fiber.StatusUnauthorized,
			description: "empty token",
		},
		{
			headers: map[string]string{
				"Authorization": authHeaderPrefix + "invalid_token",
			},
			status:      fiber.StatusUnauthorized,
			description: "invalid token",
		},
		{
			headers: map[string]string{
				"Authorization": authHeaderPrefix + invalidToken.AccessToken,
			},
			status:      fiber.StatusUnauthorized,
			description: "invalid user id",
		},
	}

	server := httptest.NewServer(adaptor.FiberApp(suite.app))
	e := httpexpect.Default(t, server.URL)

	for _, tt := range table {
		t.Run(tt.description, func(t *testing.T) {
			// Act
			resp := e.GET("/auth/userinfo").
				WithHeaders(tt.headers).
				Expect()

			// Assert
			resp.Status(tt.status)
		})
	}
}
