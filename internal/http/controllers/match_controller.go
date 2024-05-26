package controllers

import (
	"errors"
	"net/http"

	match_exception "github.com/danzbraham/cats-social/internal/commons/exceptions/match"
	http_common "github.com/danzbraham/cats-social/internal/commons/http"
	"github.com/danzbraham/cats-social/internal/commons/validator"
	match_entity "github.com/danzbraham/cats-social/internal/entities/match"
	"github.com/danzbraham/cats-social/internal/http/middlewares"
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
	r.Get("/", c.handleGetMatchCatRequests)
	r.Post("/approve", c.handleApproveMatchCatRequest)
	r.Post("/reject", c.handleRejectMatchCatRequest)
	r.Delete("/{matchId}", c.handleDeleteMatchCatRequest)

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

	err := c.Service.RequestMatchCat(r.Context(), payload)
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

func (c *MatchController) handleGetMatchCatRequests(w http.ResponseWriter, r *http.Request) {
	matchCatsResponse, err := c.Service.GetMatchCatRequests(r.Context())
	if err != nil {
		http_common.ResponseError(w, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	http_common.ResponseSuccess(w, http.StatusCreated, "successfully get match requests", matchCatsResponse)
}

func (c *MatchController) handleApproveMatchCatRequest(w http.ResponseWriter, r *http.Request) {
	payload := &match_entity.DecisionMatchRequest{}

	if err := http_common.DecodeJSON(r, payload); err != nil {
		http_common.ResponseError(w, http.StatusBadRequest, err.Error(), "Failed to decode JSON")
		return
	}

	if err := validator.ValidatePayload(payload); err != nil {
		http_common.ResponseError(w, http.StatusBadRequest, err.Error(), "Request doesn't pass validation")
		return
	}

	err := c.Service.ApproveMatchCatRequest(r.Context(), payload)
	if errors.Is(err, match_exception.ErrMatchIdIsNotFound) {
		http_common.ResponseError(w, http.StatusNotFound, "Not found error", err.Error())
		return
	}
	if errors.Is(err, match_exception.ErrMatchIdIsNoLongerValid) {
		http_common.ResponseError(w, http.StatusBadRequest, "Bad request error", err.Error())
		return
	}
	if err != nil {
		http_common.ResponseError(w, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	http_common.ResponseSuccess(w, http.StatusOK, "successfully matches the cat match request", nil)
}

func (c *MatchController) handleRejectMatchCatRequest(w http.ResponseWriter, r *http.Request) {
	payload := &match_entity.DecisionMatchRequest{}

	if err := http_common.DecodeJSON(r, payload); err != nil {
		http_common.ResponseError(w, http.StatusBadRequest, err.Error(), "Failed to decode JSON")
		return
	}

	if err := validator.ValidatePayload(payload); err != nil {
		http_common.ResponseError(w, http.StatusBadRequest, err.Error(), "Request doesn't pass validation")
		return
	}

	err := c.Service.RejectMatchCatRequest(r.Context(), payload)
	if errors.Is(err, match_exception.ErrMatchIdIsNotFound) {
		http_common.ResponseError(w, http.StatusNotFound, "Not found error", err.Error())
		return
	}
	if errors.Is(err, match_exception.ErrMatchIdIsNoLongerValid) {
		http_common.ResponseError(w, http.StatusBadRequest, "Bad request error", err.Error())
		return
	}
	if err != nil {
		http_common.ResponseError(w, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	http_common.ResponseSuccess(w, http.StatusOK, "successfully reject the cat match request", nil)
}

func (c *MatchController) handleDeleteMatchCatRequest(w http.ResponseWriter, r *http.Request) {
	matchId := chi.URLParam(r, "matchId")
	userId, ok := r.Context().Value(middlewares.ContextUserIdKey).(string)
	if !ok {
		http_common.ResponseError(w, http.StatusBadRequest, "User Id type assertion failed", "User Id not found in context")
		return
	}
	payload := &match_entity.DeleteMatchCatRequest{MatchId: matchId, Issuer: userId}

	if err := validator.ValidatePayload(payload); err != nil {
		http_common.ResponseError(w, http.StatusBadRequest, err.Error(), "Request doesn't pass validation")
		return
	}

	err := c.Service.DeleteMatchCatRequest(r.Context(), payload)
	if err != nil {
		http_common.ResponseError(w, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	http_common.ResponseSuccess(w, http.StatusOK, "successfully remove a cat match request", nil)
}
