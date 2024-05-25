package validator

import (
	"net/url"
	"path"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	Validator *validator.Validate
}

var validate = validator.New(validator.WithRequiredStructEnabled())

func InitCustomValidation() {
	validate.RegisterValidation("imageurl", validateImageURL)
}

func ValidatePayload(payload interface{}) error {
	if err := validate.Struct(payload); err != nil {
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
