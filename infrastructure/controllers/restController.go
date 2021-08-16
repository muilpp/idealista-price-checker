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
	rentalFlats := flatService.GetFlatsFromDatabase("rent", false)
	c.JSON(http.StatusOK, rentalFlats)
}

func GetRentalFlatsOncePerMonth(c *gin.Context) {
	flatService := service.NewFlatService()
	rentalFlats := flatService.GetFlatsFromDatabase("rent", true)
	c.JSON(http.StatusOK, rentalFlats)
}

func GetSaleFlats(c *gin.Context) {
	flatService := service.NewFlatService()
	saleFlats := flatService.GetFlatsFromDatabase("sale", false)
	c.JSON(http.StatusOK, saleFlats)
}

func GetSaleFlatsOncePerMonth(c *gin.Context) {
	flatService := service.NewFlatService()
	saleFlats := flatService.GetFlatsFromDatabase("sale", true)
	c.JSON(http.StatusOK, saleFlats)
}
