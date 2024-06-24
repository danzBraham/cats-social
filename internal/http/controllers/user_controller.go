package controllers

import (
	"errors"
	"net/http"
	"time"

	"github.com/danzBraham/cats-social/internal/entities/user_entity"
	"github.com/danzBraham/cats-social/internal/errors/user_error"
	"github.com/danzBraham/cats-social/internal/helpers/http_helper"
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
	err := http_helper.DecodeAndValidate(w, r, payload)
	if err != nil {
		return
	}

	userResponse, err := c.UserService.RegisterUser(r.Context(), payload)
	if errors.Is(err, user_error.ErrEmailAlreadyExists) {
		http_helper.HandleErrorResponse(w, http.StatusConflict, err)
		return
	}
	if err != nil {
		http_helper.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	cookie := &http.Cookie{
		Name:    "Authorization",
		Value:   userResponse.AccessToken,
		Expires: time.Now().Add(8 * time.Hour),
	}
	http.SetCookie(w, cookie)

	http_helper.HandleSuccessResponse(w, http.StatusCreated, "User registered successfully", userResponse)
}

func (c *UserControllerImpl) HandleLoginUser(w http.ResponseWriter, r *http.Request) {
	payload := &user_entity.LoginUserRequest{}
	err := http_helper.DecodeAndValidate(w, r, payload)
	if err != nil {
		return
	}

	userResponse, err := c.UserService.LoginUser(r.Context(), payload)
	if errors.Is(err, user_error.ErrUserNotFound) {
		http_helper.HandleErrorResponse(w, http.StatusNotFound, err)
		return
	}
	if errors.Is(err, user_error.ErrInvalidPassword) {
		http_helper.HandleErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	if err != nil {
		http_helper.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	cookie := &http.Cookie{
		Name:    "Authorizaiton",
		Value:   userResponse.AccessToken,
		Expires: time.Now().Add(8 * time.Hour),
	}
	http.SetCookie(w, cookie)

	http_helper.HandleSuccessResponse(w, http.StatusOK, "User logged successfully", userResponse)
}
