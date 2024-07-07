package main

import (
	"log"

	"github.com/danzBraham/cats-social/internal/database"
	"github.com/danzBraham/cats-social/internal/http"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	addr := ":8080"
	pool, err := database.Connect()
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}
	defer pool.Close()

	server := http.NewServer(addr, pool)
	if err := server.Launch(); err != nil {
		log.Fatal(err)
	}
}
