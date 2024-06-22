package main

import (
	"log"

	"github.com/danzBraham/cats-social/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	server := server.NewServer()

	log.Printf("Server listening on %s\n", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
