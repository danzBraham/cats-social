package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func formatValidationErrors(err error) string {
	if validationError, ok := err.(validator.ValidationErrors); ok {
		var sb strings.Builder
		for _, fieldError := range validationError {
			sb.WriteString(fmt.Sprintf("field '%s' failed on the '%s' tag; ", fieldError.Field(), fieldError.Tag()))
		}
		return sb.String()
	}
	return err.Error()
}

func ValidatePayload(payload interface{}) error {
	if err := validate.Struct(payload); err != nil {
		return fmt.Errorf(formatValidationErrors(err))
	}
	return nil
}
