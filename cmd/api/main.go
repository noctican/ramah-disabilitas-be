package main

import (
	"log"
	"ramah-disabilitas-be/internal/router"
	"ramah-disabilitas-be/pkg/database"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, finding environment variables from system")
	}

	database.Connect()

	r := router.SetupRouter()

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
