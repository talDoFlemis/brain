package postgres

import (
	"context"
	"log"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/taldoflemis/brain.test/internal/ports"
	"github.com/taldoflemis/brain.test/test/helpers"
)

var (
	testUsername = "tubias"
	testEmail    = "tubias3@gmail.com"
	testPassword = "hashedpassword"
	testUserId   = "f7396104-a636-4826-9d9f-b92ae90cea14"
)

type LocalIDPPostgresStorerTestSuite struct {
	suite.Suite
	pgContainer *testshelpers.PostgresContainer
	ctx         context.Context
	repo        *LocalIDPPostgresStorer
	pool        *pgxpool.Pool
}

func (suite *LocalIDPPostgresStorerTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	pgContainer, err := testshelpers.CreatePostgresContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}
	suite.pgContainer = pgContainer
	pool, err := NewPool(pgContainer.ConnStr)
	if err != nil {
		log.Fatal(err)
	}

	Migrate(pgContainer.ConnStr, "./migrations/")

	repository := NewLocalIDPPostgresStorer(pool)

	suite.repo = repository
	suite.pool = pool
}

func (suite *LocalIDPPostgresStorerTestSuite) SetupTest() {
	_, err := suite.pool.Exec(
		suite.ctx,
		`INSERT INTO users (id, username, email, password) VALUES ($1, $2, $3, $4)`,
		testUserId, testUsername, testEmail, testPassword,
	)
	if err != nil {
		log.Fatalf("error inserting user: %s", err)
	}
}

func (suite *LocalIDPPostgresStorerTestSuite) TearDownTest() {
	_, err := suite.pool.Exec(suite.ctx, "TRUNCATE TABLE users")
	if err != nil {
		log.Fatalf("error truncating users table: %s", err)
	}
}

func (suite *LocalIDPPostgresStorerTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating postgres container: %s", err)
	}
}

func TestLocalIDPPostgresStorer(t *testing.T) {
	if testing.Short() {
		t.Skip("too slow for testing.Short")
	}

	suite.Run(t, new(LocalIDPPostgresStorerTestSuite))
}

func (suite *LocalIDPPostgresStorerTestSuite) TestFindUserByUsername() {
	// Arrange
	t := suite.T()

	// Act
	user, err := suite.repo.FindUserByUsername(suite.ctx, testUsername)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, testUsername, user.Username)
	assert.Equal(t, testEmail, user.Email)
}

func (suite *LocalIDPPostgresStorerTestSuite) TestTryToFindUserByUsernameThatDoesNotExist() {
	// Arrange
	t := suite.T()
	random := "7ebc4755-b7cc-4963-a2b1-636949b035d6"

	// Act
	user, err := suite.repo.FindUserByUsername(suite.ctx, random)

	// Assert
	assert.Nil(t, user)
	assert.ErrorIs(t, err, ports.ErrUserNotFound)
}

func (suite *LocalIDPPostgresStorerTestSuite) TestFindUserById() {
	// Arrange
	t := suite.T()

	// Act
	user, err := suite.repo.FindUserById(suite.ctx, testUserId)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, testUsername, user.Username)
	assert.Equal(t, testEmail, user.Email)
	assert.Equal(t, testUserId, user.ID)
}

func (suite *LocalIDPPostgresStorerTestSuite) TestTryToFindUserByIdThatDoesNotExist() {
	// Arrange
	t := suite.T()

	// Act
	user, err := suite.repo.FindUserById(suite.ctx, "7ebc4755-b7cc-4963-a2b1-636949b035d6")

	// Assert
	assert.Nil(t, user)
	assert.ErrorIs(t, err, ports.ErrUserNotFound)
}

func (suite *LocalIDPPostgresStorerTestSuite) TestDeleteUser() {
	// Arrange
	t := suite.T()

	// Act
	err := suite.repo.DeleteUser(suite.ctx, testUserId)

	// Assert
	assert.NoError(t, err)
}

func (suite *LocalIDPPostgresStorerTestSuite) TestTryToDeleteUserThatDoesNotExist() {
	// Arrange
	t := suite.T()
	randomId := "d0b8b515-f46b-4179-bb26-f7833ded8f8f"

	// Act
	err := suite.repo.DeleteUser(suite.ctx, randomId)

	// Assert
	assert.ErrorIs(t, err, ports.ErrUserNotFound)
}

func (suite *LocalIDPPostgresStorerTestSuite) TestCreateUser() {
	// Arrange
	t := suite.T()
	username := "tubias2"
	email := "tubias2@gmail.com"
	password := "hashedpassword"

	// Act
	user, err := suite.repo.StoreUser(suite.ctx, username, email, password)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, username, user.Username)
	assert.Equal(t, email, user.Email)
}

func (suite *LocalIDPPostgresStorerTestSuite) TestTryToCreateUserWithExistingUsername() {
	// Arrange
	t := suite.T()
	password := "newpassword"

	// Act
	user, err := suite.repo.StoreUser(suite.ctx, testUsername, testEmail, password)

	// Assert
	assert.ErrorIs(t, err, ports.ErrUserAlreadyExists)
	assert.Nil(t, user)
}

func (suite *LocalIDPPostgresStorerTestSuite) TestUpdateUser() {
	// Arrange
	t := suite.T()
	username := "newtubias"
	email := "newtubias3@gmail.com"
	password := "hashedpassword"

	// Act
	user, err := suite.repo.UpdateUser(suite.ctx, testUserId, username, password, email)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, testUserId, user.ID)
	assert.Equal(t, username, user.Username)
	assert.Equal(t, email, user.Email)
}

func (suite *LocalIDPPostgresStorerTestSuite) TestUpdateUserWithExistingUsername() {
	// Arrange
	t := suite.T()
	newPassword := "newpassword"
	otherUsername := "otheruser"
	_, err := suite.repo.StoreUser(suite.ctx, otherUsername, testEmail, testPassword)
	assert.NoError(t, err)

	// Act
	user, err := suite.repo.UpdateUser(suite.ctx, testUserId, otherUsername, testEmail, newPassword)

	// Assert
	assert.ErrorIs(t, err, ports.ErrUserAlreadyExists)
	assert.Nil(t, user)
}

func (suite *LocalIDPPostgresStorerTestSuite) TestTryToUpdateUserThatDoesNotExist() {
	// Arrange
	t := suite.T()
	username := "tubias"
	email := "tubias@gmail.com"
	password := "hashedpassword"
	randomId := "d0b8b515-f46b-4179-bb26-f7833ded8f8f"

	// Act
	user, err := suite.repo.UpdateUser(suite.ctx, randomId, username, password, email)

	// Assert
	assert.ErrorIs(t, err, ports.ErrUserNotFound)
	assert.Nil(t, user)
}
