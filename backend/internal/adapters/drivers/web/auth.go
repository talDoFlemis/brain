package web

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"github.com/taldoflemis/brain.test/internal/core/services"
	"github.com/taldoflemis/brain.test/internal/ports"
)

// LoginRequest
//
//	@Description	Request of Login
type LoginRequest struct {
	// the username of the user
	Username string `json:"username" validate:"required"`
	// the password of the user
	Password string `json:"password" validate:"required"`
}

// TokenResponse
//
//	@Description	A Token Response
type TokenResponse struct {
	// access token
	AccessToken string `json:"access_token"`
	// refresh token
	RefreshToken string `json:"refresh_token"`
	// expired at
	ExpireAt string `json:"expire_at"`
}

// RefreshTokenRequest
//
//	@Description	Request of Refresh Token
type RefreshTokenRequest struct {
	// refresh token
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// RegisterUserRequest
type RegisterUserRequest struct {
	// the username of the user
	Username string `json:"username" validate:"required"`
	// the email of the user
	Email string `json:"email"    validate:"required"`
	// the password of the user
	Password string `json:"password" validate:"required"`
}

// UpdateAccountRequest
type UpdateAccountRequest struct {
	// the username of the user
	Username string `json:"username" validate:"required"`
	// the email of the user
	Password string `json:"password" validate:"required"`
	// the password of the user
	Email string `json:"email"    validate:"required"`
}

type authHandler struct {
	authService   *services.AuthenticationService
	valService    *services.ValidationService
	jwtMiddleware fiber.Handler
}

func NewAuthHandler(
	jwtMiddleware fiber.Handler,
	authService *services.AuthenticationService,
	valService *services.ValidationService,
) *authHandler {
	return &authHandler{
		authService:   authService,
		valService:    valService,
		jwtMiddleware: jwtMiddleware,
	}
}

func (h *authHandler) RegisterRoutes(router fiber.Router) {
	authApi := router.Group("/auth")

	authApi.Post("/", h.RegisterUser)
	authApi.Post("/login", h.Login)
	authApi.Post("/refresh", h.RefreshToken)

	authApi.Use(h.jwtMiddleware)

	authApi.Put("/", h.UpdateAccount)
	authApi.Delete("/", h.DeleteAccount)
}

// Login godoc
//
//	@Summary	Log In a User
//	@Tags		Authentication
//	@Accept		json
//	@Produce	json
//	@Param		req	body		LoginRequest	true	"login Request"
//	@Success	200	{object}	TokenResponse
//	@Failure	401	{string}	string	"Authentication Failed"
//	@Failure    422 {object}    ValidationErrorResponse
//	@Router		/auth/ [post]
func (h *authHandler) Login(c *fiber.Ctx) error {
	req := new(LoginRequest)

	err := c.BodyParser(req)
	if err != nil {
		return err
	}

	err = h.valService.Validate(req)
	if err != nil {
		return err
	}

	info, err := h.authService.AuthenticateUser(c.Context(), req.Username, req.Password)
	if err != nil {
		if errors.Is(err, ports.ErrUserNotFound) || errors.Is(err, ports.ErrInvalidPassword) {
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid username or password")
		}
		return err
	}

	token, err := h.authService.CreateToken(c.Context(), info.ID)
	if err != nil {
		return err
	}

	return c.JSON(TokenResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpireAt:     token.ExpiresAt.String(),
	})
}

// RefreshToken godoc
//
//	@Summary	Refresh an Access Token
//	@Tags		Authentication
//	@Accept		json
//	@Produce	json
//	@Param		req	body		RefreshTokenRequest	true	"Refresh Token Request"
//	@Success	200	{object}	TokenResponse
//	@Failure	400	{string}	string	"Bad Refresh Token"
//	@Failure	401	{string}	string	"Expired Token"
//	@Failure    422 {object}    ValidationErrorResponse
//	@Router		/auth/refresh [post]
func (h *authHandler) RefreshToken(c *fiber.Ctx) error {
	req := new(RefreshTokenRequest)

	err := c.BodyParser(req)
	if err != nil {
		return err
	}

	err = h.valService.Validate(req)
	if err != nil {
		return err
	}

	token, err := h.authService.RefreshToken(c.Context(), req.RefreshToken)
	if err != nil {
		if errors.Is(err, ports.ErrExpiredToken) {
			return c.Status(fiber.StatusUnauthorized).SendString(ports.ErrExpiredToken.Error())
		}
		if errors.Is(err, ports.ErrInvalidRefreshToken) {
			return c.Status(fiber.StatusBadRequest).
				SendString(ports.ErrInvalidRefreshToken.Error())
		}

		return err
	}
	return c.JSON(TokenResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpireAt:     token.ExpiresAt.String(),
	})
}

// RegisterUser godoc
//
//	@Summary	Register an User
//	@Tags		Authentication
//	@Accept		json
//	@Produce	json
//	@Param		req	body		RegisterUserRequest	true	"Register User Request"
//	@Success	201	{object}	TokenResponse
//	@Failure	409	{string}	string	"User already exists"
//	@Failure    422 {object}    ValidationErrorResponse
//	@Router		/auth/register [post]
func (h *authHandler) RegisterUser(c *fiber.Ctx) error {
	req := new(RegisterUserRequest)

	err := c.BodyParser(req)
	if err != nil {
		return err
	}

	err = h.valService.Validate(req)
	if err != nil {
		return err
	}

	user, err := h.authService.CreateUser(
		c.Context(),
		&services.CreateUserRequest{
			Username: req.Username,
			Email:    req.Email,
			Password: req.Password,
		},
	)
	if err != nil {
		if errors.Is(err, ports.ErrUserAlreadyExists) {
			return c.Status(fiber.StatusConflict).SendString("User already exists")
		}
		return err
	}

	token, err := h.authService.CreateToken(c.Context(), user.ID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(TokenResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpireAt:     token.ExpiresAt.String(),
	})
}

// UpdateAccount godoc
//
//	@Summary	Update an Account
//	@Tags		Authentication
//	@Accept		json
//	@Param		req	body	UpdateAccountRequest	true	"Update Account Request"
//	@Success	200
//	@Failure	409	{string}	string	"User already exists"
//	@Failure    422 {object}    ValidationErrorResponse
//	@Router		/auth/ [put]
func (h *authHandler) UpdateAccount(c *fiber.Ctx) error {
	id := extractTokenFromContext(c)

	req := new(UpdateAccountRequest)

	err := c.BodyParser(req)
	if err != nil {
		return err
	}

	err = h.valService.Validate(req)
	if err != nil {
		return err
	}

	_, err = h.authService.UpdateUser(
		c.Context(),
		id,
		&services.UpdateUserRequest{
			Username: req.Username,
			Email:    req.Email,
			Password: req.Password,
		},
	)
	if err != nil {
		if errors.Is(err, ports.ErrUserAlreadyExists) {
			return c.Status(fiber.StatusConflict).SendString("User already exists")
		}
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// DeleteAccount godoc
//
//	@Summary	Delete an Account
//	@Tags		Authentication
//	@Success	200
//	@Router		/auth/ [delete]
func (h *authHandler) DeleteAccount(c *fiber.Ctx) error {
	id := extractTokenFromContext(c)

	err := h.authService.DeleteUser(c.Context(), id)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}
