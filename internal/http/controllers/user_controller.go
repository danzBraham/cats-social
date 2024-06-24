package controllers

import (
	"errors"
	"net/http"
	"time"

	"github.com/danzBraham/cats-social/internal/entities/userentity"
	"github.com/danzBraham/cats-social/internal/errors/usererror"
	"github.com/danzBraham/cats-social/internal/helpers/httphelper"
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
	payload := &userentity.RegisterUserRequest{}
	err := httphelper.DecodeAndValidate(w, r, payload)
	if err != nil {
		return
	}

	userResponse, err := c.UserService.RegisterUser(r.Context(), payload)
	if errors.Is(err, usererror.ErrEmailAlreadyExists) {
		httphelper.HandleErrorResponse(w, http.StatusConflict, err)
		return
	}
	if err != nil {
		httphelper.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	cookie := &http.Cookie{
		Name:    "Authorization",
		Value:   userResponse.AccessToken,
		Expires: time.Now().Add(8 * time.Hour),
	}
	http.SetCookie(w, cookie)

	httphelper.HandleSuccessResponse(w, http.StatusCreated, "User registered successfully", userResponse)
}

func (c *UserControllerImpl) HandleLoginUser(w http.ResponseWriter, r *http.Request) {
	payload := &userentity.LoginUserRequest{}
	err := httphelper.DecodeAndValidate(w, r, payload)
	if err != nil {
		return
	}

	userResponse, err := c.UserService.LoginUser(r.Context(), payload)
	if errors.Is(err, usererror.ErrUserNotFound) {
		httphelper.HandleErrorResponse(w, http.StatusNotFound, err)
		return
	}
	if errors.Is(err, usererror.ErrInvalidPassword) {
		httphelper.HandleErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	if err != nil {
		httphelper.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	cookie := &http.Cookie{
		Name:    "Authorizaiton",
		Value:   userResponse.AccessToken,
		Expires: time.Now().Add(8 * time.Hour),
	}
	http.SetCookie(w, cookie)

	httphelper.HandleSuccessResponse(w, http.StatusOK, "User logged successfully", userResponse)
}
