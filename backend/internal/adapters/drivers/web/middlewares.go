package web

import (
	"fmt"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"github.com/taldoflemis/brain.test/internal/core/services"
	"github.com/taldoflemis/brain.test/internal/ports"
)

func NewJWTMiddleware(authManager ports.AuthenticationManager) fiber.Handler {
	return jwtware.New(jwtware.Config{
		KeyFunc: customKeyFunc(authManager),
	})
}

func customKeyFunc(authManager ports.AuthenticationManager) jwt.Keyfunc {
	return func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != authManager.GetAlgorithm() {
			return nil, fmt.Errorf("Unexpected jwt signing method=%v", t.Header["alg"])
		}

		return authManager.GetPublicKey(), nil
	}
}

type ValidationErrorResponse struct {
	Errors []string `json:"errors"`
}

func ErrorHandlerMiddleware(c *fiber.Ctx, err error) error {
	if err != nil {
		validationErrors, ok := err.(*services.ValidationError)
		if ok {
			resp := convertValidationErrorsToResponse(validationErrors)
			return c.Status(fiber.StatusUnprocessableEntity).JSON(resp)
		}
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return nil
}

func convertValidationErrorsToResponse(
	validationErrors *services.ValidationError,
) ValidationErrorResponse {
	messages := validationErrors.GetMessages()
	error := ValidationErrorResponse{
		Errors: make([]string, len(messages)),
	}

	for i, message := range messages {
		error.Errors[i] = message.Message
	}

	return error
}
