package controllers

import (
	"net/http"

	"idealista/adapter/out/persistance"
	"idealista/application/authentication"
	service "idealista/application/flats"

	"github.com/gin-gonic/gin"
)

var flatService = service.NewFlatService(persistance.NewFlatRepository(), authentication.NewAuthenticationService())

func AddFlat(c *gin.Context) {
	flatService.AddNewFlats()
	c.JSON(http.StatusOK, true)
}

func GetRentalFlats(c *gin.Context) {
	rentalFlats := flatService.GetFlatsFromDatabase("rent", false)
	c.JSON(http.StatusOK, rentalFlats)
}

func GetRentalFlatsOncePerMonth(c *gin.Context) {
	rentalFlats := flatService.GetFlatsFromDatabase("rent", true)
	c.JSON(http.StatusOK, rentalFlats)
}

func GetSaleFlats(c *gin.Context) {
	saleFlats := flatService.GetFlatsFromDatabase("sale", false)
	c.JSON(http.StatusOK, saleFlats)
}

func GetSaleFlatsOncePerMonth(c *gin.Context) {
	saleFlats := flatService.GetFlatsFromDatabase("sale", true)
	c.JSON(http.StatusOK, saleFlats)
}
