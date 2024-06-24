package controllers

import (
	"net/http"
	"strconv"

	"github.com/danzBraham/cats-social/internal/entities/catentity"
	"github.com/danzBraham/cats-social/internal/errors/autherror"
	"github.com/danzBraham/cats-social/internal/helpers/httphelper"
	"github.com/danzBraham/cats-social/internal/http/middlewares"
	"github.com/danzBraham/cats-social/internal/services"
)

type CatController interface {
	HandleCreateCat(w http.ResponseWriter, r *http.Request)
	HandleGetCats(w http.ResponseWriter, r *http.Request)
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

func (c *CatControllerImpl) HandleGetCats(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(middlewares.ContextUserIdKey).(string)
	if !ok {
		httphelper.HandleErrorResponse(w, http.StatusUnauthorized, autherror.ErrUserIdNotFoundInTheContext)
		return
	}

	query := r.URL.Query()
	params := &catentity.CatQueryParams{
		Id:         query.Get("id"),
		Limit:      5,
		Offset:     0,
		Race:       catentity.Race(query.Get("race")),
		Sex:        catentity.Sex(query.Get("sex")),
		AgeInMonth: query.Get("ageInMonth"),
		Search:     query.Get("search"),
	}

	if limit := query.Get("limit"); limit != "" {
		params.Limit, _ = strconv.Atoi(limit)
	}

	if offset := query.Get("offset"); offset != "" {
		params.Offset, _ = strconv.Atoi(offset)
	}

	if hasMatched := query.Get("hasMatched"); hasMatched != "" {
		params.HasMatched, _ = strconv.ParseBool(hasMatched)
	}

	if owned := query.Get("owned"); owned != "" {
		params.Owned, _ = strconv.ParseBool(owned)
	}

	catsResponse, err := c.CatService.GetCats(r.Context(), userId, params)
	if err != nil {
		httphelper.HandleErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	httphelper.HandleSuccessResponse(w, http.StatusOK, "success", catsResponse)
}
