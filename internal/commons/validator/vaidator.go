package validator_common

import "github.com/go-playground/validator/v10"

var validate = validator.New(validator.WithRequiredStructEnabled())

func ValidatePayload(payload interface{}) error {
	if err := validate.Struct(payload); err != nil {
		return err.(validator.ValidationErrors)
	}
	return nil
}
