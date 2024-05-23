package securities_impl

import (
	"github.com/danzbraham/cats-social/internal/applications/securities"
	"github.com/go-playground/validator/v10"
)

type GoValidator struct {
	Validator *validator.Validate
}

func NewGoValidator(validator *validator.Validate) securities.Validator {
	return &GoValidator{Validator: validator}
}

func (gv *GoValidator) ValidatePayload(payload interface{}) error {
	if err := gv.Validator.Struct(payload); err != nil {
		return err.(validator.ValidationErrors)
	}
	return nil
}
