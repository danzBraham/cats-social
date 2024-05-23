package controllers

import (
	"errors"
	"net/http"
	"time"

	"github.com/danzbraham/cats-social/internal/applications/usecases"
	user_exception "github.com/danzbraham/cats-social/internal/commons/exceptions/users"
	http_common "github.com/danzbraham/cats-social/internal/commons/http"
	validator_common "github.com/danzbraham/cats-social/internal/commons/validator"
	user_entity "github.com/danzbraham/cats-social/internal/domains/entities/users"
	"github.com/go-chi/chi/v5"
)

type UserController struct {
	Usecase usecases.UserUsecase
}

func NewUserController(usecase usecases.UserUsecase) *UserController {
	return &UserController{Usecase: usecase}
}

func (c *UserController) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/register", c.handleRegisterUser)
	r.Post("/login", c.handleLoginUser)

	return r
}

func (c *UserController) handleRegisterUser(w http.ResponseWriter, r *http.Request) {
	payload := &user_entity.RegisterUserRequest{}

	if err := http_common.DecodeJSON(r, payload); err != nil {
		http_common.ResponseError(w, http.StatusBadRequest, err.Error(), "Failed to decode JSON")
		return
	}

	if err := validator_common.ValidatePayload(payload); err != nil {
		http_common.ResponseError(w, http.StatusBadRequest, err.Error(), "Request doesn't pass validation")
		return
	}

	userResponse, err := c.Usecase.RegisterUser(r.Context(), payload)
	if errors.Is(err, user_exception.ErrEmailAlreadyExists) {
		http_common.ResponseError(w, http.StatusConflict, "Conflict error", err.Error())
		return
	}
	if err != nil {
		http_common.ResponseError(w, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	cookie := &http.Cookie{
		Name:    "Authorization",
		Value:   userResponse.AccessToken,
		Expires: time.Now().Add(2 * time.Hour),
	}
	http.SetCookie(w, cookie)

	http_common.ResponseSuccess(w, http.StatusCreated, "User successfully registered", userResponse)
}

func (c *UserController) handleLoginUser(w http.ResponseWriter, r *http.Request) {
	payload := &user_entity.LoginUserRequest{}

	if err := http_common.DecodeJSON(r, payload); err != nil {
		http_common.ResponseError(w, http.StatusBadRequest, err.Error(), "Failed to decode JSON")
		return
	}

	if err := validator_common.ValidatePayload(payload); err != nil {
		http_common.ResponseError(w, http.StatusBadRequest, err.Error(), "Request doesn't pass validation")
		return
	}

	userResponse, err := c.Usecase.LoginUser(r.Context(), payload)
	if errors.Is(err, user_exception.ErrUserNotFound) {
		http_common.ResponseError(w, http.StatusNotFound, "Not found error", err.Error())
		return
	}
	if errors.Is(err, user_exception.ErrInvalidPassword) {
		http_common.ResponseError(w, http.StatusBadRequest, "Bad request error", err.Error())
		return
	}
	if err != nil {
		http_common.ResponseError(w, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	cookie := &http.Cookie{
		Name:    "Authorization",
		Value:   userResponse.AccessToken,
		Expires: time.Now().Add(2 * time.Hour),
	}
	http.SetCookie(w, cookie)

	http_common.ResponseSuccess(w, http.StatusCreated, "User successfully logged", userResponse)
}
