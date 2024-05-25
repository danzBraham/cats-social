package controllers

import (
	"errors"
	"net/http"

	match_exception "github.com/danzbraham/cats-social/internal/commons/exceptions/match"
	http_common "github.com/danzbraham/cats-social/internal/commons/http"
	"github.com/danzbraham/cats-social/internal/commons/validator"
	match_entity "github.com/danzbraham/cats-social/internal/entities/match"
	"github.com/danzbraham/cats-social/internal/middlewares"
	"github.com/danzbraham/cats-social/internal/services"
	"github.com/go-chi/chi/v5"
)

type MatchController struct {
	Service services.MatchService
}

func NewMatchController(service services.MatchService) *MatchController {
	return &MatchController{Service: service}
}

func (c *MatchController) Routes() chi.Router {
	r := chi.NewRouter()

	r.Use(middlewares.AuthMiddleware)
	r.Post("/", c.handleMatchCatRequest)

	return r
}

func (c *MatchController) handleMatchCatRequest(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(middlewares.ContextUserIdKey).(string)
	if !ok {
		http_common.ResponseError(w, http.StatusBadRequest, "User Id type assertion failed", "User Id not found in context")
		return
	}
	payload := &match_entity.MatchCatRequest{Issuer: userId}

	if err := http_common.DecodeJSON(r, payload); err != nil {
		http_common.ResponseError(w, http.StatusBadRequest, err.Error(), "Failed to decode JSON")
		return
	}

	if err := validator.ValidatePayload(payload); err != nil {
		http_common.ResponseError(w, http.StatusBadRequest, err.Error(), "Request doesn't pass validation")
		return
	}

	err := c.Service.MatchCat(r.Context(), payload)
	if errors.Is(err, match_exception.ErrMatchCatIdIsNotFound) {
		http_common.ResponseError(w, http.StatusNotFound, "Not found error", err.Error())
		return
	}
	if errors.Is(err, match_exception.ErrUserCatIdIsNotFound) {
		http_common.ResponseError(w, http.StatusNotFound, "Not found error", err.Error())
		return
	}
	if errors.Is(err, match_exception.ErrUserCatIdNotBelongTheUser) {
		http_common.ResponseError(w, http.StatusNotFound, "Not found error", err.Error())
		return
	}
	if errors.Is(err, match_exception.ErrBothCatHaveSameGender) {
		http_common.ResponseError(w, http.StatusBadRequest, "Bad request error", err.Error())
		return
	}
	if errors.Is(err, match_exception.ErrBothCatAlreadyMatched) {
		http_common.ResponseError(w, http.StatusBadRequest, "Bad request error", err.Error())
		return
	}
	if errors.Is(err, match_exception.ErrBothCatsHaveTheSameOwner) {
		http_common.ResponseError(w, http.StatusBadRequest, "Bad request error", err.Error())
		return
	}
	if err != nil {
		http_common.ResponseError(w, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	http_common.ResponseSuccess(w, http.StatusCreated, "successfully send match request", nil)
}
