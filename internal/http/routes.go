package http

import (
	"net/http"

	"github.com/danzBraham/cats-social/internal/errors/commonerror"
	"github.com/danzBraham/cats-social/internal/helpers/httphelper"
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
		httphelper.EncodeJSON(w, http.StatusOK, httphelper.ResponseBody{
			Message: "Welcome to Cats Social API",
		})
	})

	// repositories
	userRepository := repositories.NewUserRepository(s.DB)
	catRepository := repositories.NewCatRepository(s.DB)
	matchRepository := repositories.NewMatchRepository(s.DB)

	// services
	userService := services.NewUserService(userRepository)
	catService := services.NewCatService(catRepository, matchRepository)
	matchService := services.NewMatchService(matchRepository, catRepository, userRepository)

	// controllers
	userController := controllers.NewUserController(userService)
	catController := controllers.NewCatController(catService)
	matchController := controllers.NewMatchController(matchService)

	r.Route("/v1", func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {
			r.Post("/register", userController.HandleRegisterUser)
			r.Post("/login", userController.HandleLoginUser)
		})

		r.Group(func(r chi.Router) {
			r.Use(middlewares.Auth)

			r.Route("/cat", func(r chi.Router) {
				r.Post("/", catController.HandleCreateCat)
				r.Get("/", catController.HandleGetCats)
				r.Put("/{id}", catController.HandleUpdateCatById)
				r.Delete("/{id}", catController.HandleDeleteCatById)

				r.Route("/match", func(r chi.Router) {
					r.Post("/", matchController.HandleCreateMatch)
					r.Get("/", matchController.HandleGetMatches)
					r.Post("/approve", matchController.HandleApproveMatch)
					r.Post("/reject", matchController.HandleRejectMatch)
					r.Delete("/{id}", matchController.HandleDeleteMatch)
				})
			})
		})
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		httphelper.ErrorResponse(w, http.StatusNotFound, commonerror.ErrRouteDoesNotExist)
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		httphelper.ErrorResponse(w, http.StatusMethodNotAllowed, commonerror.ErrMethodNotAllowed)
	})

	return r
}
