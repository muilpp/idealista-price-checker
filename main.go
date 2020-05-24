package main

import (
	service "idealista/application/flats"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("credentials.env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	flatService := service.NewFlatService()
	flatService.GetAllFlats()
}
