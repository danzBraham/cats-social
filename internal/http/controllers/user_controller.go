package controllers

import (
	"errors"
	"net/http"
	"time"

	user_entity "github.com/danzBraham/cats-social/internal/entities/user"
	user_error "github.com/danzBraham/cats-social/internal/errors/user"
	http_helper "github.com/danzBraham/cats-social/internal/helpers/http"
	"github.com/danzBraham/cats-social/internal/helpers/validator"
	"github.com/danzBraham/cats-social/internal/services"
)

type UserController interface {
	HandleRegisterUser(w http.ResponseWriter, r *http.Request)
	HandleLoginUser(w http.ResponseWriter, r *http.Request)
}

type UserControllerImpl struct {
	UserService services.UserService
}

func NewUserController(userService services.UserService) UserController {
	return &UserControllerImpl{UserService: userService}
}

func (c *UserControllerImpl) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	payload := &user_entity.RegisterUserRequest{}

	err := http_helper.DecodeJSON(r, payload)
	if err != nil {
		http_helper.EncodeJSON(w, http.StatusBadRequest, http_helper.ResponseBody{
			Error:   "bad request",
			Message: err.Error(),
		})
		return
	}

	err = validator.ValidatePayload(payload)
	if err != nil {
		http_helper.EncodeJSON(w, http.StatusBadRequest, http_helper.ResponseBody{
			Error:   "request doesn't pass validation",
			Message: err.Error(),
		})
		return
	}

	userResponse, err := c.UserService.RegisterUser(r.Context(), payload)
	if errors.Is(err, user_error.ErrEmailAlreadyExists) {
		http_helper.EncodeJSON(w, http.StatusConflict, http_helper.ResponseBody{
			Error:   "conflict",
			Message: err.Error(),
		})
		return
	}
	if err != nil {
		http_helper.EncodeJSON(w, http.StatusInternalServerError, http_helper.ResponseBody{
			Error:   "internal server",
			Message: err.Error(),
		})
		return
	}

	cookie := &http.Cookie{
		Name:    "Authorization",
		Value:   userResponse.AccessToken,
		Expires: time.Now().Add(8 * time.Hour),
	}
	http.SetCookie(w, cookie)

	http_helper.EncodeJSON(w, http.StatusCreated, http_helper.ResponseBody{
		Message: "User registered successfully",
		Data:    userResponse,
	})
}

func (c *UserControllerImpl) HandleLoginUser(w http.ResponseWriter, r *http.Request) {
	payload := &user_entity.LoginUserRequest{}

	err := http_helper.DecodeJSON(r, payload)
	if err != nil {
		http_helper.EncodeJSON(w, http.StatusBadRequest, http_helper.ResponseBody{
			Error:   "bad request",
			Message: err.Error(),
		})
		return
	}

	err = validator.ValidatePayload(payload)
	if err != nil {
		http_helper.EncodeJSON(w, http.StatusBadRequest, http_helper.ResponseBody{
			Error:   "request doesn't pass validation",
			Message: err.Error(),
		})
		return
	}

	userResponse, err := c.UserService.LoginUser(r.Context(), payload)
	if errors.Is(err, user_error.ErrUserNotFound) {
		http_helper.EncodeJSON(w, http.StatusNotFound, http_helper.ResponseBody{
			Error:   "not found",
			Message: err.Error(),
		})
		return
	}
	if errors.Is(err, user_error.ErrInvalidPassword) {
		http_helper.EncodeJSON(w, http.StatusBadRequest, http_helper.ResponseBody{
			Error:   "bad request",
			Message: err.Error(),
		})
		return
	}
	if err != nil {
		http_helper.EncodeJSON(w, http.StatusInternalServerError, http_helper.ResponseBody{
			Error:   "internal server",
			Message: err.Error(),
		})
		return
	}

	cookie := &http.Cookie{
		Name:    "Authorizaiton",
		Value:   userResponse.AccessToken,
		Expires: time.Now().Add(8 * time.Hour),
	}
	http.SetCookie(w, cookie)

	http_helper.EncodeJSON(w, http.StatusOK, http_helper.ResponseBody{
		Message: "User logged successfully",
		Data:    userResponse,
	})
}
