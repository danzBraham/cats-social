package server

import (
	"log"
	"net/http"
	"os"

	"github.com/danzBraham/cats-social/internal/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	Addr string
	DB   *pgxpool.Pool
}

func NewServer() *http.Server {
	addr := os.Getenv("APP_HOST") + ":" + os.Getenv("APP_PORT")
	pool, err := database.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer pool.Close()

	NewServer := &Server{
		Addr: addr,
		DB:   pool,
	}

	server := &http.Server{
		Addr:    NewServer.Addr,
		Handler: NewServer.RegisterRoutes(),
	}

	return server
}
