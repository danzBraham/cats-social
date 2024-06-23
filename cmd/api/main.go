package main

import (
	"log"
	"os"

	"github.com/danzBraham/cats-social/internal/database"
	"github.com/danzBraham/cats-social/internal/http"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading .env file")
	}

	addr := os.Getenv("APP_HOST") + ":" + os.Getenv("APP_PORT")
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
