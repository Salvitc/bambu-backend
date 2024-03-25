package main

import (
	"backbu/internal/api"
	"backbu/pkg/database"

	"github.com/joho/godotenv"
)

func main() {
	//Carga las variables de entorno en desarrollo
	godotenv.Load(".env")

	db.GetClient()

	router := api.Router()
	router.Run("localhost:8080")
}
