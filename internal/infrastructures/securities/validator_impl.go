package securities_impl

import (
	"net/url"
	"path"

	"github.com/danzbraham/cats-social/internal/applications/securities"
	"github.com/go-playground/validator/v10"
)

type GoValidator struct {
	Validator *validator.Validate
}

func NewGoValidator(validator *validator.Validate) securities.Validator {
	goValidator := &GoValidator{Validator: validator}
	goValidator.Validator.RegisterValidation("imageurl", validateImageURL)
	return goValidator
}

func (gv *GoValidator) ValidatePayload(payload interface{}) error {
	if err := gv.Validator.Struct(payload); err != nil {
		return err.(validator.ValidationErrors)
	}
	return nil
}

func validateImageURL(fl validator.FieldLevel) bool {
	u, err := url.ParseRequestURI(fl.Field().String())
	if err != nil {
		return false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	if u.Host == "" {
		return false
	}
	ext := path.Ext(u.Path)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return false
	}
	return true
}
