package services

import (
	"context"

	"github.com/taldoflemis/brain.test/internal/ports"
)

type CreateUserRequest struct {
	Username string `validate:"required"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8,max=72"`
}

type UpdateUserRequest struct {
	Username string `validate:"required"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8,max=72"`
}

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
	req *CreateUserRequest,
) (*ports.TokenResponse, error) {
	err := s.validationService.Validate(req)
	if err != nil {
		s.logger.Error("failed to validate struct")
		return nil, err
	}

	info, err := s.authManager.CreateUser(ctx, req.Username, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return s.authManager.CreateToken(ctx, info.ID)
}

func (s *AuthenticationService) AuthenticateUser(
	ctx context.Context,
	username, password string,
) (*ports.TokenResponse, error) {
	info, err := s.authManager.AuthenticateUser(ctx, username, password)
	if err != nil {
		return nil, err
	}

	return s.authManager.CreateToken(ctx, info.ID)
}

func (s *AuthenticationService) DeleteUser(ctx context.Context, userId string) error {
	return s.authManager.DeleteUser(ctx, userId)
}

func (s *AuthenticationService) UpdateUser(
	ctx context.Context,
	userId string,
	req *UpdateUserRequest,
) (*ports.UserIdentityInfo, error) {
	err := s.validationService.Validate(req)
	if err != nil {
		s.logger.Error("failed to validate struct")
		return nil, err
	}

	return s.authManager.UpdateUser(ctx, userId, req.Username, req.Password, req.Email)
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

func (s *AuthenticationService) GetUserInfo(
	ctx context.Context,
	token string,
) (*ports.UserIdentityInfo, error) {
	return s.authManager.GetUserInfo(ctx, token)
}
