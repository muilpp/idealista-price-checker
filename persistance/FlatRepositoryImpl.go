package persistance

import (
	"database/sql"
	"fmt"
	"idealista/domain"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type FlatRepository interface {
	Add([]domain.Flat, string) bool
}

type flatRepositoryImpl struct{}

func NewFlatRepository() FlatRepository {
	return &flatRepositoryImpl{}
}

func (f flatRepositoryImpl) Add(flats []domain.Flat, operation string) bool {
	db, err := sql.Open("mysql", os.Getenv("DB_DATA_SOURCE"))

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	var totalSumPrice, totalSumAreaPrice float64
	for _, flat := range flats {
		insert, err := db.Query("INSERT INTO " + operation + "_flat_prices (price, area_price) VALUES ('" + fmt.Sprintf("%f", flat.Price) + "', '" + fmt.Sprintf("%f", flat.AreaPrice) + "')")
		if err != nil {
			panic(err.Error())
		}
		defer insert.Close()
		totalSumPrice += flat.Price
		totalSumAreaPrice += flat.AreaPrice
	}

	averagePrice := totalSumPrice / float64(len(flats))
	averageAreaPrice := totalSumAreaPrice / float64(len(flats))
	insert, err := db.Query("INSERT INTO " + operation + "_average_price (average, area_average) VALUES ('" + fmt.Sprintf("%f", averagePrice) + "', '" + fmt.Sprintf("%f", averageAreaPrice) + "')")
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()

	return true
}
