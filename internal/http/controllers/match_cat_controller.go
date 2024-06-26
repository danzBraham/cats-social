package controllers

import (
	"errors"
	"net/http"

	"github.com/danzBraham/cats-social/internal/entities/matchcatentity"
	"github.com/danzBraham/cats-social/internal/errors/autherror"
	"github.com/danzBraham/cats-social/internal/errors/matchcaterror"
	"github.com/danzBraham/cats-social/internal/helpers/httphelper"
	"github.com/danzBraham/cats-social/internal/http/middlewares"
	"github.com/danzBraham/cats-social/internal/services"
)

type MatchCatController interface {
	HandleCreateMatchCat(w http.ResponseWriter, r *http.Request)
	HandleGetMatchCats(w http.ResponseWriter, r *http.Request)
	HandleApproveMatchCat(w http.ResponseWriter, r *http.Request)
}

type MatchCatControllerImpl struct {
	MatchCatService services.MatchCatService
}

func NewMatchCatController(matchCatService services.MatchCatService) MatchCatController {
	return &MatchCatControllerImpl{MatchCatService: matchCatService}
}

func (c *MatchCatControllerImpl) HandleCreateMatchCat(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(middlewares.ContextUserIdKey).(string)
	if !ok {
		httphelper.ErrorResponse(w, http.StatusUnauthorized, autherror.ErrUserIdNotFoundInTheContext)
		return
	}

	payload := &matchcatentity.CreateMatchCatRequest{}
	err := httphelper.DecodeAndValidate(w, r, payload)
	if err != nil {
		return
	}

	err = c.MatchCatService.CreateMatchCat(r.Context(), userId, payload)
	for errType, statusCode := range matchcaterror.MatchCatErrorMap {
		if errors.Is(err, errType) {
			httphelper.ErrorResponse(w, statusCode, err)
			return
		}
	}
	if err != nil {
		httphelper.ErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	httphelper.SuccessResponse(w, http.StatusCreated, "successfully send match request", nil)
}

func (c *MatchCatControllerImpl) HandleGetMatchCats(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(middlewares.ContextUserIdKey).(string)
	if !ok {
		httphelper.ErrorResponse(w, http.StatusUnauthorized, autherror.ErrUserIdNotFoundInTheContext)
		return
	}

	matchCatResponses, err := c.MatchCatService.GetMatchCats(r.Context(), userId)
	if err != nil {
		httphelper.ErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	httphelper.SuccessResponse(w, http.StatusOK, "successfully get match requests", matchCatResponses)
}

func (c *MatchCatControllerImpl) HandleApproveMatchCat(w http.ResponseWriter, r *http.Request) {
	payload := &matchcatentity.ApproveMatchCatRequest{}
	err := httphelper.DecodeAndValidate(w, r, payload)
	if err != nil {
		return
	}

	err = c.MatchCatService.ApproveMatchCat(r.Context(), payload)
	if errors.Is(err, matchcaterror.ErrMatchIdNotFound) {
		httphelper.ErrorResponse(w, http.StatusNotFound, err)
		return
	}
	if errors.Is(err, matchcaterror.ErrMatchIdIsNoLongerValid) {
		httphelper.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	if err != nil {
		httphelper.ErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	httphelper.SuccessResponse(w, http.StatusOK, "successfully matches the cat match request", nil)
}
