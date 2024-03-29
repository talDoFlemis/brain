package services

import (
	"context"

	"github.com/taldoflemis/brain.test/internal/core/domain"
	"github.com/taldoflemis/brain.test/internal/ports"
)

type AuthenticationService struct {
	logger            ports.Logger
	authManager       ports.AuthenticationManager
	validationService *ValidationService
}

func NewAuthenticationService(
	logger ports.Logger,
	authManager ports.AuthenticationManager,
	validationService *ValidationService,
) *AuthenticationService {
	return &AuthenticationService{
		logger:            logger,
		authManager:       authManager,
		validationService: validationService,
	}
}

func (s *AuthenticationService) CreateUser(
	ctx context.Context,
	req *domain.CreateUserRequest,
) (*ports.UserIdentityInfo, error) {
	err := s.validationService.Validate(req)
	if err != nil {
		s.logger.Error("failed to validate struct")
		return nil, err
	}

	return s.authManager.CreateUser(ctx, req.Username, req.Email, req.Password)
}

func (s *AuthenticationService) AuthenticateUser(
	ctx context.Context,
	username, password string,
) (*ports.UserIdentityInfo, error) {
	return s.authManager.AuthenticateUser(ctx, username, password)
}

func (s *AuthenticationService) DeleteUser(ctx context.Context, userId string) error {
	return s.authManager.DeleteUser(ctx, userId)
}

func (s *AuthenticationService) UpdateUser(
	ctx context.Context,
	userId string,
	req *domain.UpdateUserRequest,
) (*ports.UserIdentityInfo, error) {
	err := s.validationService.Validate(req)
	if err != nil {
		s.logger.Error("failed to validate struct")
		return nil, err
	}

	return s.authManager.UpdateUser(ctx, userId, req.Username, req.Password, req.Email)
}

func (s *AuthenticationService) CreateToken(
	ctx context.Context,
	userId string,
) (*ports.TokenResponse, error) {
	return s.authManager.CreateToken(ctx, userId)
}

func (s *AuthenticationService) RefreshToken(
	ctx context.Context,
	refreshToken string,
) (*ports.TokenResponse, error) {
	return s.authManager.RefreshToken(ctx, refreshToken)
}

func (s *AuthenticationService) GetPublicKey() interface{} {
	return s.authManager.GetPublicKey()
}

func (s *AuthenticationService) GetAlgorithm() string {
	return s.authManager.GetAlgorithm()
}
