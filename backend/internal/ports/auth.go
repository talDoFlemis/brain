package ports

import (
	"context"
	"errors"
	"time"
)

var (
	ErrFailedToSignToken   = errors.New("token sign failed")
	ErrUserNotFound        = errors.New("user not found")
	ErrInvalidPassword     = errors.New("invalid password")
	ErrExpiredToken        = errors.New("token expired")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrUserAlreadyExists   = errors.New("user already exists")
)

type TokenResponse struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

type UserIdentityInfo struct {
	ID       string
	Username string
	Email    string
}

type AuthenticationManager interface {
	CreateUser(ctx context.Context, username, email, password string) (*UserIdentityInfo, error)
	AuthenticateUser(ctx context.Context, username, password string) (*UserIdentityInfo, error)
	DeleteUser(ctx context.Context, userId string) error
	UpdateUser(
		ctx context.Context,
		userId, username, password, email string,
	) (*UserIdentityInfo, error)
	CreateToken(ctx context.Context, userId string) (*TokenResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error)
	GetPublicKey() interface{}
	GetAlgorithm() string
}

type LocalIDPUserEntity struct {
	ID             string
	Username       string
	Email          string
	HashedPassword string
}

type LocalIDPStorer interface {
	StoreUser(
		ctx context.Context,
		username, email, password string,
	) (*LocalIDPUserEntity, error)
	UpdateUser(
		ctx context.Context,
		userId, username, password, email string,
	) (*LocalIDPUserEntity, error)
	DeleteUser(ctx context.Context, userId string) error
	FindUserByUsername(ctx context.Context, username string) (*LocalIDPUserEntity, error)
}
