package server

import (
	"net/http"

	http_helper "github.com/danzBraham/cats-social/internal/helpers/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http_helper.EncodeJSON(w, http.StatusOK, http_helper.ResponseBody{
			Message: "welcome to Cats Social API",
		})
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http_helper.EncodeJSON(w, http.StatusNotFound, http_helper.ResponseBody{
			Error:   "not found",
			Message: "route does not exist",
		})
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		http_helper.EncodeJSON(w, http.StatusMethodNotAllowed, http_helper.ResponseBody{
			Error:   "method not allowed",
			Message: "method is not allowed",
		})
	})

	return r
}
