package controllers

import (
	"net/http"

	service "idealista/application/flats"

	"github.com/gin-gonic/gin"
)

func AddFlat(c *gin.Context) {
	flatService := service.NewFlatService()
	flatService.AddNewFlats()
	c.JSON(http.StatusOK, true)
}

func GetRentalFlats(c *gin.Context) {
	flatService := service.NewFlatService()
	rentalFlats, _ := flatService.GetFlatsFromDatabase("rent", false, false)
	c.JSON(http.StatusOK, rentalFlats)
}

func GetRentalFlatsOncePerMonth(c *gin.Context) {
	flatService := service.NewFlatService()
	rentalFlats, _ := flatService.GetFlatsFromDatabase("rent", true, false)
	c.JSON(http.StatusOK, rentalFlats)
}

func GetSaleFlats(c *gin.Context) {
	flatService := service.NewFlatService()
	saleFlats, _ := flatService.GetFlatsFromDatabase("sale", false, false)
	c.JSON(http.StatusOK, saleFlats)
}

func GetSaleFlatsOncePerMonth(c *gin.Context) {
	flatService := service.NewFlatService()
	saleFlats, _ := flatService.GetFlatsFromDatabase("sale", true, false)
	c.JSON(http.StatusOK, saleFlats)
}
