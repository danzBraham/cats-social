package http

import (
	"net/http"

	"github.com/danzBraham/cats-social/internal/helpers/http_helper"
	"github.com/danzBraham/cats-social/internal/http/controllers"
	"github.com/danzBraham/cats-social/internal/http/middlewares"
	"github.com/danzBraham/cats-social/internal/repositories"
	"github.com/danzBraham/cats-social/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http_helper.EncodeJSON(w, http.StatusOK, http_helper.ResponseBody{
			Message: "Welcome to Cats Social API",
		})
	})

	// user domain
	userRepository := repositories.NewUserRepository(s.DB)
	userService := services.NewUserService(userRepository)
	userController := controllers.NewUserController(userService)

	r.Route("/v1", func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {
			r.Post("/register", userController.HandleRegisterUser)
			r.Post("/login", userController.HandleLoginUser)
		})

		r.Route("/cat", func(r chi.Router) {
			r.Use(middlewares.Auth)
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
