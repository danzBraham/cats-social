package controllers

import (
	"net/http"

	"github.com/danzBraham/cats-social/internal/entities/catentity"
	"github.com/danzBraham/cats-social/internal/errors/autherror"
	"github.com/danzBraham/cats-social/internal/helpers/httphelper"
	"github.com/danzBraham/cats-social/internal/http/middlewares"
	"github.com/danzBraham/cats-social/internal/services"
)

type CatController interface {
	HandleCreateCat(w http.ResponseWriter, r *http.Request)
}

type CatControllerImpl struct {
	CatService services.CatService
}

func NewCatController(catService services.CatService) CatController {
	return &CatControllerImpl{CatService: catService}
}

func (c *CatControllerImpl) HandleCreateCat(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(middlewares.ContextUserIdKey).(string)
	if !ok {
		httphelper.HandleErrorResponse(w, http.StatusUnauthorized, autherror.ErrUserIdNotFoundInTheContext)
		return
	}

	payload := &catentity.CreateCatRequest{}
	httphelper.DecodeAndValidate(w, r, payload)

	catResponse, err := c.CatService.CreateCat(r.Context(), userId, payload)
	if err != nil {
		httphelper.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	httphelper.HandleSuccessResponse(w, http.StatusCreated, "success", catResponse)
}
