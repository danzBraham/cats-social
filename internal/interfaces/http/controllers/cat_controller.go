package controllers

import (
	"net/http"

	"github.com/danzbraham/cats-social/internal/applications/securities"
	"github.com/danzbraham/cats-social/internal/applications/usecases"
	http_common "github.com/danzbraham/cats-social/internal/commons/http"
	cat_entity "github.com/danzbraham/cats-social/internal/domains/entities/cats"
	"github.com/go-chi/chi/v5"
)

type CatController struct {
	Usecase   usecases.CatUsecase
	Validator securities.Validator
}

func NewCatController(usecase usecases.CatUsecase, validator securities.Validator) *CatController {
	return &CatController{
		Usecase:   usecase,
		Validator: validator,
	}
}

func (c *CatController) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", c.handleAddCat)

	return r
}

func (c *CatController) handleAddCat(w http.ResponseWriter, r *http.Request) {
	payload := &cat_entity.AddCatRequest{}

	if err := http_common.DecodeJSON(r, payload); err != nil {
		http_common.ResponseError(w, http.StatusBadRequest, err.Error(), "Failed to decode JSON")
		return
	}

	if err := c.Validator.ValidatePayload(payload); err != nil {
		http_common.ResponseError(w, http.StatusBadRequest, err.Error(), "Request doesn't pass validation")
		return
	}

	catResponse, err := c.Usecase.AddCat(r.Context(), payload)
	if err != nil {
		http_common.ResponseError(w, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	http_common.ResponseSuccess(w, http.StatusCreated, "successfully add cat", catResponse)
}
