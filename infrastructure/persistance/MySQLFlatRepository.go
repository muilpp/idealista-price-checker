package persistance

import (
	"database/sql"
	"fmt"
	"idealista/domain"
	"idealista/domain/ports"
	"log"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type mysqlFlatRepository struct{}

func NewFlatRepository() ports.FlatRepository {
	return &mysqlFlatRepository{}
}

func (f mysqlFlatRepository) Add(flats []domain.Flat, operation string, flatSize int) bool {
	db := openDB()
	defer db.Close()

	var totalSumPrice, totalSumAreaPrice float64
	for _, flat := range flats {
		totalSumPrice += flat.Price
		totalSumAreaPrice += flat.AreaPrice
	}

	averagePrice := totalSumPrice / float64(len(flats))
	averageAreaPrice := totalSumAreaPrice / float64(len(flats))

	insert, err := db.Query("INSERT INTO " + operation + "_average_price (average, area_average, size) VALUES ('" + fmt.Sprintf("%f", averagePrice) + "', '" + fmt.Sprintf("%f", averageAreaPrice) + "', '" + strconv.Itoa(flatSize) + "')")
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()

	return true
}

func (f mysqlFlatRepository) Get(operation string, getOncePerMonthOnly bool, isFormatDate bool, flatSize int) []domain.Flat {
	db := openDB()
	defer db.Close()

	var query string
	if isFormatDate {
		query = "select average, area_average, size, DATE_FORMAT(added,'%b %y') from " + operation + "_average_price where size = '" + strconv.Itoa(flatSize) + "'"
	} else {
		query = "select average, area_average, size, added from " + operation + "_average_price where size = '" + strconv.Itoa(flatSize) + "'"
	}

	if getOncePerMonthOnly {
		query += " and DAY(added) = 1"
	}

	log.Println("Get flats -> ", query)
	rows, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}

	defer rows.Close()

	var average float64
	var areaAverage float64
	var added string
	var size string
	var flats []domain.Flat

	for rows.Next() {
		err := rows.Scan(&average, &areaAverage, &size, &added)
		if err != nil {
			log.Fatal(err)
		}

		flat := domain.NewFlatWithDate(average, areaAverage, added)
		flats = append(flats, *flat)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	if len(flats) > 0 {
		return flats
	}

	return nil
}

func openDB() *sql.DB {
	db, err := sql.Open("mysql", os.Getenv("DB_DATA_SOURCE"))

	if err != nil {
		panic(err.Error())
	}

	return db
}
