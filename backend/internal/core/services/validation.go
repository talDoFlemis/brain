package services

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"

	game "github.com/taldoflemis/brain.test/internal/core/domain/game_aggregate"
)

type ValidationService struct {
	validate *validator.Validate
}

type ErrorMessage struct {
	Message string
}

type ValidationError struct {
	errors []ErrorMessage
}

func (validationerror *ValidationError) Error() string {
	return "failed to validate struct"
}

func (validationerror *ValidationError) GetMessages() []ErrorMessage {
	return validationerror.errors
}

func (v *ValidationError) AddNewMessage(e ErrorMessage) {
	v.errors = append(v.errors, e)
}

func NewValidationService() *ValidationService {
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.RegisterValidation("atleastonecorrect", game.ValidateAtLeastOneCorrect)
	if err != nil {
		panic(err)
	}
	return &ValidationService{
		validate: validate,
	}
}

func (v ValidationService) Validate(i interface{}) error {
	validationErrors := make([]ErrorMessage, 0)

	errs := v.validate.Struct(i)
	if errs != nil {
		var ve validator.ValidationErrors
		if errors.As(errs, &ve) {
			for _, err := range errs.(validator.ValidationErrors) {
				var elem ErrorMessage

				elem.Message = fmt.Sprintf(
					"[%s]: '%v' | Needs to implement '%s'",
					err.Field(),
					err.Value(),
					err.Tag(),
				)
				validationErrors = append(validationErrors, elem)
			}
		} else {
			return errs
		}

		return &ValidationError{errors: validationErrors}
	}

	return nil
}
