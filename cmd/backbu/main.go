package main

import (
	"backbu/internal/api"
	"log"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("../../.env"); err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

	router := api.Router()
	go router.Run("localhost:8080")
}
