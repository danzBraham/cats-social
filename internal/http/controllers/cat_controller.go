package controllers

import (
	"errors"
	"net/http"

	cat_exception "github.com/danzbraham/cats-social/internal/commons/exceptions/cats"
	http_common "github.com/danzbraham/cats-social/internal/commons/http"
	"github.com/danzbraham/cats-social/internal/commons/validator"
	cat_entity "github.com/danzbraham/cats-social/internal/entities/cat"
	"github.com/danzbraham/cats-social/internal/http/middlewares"
	"github.com/danzbraham/cats-social/internal/services"
	"github.com/go-chi/chi/v5"
)

type CatController struct {
	Service services.CatService
}

func NewCatController(service services.CatService) *CatController {
	return &CatController{Service: service}
}

func (c *CatController) Routes() chi.Router {
	r := chi.NewRouter()

	r.Use(middlewares.AuthMiddleware)
	r.Post("/", c.handleAddCat)
	r.Get("/", c.handleGetCats)
	r.Put("/{catId}", c.handleUpdateCat)
	r.Delete("/{catId}", c.handleDeleteCat)

	return r
}

func (c *CatController) handleAddCat(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(middlewares.ContextUserIdKey).(string)
	if !ok {
		http_common.ResponseError(w, http.StatusBadRequest, "User Id type assertion failed", "User Id not found in context")
		return
	}
	payload := &cat_entity.AddCatRequest{OwnerId: userId}

	if err := http_common.DecodeJSON(r, payload); err != nil {
		http_common.ResponseError(w, http.StatusBadRequest, err.Error(), "Failed to decode JSON")
		return
	}

	if err := validator.ValidatePayload(payload); err != nil {
		http_common.ResponseError(w, http.StatusBadRequest, err.Error(), "Request doesn't pass validation")
		return
	}

	catResponse, err := c.Service.AddCat(r.Context(), payload)
	if err != nil {
		http_common.ResponseError(w, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	http_common.ResponseSuccess(w, http.StatusCreated, "successfully add cat", catResponse)
}

func (c *CatController) handleGetCats(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	params := &cat_entity.CatQueryParams{
		Id:         query.Get("id"),
		Limit:      "5",
		Offset:     "0",
		Race:       query.Get("race"),
		Sex:        query.Get("sex"),
		HasMatched: query.Get("hasMatched"),
		AgeInMonth: query.Get("ageInMonth"),
		Owned:      query.Get("owned"),
		Search:     query.Get("search"),
	}

	if limit := query.Get("limit"); limit != "" {
		params.Limit = limit
	}

	if offset := query.Get("offset"); offset != "" {
		params.Offset = offset
	}

	catsResponse, err := c.Service.GetCats(r.Context(), params)
	if err != nil {
		http_common.ResponseError(w, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	http_common.ResponseSuccess(w, http.StatusOK, "successfully get cats", catsResponse)
}

func (c *CatController) handleUpdateCat(w http.ResponseWriter, r *http.Request) {
	catId := chi.URLParam(r, "catId")
	payload := &cat_entity.UpdateCatRequest{}

	if err := http_common.DecodeJSON(r, payload); err != nil {
		http_common.ResponseError(w, http.StatusBadRequest, err.Error(), "Failed to decode JSON")
		return
	}

	if err := validator.ValidatePayload(payload); err != nil {
		http_common.ResponseError(w, http.StatusBadRequest, err.Error(), "Request doesn't pass validation")
		return
	}

	err := c.Service.UpdateCat(r.Context(), catId, payload)
	if errors.Is(err, cat_exception.ErrCatIdIsNotFound) {
		http_common.ResponseError(w, http.StatusNotFound, "Not found error", err.Error())
		return
	}
	if err != nil {
		http_common.ResponseError(w, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	http_common.ResponseSuccess(w, http.StatusOK, "successfully update cat", nil)
}

func (c *CatController) handleDeleteCat(w http.ResponseWriter, r *http.Request) {
	catId := chi.URLParam(r, "catId")

	err := c.Service.DeleteCat(r.Context(), catId)
	if errors.Is(err, cat_exception.ErrCatIdIsNotFound) {
		http_common.ResponseError(w, http.StatusNotFound, "Not found error", err.Error())
		return
	}
	if err != nil {
		http_common.ResponseError(w, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	http_common.ResponseSuccess(w, http.StatusOK, "successfully delete cat", nil)
}
