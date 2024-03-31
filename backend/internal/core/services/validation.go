package services

import (
	"fmt"

	"github.com/go-playground/validator/v10"
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

func NewValidationService() *ValidationService {
	return &ValidationService{
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (v ValidationService) Validate(i interface{}) error {
	validationErrors := make([]ErrorMessage, 0)

	errs := v.validate.Struct(i)
	if errs != nil {
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

		return &ValidationError{errors: validationErrors}
	}

	return nil
}
