package main

import (
	service "idealista/application/flats"
	"idealista/infrastructure"
	"idealista/infrastructure/controllers"
	"idealista/infrastructure/notification"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("credentials.env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	if len(os.Args) < 2 {
		r := gin.Default()
		r.Use(cors.Default())
		r.GET("/idealista/add", controllers.AddFlat)
		r.GET("/idealista/get-rental-flats", controllers.GetRentalFlats)
		r.GET("/idealista/get-rental-flats/once-per-month", controllers.GetRentalFlatsOncePerMonth)
		r.GET("/idealista/get-sale-flats", controllers.GetSaleFlats)
		r.GET("/idealista/get-sale-flats/once-per-month", controllers.GetSaleFlatsOncePerMonth)
		r.Run(":8383")
	} else {
		executionType := os.Args[1]

		if executionType == "sendMonthlyReports" {
			flatService := service.NewFlatService()
			reportsService := infrastructure.NewReportsService()
			reportsService.GetMonthlyRentalReports(flatService.GetFlatsFromDatabase("rent", true, true))
			reportsService.GetMonthlySaleReports(flatService.GetFlatsFromDatabase("sale", true, true))

			telegramService := notification.NewTelegramNotification()
			telegramService.SendReports()
		}
	}
}
