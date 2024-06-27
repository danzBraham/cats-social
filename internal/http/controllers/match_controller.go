package controllers

import (
	"errors"
	"net/http"

	"github.com/danzBraham/cats-social/internal/entities/matchentity"
	"github.com/danzBraham/cats-social/internal/errors/autherror"
	"github.com/danzBraham/cats-social/internal/errors/matcherror"
	"github.com/danzBraham/cats-social/internal/helpers/httphelper"
	"github.com/danzBraham/cats-social/internal/http/middlewares"
	"github.com/danzBraham/cats-social/internal/services"
)

type MatchController interface {
	HandleCreateMatch(w http.ResponseWriter, r *http.Request)
	HandleGetMatches(w http.ResponseWriter, r *http.Request)
	HandleApproveMatch(w http.ResponseWriter, r *http.Request)
	HandleRejectMatch(w http.ResponseWriter, r *http.Request)
}

type MatchControllerImpl struct {
	MatchService services.MatchService
}

func NewMatchController(matchService services.MatchService) MatchController {
	return &MatchControllerImpl{MatchService: matchService}
}

func (c *MatchControllerImpl) HandleCreateMatch(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(middlewares.ContextUserIdKey).(string)
	if !ok {
		httphelper.ErrorResponse(w, http.StatusUnauthorized, autherror.ErrUserIdNotFoundInTheContext)
		return
	}

	payload := &matchentity.CreateMatchRequest{}
	err := httphelper.DecodeAndValidate(w, r, payload)
	if err != nil {
		return
	}

	err = c.MatchService.CreateMatch(r.Context(), userId, payload)
	if errors.Is(err, matcherror.ErrMatchCatIdNotFound) {
		httphelper.ErrorResponse(w, http.StatusNotFound, err)
		return
	}
	if errors.Is(err, matcherror.ErrUserCatIdNotFound) {
		httphelper.ErrorResponse(w, http.StatusNotFound, err)
		return
	}
	if errors.Is(err, matcherror.ErrUserCatIdNotBelongToTheUser) {
		httphelper.ErrorResponse(w, http.StatusNotFound, err)
		return
	}
	if errors.Is(err, matcherror.ErrBothCatsHaveSameGender) {
		httphelper.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	if errors.Is(err, matcherror.ErrBothCatsHaveSameOwner) {
		httphelper.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	if errors.Is(err, matcherror.ErrBothCatsHaveAlreadyMatched) {
		httphelper.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	if err != nil {
		httphelper.ErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	httphelper.SuccessResponse(w, http.StatusCreated, "successfully send match request", nil)
}

func (c *MatchControllerImpl) HandleGetMatches(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(middlewares.ContextUserIdKey).(string)
	if !ok {
		httphelper.ErrorResponse(w, http.StatusUnauthorized, autherror.ErrUserIdNotFoundInTheContext)
		return
	}

	matchResponses, err := c.MatchService.GetMatches(r.Context(), userId)
	if err != nil {
		httphelper.ErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	httphelper.SuccessResponse(w, http.StatusOK, "successfully get match requests", matchResponses)
}

func (c *MatchControllerImpl) HandleApproveMatch(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(middlewares.ContextUserIdKey).(string)
	if !ok {
		httphelper.ErrorResponse(w, http.StatusUnauthorized, autherror.ErrUserIdNotFoundInTheContext)
		return
	}

	payload := &matchentity.ApproveMatchRequest{}
	err := httphelper.DecodeAndValidate(w, r, payload)
	if err != nil {
		return
	}

	err = c.MatchService.ApproveMatch(r.Context(), userId, payload)
	if errors.Is(err, matcherror.ErrMatchIdNotFound) {
		httphelper.ErrorResponse(w, http.StatusNotFound, err)
		return
	}
	if errors.Is(err, matcherror.ErrMatchIdIsNoLongerValid) {
		httphelper.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	if errors.Is(err, matcherror.ErrUnauthorizedDecision) {
		httphelper.ErrorResponse(w, http.StatusForbidden, err)
		return
	}
	if err != nil {
		httphelper.ErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	httphelper.SuccessResponse(w, http.StatusOK, "successfully matches the cat match request", nil)
}

func (c *MatchControllerImpl) HandleRejectMatch(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(middlewares.ContextUserIdKey).(string)
	if !ok {
		httphelper.ErrorResponse(w, http.StatusUnauthorized, autherror.ErrUserIdNotFoundInTheContext)
		return
	}

	payload := &matchentity.RejectMatchRequest{}
	err := httphelper.DecodeAndValidate(w, r, payload)
	if err != nil {
		return
	}

	err = c.MatchService.RejectMatch(r.Context(), userId, payload)
	if errors.Is(err, matcherror.ErrMatchIdNotFound) {
		httphelper.ErrorResponse(w, http.StatusNotFound, err)
		return
	}
	if errors.Is(err, matcherror.ErrMatchIdIsNoLongerValid) {
		httphelper.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	if errors.Is(err, matcherror.ErrUnauthorizedDecision) {
		httphelper.ErrorResponse(w, http.StatusForbidden, err)
		return
	}
	if err != nil {
		httphelper.ErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	httphelper.SuccessResponse(w, http.StatusOK, "successfully reject the cat match request", nil)
}
