package services_test

import (
	"context"
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/taldoflemis/brain.test/internal/adapters/driven/auth"
	"github.com/taldoflemis/brain.test/internal/adapters/driven/postgres"
	"github.com/taldoflemis/brain.test/internal/core/services"
	"github.com/taldoflemis/brain.test/internal/ports"
	testshelpers "github.com/taldoflemis/brain.test/test/helpers"
)

const (
	seed = "SzceVsT4GdFOlrZn60XMgrFcvMNUMuuJ"
)

var (
	testUsername         = "gepeto"
	testEmail            = "gepeto@gmail.com"
	testPassword         = "mypassword"
	testHashedPassword   = "$2a$12$TSjLw2cqeD5bcjPUgOWaaew3xP88soPytNTnMi27vxcNMCDaLFkBa"
	testUserId           = uuid.New().String()
	accessMaxAgeInMin    = 15
	refreshMaxAgeInHours = 24
)

type AuthenticationServiceIntegrationTestSuite struct {
	suite.Suite
	pgContainer *testshelpers.PostgresContainer
	ctx         context.Context
	pool        *pgxpool.Pool
	svc         *services.AuthenticationService
	idp         ports.AuthenticationManager
}

func (suite *AuthenticationServiceIntegrationTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	pgContainer, err := testshelpers.CreatePostgresContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}
	suite.pgContainer = pgContainer
	pool, err := postgres.NewPool(pgContainer.ConnStr)
	if err != nil {
		log.Fatal(err)
	}

	postgres.Migrate(pgContainer.ConnStr, "../../adapters/driven/postgres/migrations/")

	repository := postgres.NewLocalIDPPostgresStorer(pool)
	logger := testshelpers.NewDummyLogger(log.Writer())
	cfg := auth.NewLocalIdpConfig(
		seed,
		"issuer",
		"audience",
		accessMaxAgeInMin,
		refreshMaxAgeInHours,
	)
	adapter := auth.NewLocalIdp(*cfg, logger, repository)
	svc := services.NewAuthenticationService(logger, adapter, services.NewValidationService())

	suite.svc = svc
	suite.pool = pool
	suite.idp = adapter
}

func (suite *AuthenticationServiceIntegrationTestSuite) SetupTest() {
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

func (suite *AuthenticationServiceIntegrationTestSuite) TearDownTest() {
	_, err := suite.pool.Exec(suite.ctx, "TRUNCATE TABLE users")
	if err != nil {
		log.Fatalf("error truncating users table: %s", err)
	}
}

func (suite *AuthenticationServiceIntegrationTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating postgres container: %s", err)
	}
}

func TestAuthenticationServiceIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("too slow for testing.Short")
	}

	suite.Run(t, new(AuthenticationServiceIntegrationTestSuite))
}

func (suite *AuthenticationServiceIntegrationTestSuite) TestCreateUser() {
	// Arrange
	t := suite.T()

	req := &services.CreateUserRequest{
		Email:    "newemail@gmail.com",
		Password: "newpassword",
		Username: "newusername",
	}

	// Act
	user, err := suite.svc.CreateUser(suite.ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
}

func (suite *AuthenticationServiceIntegrationTestSuite) TestCreateUserWithExistingUsername() {
	// Arrange
	t := suite.T()

	req := &services.CreateUserRequest{
		Username: testUsername,
		Email:    "newmail@gmail.com",
		Password: "newpassword",
	}

	// Act
	user, err := suite.svc.CreateUser(suite.ctx, req)

	// Assert
	assert.ErrorIs(t, err, ports.ErrUserAlreadyExists)
	assert.Nil(t, user)
}

func (suite *AuthenticationServiceIntegrationTestSuite) TestCreateUserWithBadInput() {
	// Arrange
	t := suite.T()

	req := &services.CreateUserRequest{
		Password: "",
		Email:    "",
		Username: testUsername,
	}

	table := []struct {
		badPassword string
		badEmail    string
		description string
	}{
		{
			badEmail:    "bademail",
			badPassword: testPassword,
			description: "email: must be a valid email address",
		},
		{
			badEmail:    "",
			badPassword: testPassword,
			description: "email: cannot be blank",
		},
		{
			badPassword: "123",
			badEmail:    testEmail,
			description: "password: the length must be greater or equal than 8",
		},
		{
			badPassword: "",
			badEmail:    testEmail,
			description: "password: cannot be blank",
		},
		{
			badPassword: "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890asd",
			badEmail:    testEmail,
			description: "password: the length must be less or equal than 72",
		},
	}

	validatorError := &services.ValidationError{}

	for _, tt := range table {
		req.Password = tt.badPassword
		req.Email = tt.badEmail
		t.Run(tt.description, func(t *testing.T) {
			// Act
			user, err := suite.svc.CreateUser(suite.ctx, req)

			// Assert
			assert.ErrorContains(t, err, validatorError.Error())
			assert.Nil(t, user)
		})
	}
}

func (suite *AuthenticationServiceIntegrationTestSuite) TestCreateToken() {
	// Arrange
	t := suite.T()

	// Act
	tokenResponse, err := suite.idp.CreateToken(suite.ctx, testUserId)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, tokenResponse)
}

func (suite *AuthenticationServiceIntegrationTestSuite) TestAuthenticateUser() {
	// Arrange
	t := suite.T()

	// Act
	user, err := suite.svc.AuthenticateUser(suite.ctx, testUsername, testPassword)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
}

func (suite *AuthenticationServiceIntegrationTestSuite) TestAuthenticateUserWithBadInput() {
	// Arrange
	t := suite.T()

	table := []struct {
		username    string
		password    string
		errorString string
		description string
	}{
		{
			username:    "nonexistentuser",
			password:    testPassword,
			errorString: ports.ErrUserNotFound.Error(),
			description: "nonexistent user",
		},
		{
			username:    testUsername,
			password:    "wrongpassword",
			errorString: ports.ErrInvalidPassword.Error(),
			description: "wrong password",
		},
	}

	for _, tt := range table {
		t.Run(tt.description, func(t *testing.T) {
			// Act
			user, err := suite.svc.AuthenticateUser(suite.ctx, tt.username, tt.password)

			// Assert
			assert.ErrorContains(t, err, tt.errorString)
			assert.Nil(t, user)
		})
	}
}

func (suite *AuthenticationServiceIntegrationTestSuite) TestDeleteUser() {
	// Arrange
	t := suite.T()

	// Act
	err := suite.svc.DeleteUser(suite.ctx, testUserId)

	// Assert
	assert.NoError(t, err)
}

func (suite *AuthenticationServiceIntegrationTestSuite) TestUpdateUser() {
	// Arrange
	t := suite.T()
	newEmail := "newemail@gmail.com"

	req := &services.UpdateUserRequest{
		Email:    newEmail,
		Password: testPassword,
		Username: testUsername,
	}

	// Act
	user, err := suite.svc.UpdateUser(suite.ctx, testUserId, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, testUserId, user.ID)
	assert.Equal(t, newEmail, user.Email)
	assert.Equal(t, testUsername, user.Username)
}

func (suite *AuthenticationServiceIntegrationTestSuite) TestUpdateUserWithBadInput() {
	// Arrange
	t := suite.T()

	req := &services.UpdateUserRequest{
		Password: "",
		Email:    "",
		Username: testUsername,
	}

	table := []struct {
		badPassword string
		badEmail    string
		description string
	}{
		{
			badEmail:    "bademail",
			badPassword: testPassword,
			description: "email: must be a valid email address",
		},
		{
			badEmail:    "",
			badPassword: testPassword,
			description: "email: cannot be blank",
		},
		{
			badPassword: "123",
			badEmail:    testEmail,
			description: "password: the length must be greater or equal than 8",
		},
		{
			badPassword: "",
			badEmail:    testEmail,
			description: "password: cannot be blank",
		},
		{
			badPassword: "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890asd",
			badEmail:    testEmail,
			description: "password: the length must be less or equal than 72",
		},
	}

	validatorError := &services.ValidationError{}

	for _, tt := range table {
		req.Password = tt.badPassword
		req.Email = tt.badEmail
		t.Run(tt.description, func(t *testing.T) {
			// Act
			user, err := suite.svc.UpdateUser(suite.ctx, testUserId, req)

			// Assert
			assert.ErrorContains(t, err, validatorError.Error())
			assert.Nil(t, user)
		})
	}
}

func (suite *AuthenticationServiceIntegrationTestSuite) TestRefreshToken() {
	// Arrange
	t := suite.T()
	tok, err := suite.idp.CreateToken(suite.ctx, testUserId)
	assert.NoError(t, err)

	// Act
	tok, err = suite.svc.RefreshToken(suite.ctx, tok.RefreshToken)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, tok)
}

func (suite *AuthenticationServiceIntegrationTestSuite) TestRefreshTokenWithBadInput() {
	// Arrange
	t := suite.T()

	table := []struct {
		refreshToken string
		description  string
	}{
		{
			refreshToken: "invalidtoken",
			description:  "invalid token",
		},
		{
			refreshToken: "",
			description:  "empty token",
		},
		{
			refreshToken: "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE3MTE3Mzk2ODEsImV4cCI6MTc0MzI3NTY4MSwiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoianJvY2tldEBleGFtcGxlLmNvbSJ9.W07mltoo-kXL4yHNqDhwoyueGI4HWwQeHdxH-m7znnU",
			description:  "token with invalid signature",
		},
	}

	for _, tt := range table {
		t.Run(tt.description, func(t *testing.T) {
			// Act
			tok, err := suite.svc.RefreshToken(suite.ctx, tt.refreshToken)

			// Assert
			assert.ErrorContains(t, err, ports.ErrInvalidRefreshToken.Error())
			assert.Nil(t, tok)
		})
	}
}

func (suite *AuthenticationServiceIntegrationTestSuite) TestGetUserInfo() {
	// Arrange
	t := suite.T()
	toks, err := suite.idp.CreateToken(suite.ctx, testUserId)
	assert.NoError(t, err)
	expectedUser := &ports.UserIdentityInfo{
		ID:       testUserId,
		Email:    testEmail,
		Username: testUsername,
	}

	// Act
	info, err := suite.svc.GetUserInfo(suite.ctx, toks.AccessToken)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, info)
}

func (suite *AuthenticationServiceIntegrationTestSuite) TestGetUserInfoWithUnknownUser() {
	// Arrange
	t := suite.T()
	unexistentUserId := uuid.New().String()
	toks, err := suite.idp.CreateToken(suite.ctx, unexistentUserId)
	assert.NoError(t, err)

	// Act
	info, err := suite.svc.GetUserInfo(suite.ctx, toks.AccessToken)

	// Assert
	assert.ErrorIs(t, err, ports.ErrUserNotFound)
	assert.Nil(t, info)
}
