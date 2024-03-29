package auth

import (
	"context"
	"crypto"
	"crypto/ed25519"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/taldoflemis/brain.test/internal/ports"
)

const (
	defaultBCryptCost = 12
)

type LocalIdpConfig struct {
	privateKey         ed25519.PrivateKey
	publicKey          crypto.PublicKey
	issuer             string
	audience           string
	accessTokenMaxAge  time.Duration
	refreshTokenMaxAge time.Duration
}

func NewLocalIdpConfig(
	seed, issuer, audience string,
	accessTimeInMin, refreshTimeInHours int,
) *LocalIdpConfig {
	privateKey := ed25519.NewKeyFromSeed([]byte(seed))

	return &LocalIdpConfig{
		privateKey:         privateKey,
		publicKey:          privateKey.Public(),
		issuer:             issuer,
		audience:           audience,
		accessTokenMaxAge:  time.Duration(accessTimeInMin) * time.Minute,
		refreshTokenMaxAge: time.Duration(refreshTimeInHours) * time.Hour,
	}
}

type localIDP struct {
	cfg    LocalIdpConfig
	logger ports.Logger
	repo   ports.LocalIDPStorer
}

func NewLocalIdp(cfg LocalIdpConfig, logger ports.Logger, repo ports.LocalIDPStorer) *localIDP {
	return &localIDP{
		cfg:    cfg,
		logger: logger,
		repo:   repo,
	}
}

func (i *localIDP) CreateUser(
	ctx context.Context,
	username, email, password string,
) (*ports.UserIdentityInfo, error) {
	i.logger.Debug("Creating user", "username", username, "email", email)

	hashedPassword, err := i.hashPassword(password)
	if err != nil {
		i.logger.Error("Failed to hash password", err)
		return nil, err
	}

	user, err := i.repo.StoreUser(ctx, username, email, hashedPassword)
	if err != nil {
		i.logger.Error("Failed to store user", err)
		return nil, err
	}

	i.logger.Info("User created", "username", username, "email", email)

	return &ports.UserIdentityInfo{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
	}, nil
}

func (i *localIDP) CreateToken(
	ctx context.Context,
	userId string,
) (*ports.TokenResponse, error) {
	i.logger.Debug("Creating token", "userId", userId)
	accessToken, err := i.generateToken(ctx, userId, i.cfg.accessTokenMaxAge)
	if err != nil {
		i.logger.Error("Failed to sign access token", err)
		return nil, ports.ErrFailedToSignToken
	}

	refreshToken, err := i.generateToken(ctx, userId, i.cfg.refreshTokenMaxAge)
	if err != nil {
		i.logger.Error("Failed to sign refresh token", err)
		return nil, ports.ErrFailedToSignToken
	}

	i.logger.Info("Token created", "userId", userId)
	return &ports.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(i.cfg.accessTokenMaxAge),
	}, nil
}

func (i *localIDP) RefreshToken(
	ctx context.Context,
	refreshToken string,
) (*ports.TokenResponse, error) {
	i.logger.Debug("Refreshing token")
	token, err := jwt.ParseWithClaims(
		refreshToken,
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return i.cfg.publicKey, nil
		},
		jwt.WithAudience(i.cfg.audience),
		jwt.WithIssuer(i.cfg.issuer),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		i.logger.Error("Failed to parse refresh token", err)
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ports.ErrExpiredToken
		}
		return nil, ports.ErrInvalidRefreshToken
	}
	claims := token.Claims.(*jwt.RegisteredClaims)

	accessToken, err := i.generateToken(ctx, claims.Subject, i.cfg.accessTokenMaxAge)
	if err != nil {
		return nil, err
	}

	i.logger.Info("Token refreshed", "userId", claims.Subject)

	return &ports.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(i.cfg.accessTokenMaxAge),
	}, nil
}

func (i *localIDP) AuthenticateUser(
	ctx context.Context,
	username string,
	password string,
) (*ports.UserIdentityInfo, error) {
	user, err := i.repo.FindUserByUsername(ctx, username)
	if err != nil {
		i.logger.Error("Failed to find user", err)
		return nil, err
	}

	if i.comparePassword(user.HashedPassword, password) != nil {
		i.logger.Error("Invalid password")
		return nil, ports.ErrInvalidPassword
	}

	return &ports.UserIdentityInfo{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
	}, nil
}

func (i *localIDP) DeleteUser(ctx context.Context, userId string) error {
	i.logger.Debug("Deleting user", "userId", userId)

	err := i.repo.DeleteUser(ctx, userId)
	if err != nil {
		i.logger.Error("Failed to delete user", err)
		return err
	}

	i.logger.Info("User deleted", "userId", userId)
	return nil
}

func (i *localIDP) UpdateUser(
	ctx context.Context,
	userId,
	username,
	password,
	email string,
) (*ports.UserIdentityInfo, error) {
	i.logger.Debug("Updating user", "userId", userId)

	hashedPassword, err := i.hashPassword(password)
	if err != nil {
		i.logger.Error("Failed to hash password", err)
		return nil, err
	}

	user, err := i.repo.UpdateUser(ctx, userId, username, hashedPassword, email)
	if err != nil {
		i.logger.Error("Failed to update user", err)
		return nil, err
	}

	i.logger.Info("User updated", "userId", userId)
	return &ports.UserIdentityInfo{
		Email:    user.Email,
		Username: user.Username,
		ID:       userId,
	}, nil
}

func (i *localIDP) GetPublicKey() interface{} {
	return i.cfg.publicKey
}

func (i *localIDP) GetAlgorithm() string {
	return jwt.SigningMethodEdDSA.Alg()
}

func (i *localIDP) generateToken(
	ctx context.Context,
	userId string,
	expireDate time.Duration,
) (string, error) {
	accessTokenClaims := jwt.RegisteredClaims{
		Subject:   userId,
		Issuer:    i.cfg.issuer,
		Audience:  []string{i.cfg.audience},
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireDate)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, accessTokenClaims)
	accessToken, err := token.SignedString(i.cfg.privateKey)
	if err != nil {
		return "", ports.ErrFailedToSignToken
	}
	return accessToken, nil
}

func (i *localIDP) comparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (i *localIDP) hashPassword(password string) (string, error) {
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(password), defaultBCryptCost)
	if err != nil {
		return "", err
	}
	return string(hashedPasswordBytes), nil
}
