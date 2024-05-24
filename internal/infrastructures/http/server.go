package http

import (
	"log"
	"net/http"
	"os"

	"github.com/danzbraham/cats-social/internal/applications/usecases"
	http_common "github.com/danzbraham/cats-social/internal/commons/http"
	repositories_impl "github.com/danzbraham/cats-social/internal/infrastructures/repositories"
	securities_impl "github.com/danzbraham/cats-social/internal/infrastructures/securities"
	"github.com/danzbraham/cats-social/internal/interfaces/http/controllers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
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

	// Helpers
	passwordHasher := securities_impl.NewBcryptPasswordHasher()
	authTokenManager := securities_impl.NewJWTTokenManager([]byte(os.Getenv("JWT_SECRET")))
	validator := securities_impl.NewGoValidator(validator.New(validator.WithRequiredStructEnabled()))

	// User domain
	userRepository := repositories_impl.NewUserRepositoryPostgres(s.DB)
	userUsecase := usecases.NewUserUsecase(userRepository, passwordHasher, authTokenManager)
	userController := controllers.NewUserController(userUsecase, validator)

	// Cat domain
	catRepository := repositories_impl.NewCatRepositoryImpl(s.DB)
	catUsecase := usecases.NewCatUsecase(catRepository)
	catController := controllers.NewCatController(catUsecase, validator, authTokenManager)

	r.Route("/v1", func(r chi.Router) {
		r.Mount("/user", userController.Routes())
		r.Mount("/cat", catController.Routes())
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
