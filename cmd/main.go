package main

import (
	service "idealista/application/flats"
	"idealista/infrastructure"
	"idealista/infrastructure/notification"
	"log"
	"net/http"
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

	flatService := service.NewFlatService()

	if len(os.Args) < 2 {
		r := gin.Default()
		r.Use(cors.Default())
		r.GET("/idealista/add", func(c *gin.Context) {
			flatService.AddNewFlats()
			c.JSON(http.StatusOK, true)
		})

		r.GET("/idealista/get-rental-flats", func(c *gin.Context) {
			rentalFlats := flatService.GetFlatsFromDatabase("rent", false, false)
			c.JSON(http.StatusOK, rentalFlats)
		})

		r.GET("/idealista/get-rental-flats/once-per-month", func(c *gin.Context) {
			rentalFlats := flatService.GetFlatsFromDatabase("rent", true, false)
			c.JSON(http.StatusOK, rentalFlats)
		})

		r.GET("/idealista/get-sale-flats", func(c *gin.Context) {
			saleFlats := flatService.GetFlatsFromDatabase("sale", false, false)
			c.JSON(http.StatusOK, saleFlats)
		})

		r.GET("/idealista/get-sale-flats/once-per-month", func(c *gin.Context) {
			saleFlats := flatService.GetFlatsFromDatabase("sale", true, false)
			c.JSON(http.StatusOK, saleFlats)
		})

		r.Run(":8383")
	} else {
		executionType := os.Args[1]

		if executionType == "sendMonthlyReports" {
			reportsService := infrastructure.NewReportsService()
			reportsService.GetMonthlyRentalReports(flatService.GetFlatsFromDatabase("rent", true, true))
			reportsService.GetMonthlySaleReports(flatService.GetFlatsFromDatabase("sale", true, true))

			telegramService := notification.NewTelegramNotification()
			telegramService.SendReports()
		}
	}
}
