package http

import (
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	Addr string
	DB   *pgxpool.Pool
}

func NewServer(addr string, db *pgxpool.Pool) *Server {
	return &Server{
		Addr: addr,
		DB:   db,
	}
}

func (s *Server) Launch() error {
	server := &http.Server{
		Addr:    s.Addr,
		Handler: s.RegisterRoutes(),
	}

	log.Printf("Server listening on %s\n", server.Addr)
	return server.ListenAndServe()
}
