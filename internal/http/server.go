package http

import (
	"log"
	"net/http"

	http_common "github.com/danzbraham/cats-social/internal/commons/http"
	"github.com/danzbraham/cats-social/internal/commons/validator"
	"github.com/danzbraham/cats-social/internal/http/controllers"
	"github.com/danzbraham/cats-social/internal/repositories"
	"github.com/danzbraham/cats-social/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

type APIServer struct {
	Addr string
	DB   *pgxpool.Pool
}

func NewAPIServer(addr string, db *pgxpool.Pool) *APIServer {
	return &APIServer{
		Addr: addr,
		DB:   db,
	}
}

func (s *APIServer) Launch() error {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Cats Social API"))
	})

	validator.InitCustomValidation()

	// User domain
	userRepository := repositories.NewUserRepository(s.DB)
	userService := services.NewUserService(userRepository)
	userController := controllers.NewUserController(userService)

	// Cat domain
	catRepository := repositories.NewCatRepository(s.DB)
	catService := services.NewCatService(catRepository)
	catController := controllers.NewCatController(catService)

	// Match domain
	matchRepository := repositories.NewMatchRepository(s.DB)
	matchService := services.NewMatchService(catRepository, matchRepository)
	matchController := controllers.NewMatchController(matchService)

	r.Route("/v1", func(r chi.Router) {
		r.Mount("/user", userController.Routes())
		r.Mount("/cat", catController.Routes())
		r.Mount("/match", matchController.Routes())
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http_common.ResponseError(w, http.StatusNotFound, "Not found error", "Route does not exists")
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		http_common.ResponseError(w, http.StatusMethodNotAllowed, "Method not allowed error", "Method is not allowed")
	})

	server := http.Server{
		Addr:    s.Addr,
		Handler: r,
	}

	log.Printf("Server listening on %s\n", s.Addr)
	return server.ListenAndServe()
}
