package main

import (
	"log"
	"os"

	"github.com/danzbraham/cats-social/db"
	"github.com/danzbraham/cats-social/internal/http"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	pool, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer pool.Close()

	addr := os.Getenv("APP_HOST") + ":" + os.Getenv("APP_PORT")
	server := http.NewAPIServer(addr, pool)
	if err := server.Launch(); err != nil {
		log.Fatal(err)
	}
}
