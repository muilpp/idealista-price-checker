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

func (f mysqlFlatRepository) Add(allFlats [][]domain.Flat) bool {
	db := openDB()
	defer db.Close()

	var totalSumPrice, totalSumAreaPrice float64
	for _, flatSlice := range allFlats {

		for _, flat := range flatSlice {
			totalSumPrice += flat.Price
			totalSumAreaPrice += flat.AreaPrice
		}

		averagePrice := totalSumPrice / float64(len(flatSlice))
		averageAreaPrice := totalSumAreaPrice / float64(len(flatSlice))

		insert, err := db.Query("INSERT INTO sale_average_price (place_id, average, area_average, min_size, max_size) VALUES ('" + strconv.Itoa(flatSlice[0].AreaId) + "', '" + fmt.Sprintf("%f", averagePrice) + "', '" + fmt.Sprintf("%f", averageAreaPrice) + "', '" + strconv.Itoa(flatSlice[0].Size.GetMinSize()) + "', '" + strconv.Itoa(flatSlice[0].Size.GetMaxSize()) + "'")
		if err != nil {
			panic(err.Error())
		}
		defer insert.Close()
	}

	return true
}

func (f mysqlFlatRepository) Get(operation string, getOncePerMonthOnly bool) [][]domain.Flat {
	db := openDB()
	defer db.Close()

	var allFlats [][]domain.Flat
	places := f.GetPlacesToSearch()

	for _, place := range places {
		var query string
		query = "select p.name, s.average, s.area_average, s.min_size, s.max_size, DATE_FORMAT(s.added,'%b %y') "
		query += "FROM " + operation + "_average_price s, places_to_search p "
		query += "WHERE s.place_id = p.id "
		query += "AND p.id = '" + strconv.Itoa(place.GetId()) + "' "

		//query = "select average, area_average, size, DATE_FORMAT(added,'%b %y') from " + operation + "_average_price where size = '" + strconv.Itoa(flatSize) + "'"

		if getOncePerMonthOnly {
			query += " and DAY(added) = 1"
		}

		log.Println("Get flats -> ", query)
		rows, err := db.Query(query)
		if err != nil {
			panic(err.Error())
		}

		defer rows.Close()

		var name string
		var average float64
		var areaAverage float64
		var added string
		var size string
		var flats []domain.Flat

		for rows.Next() {
			err := rows.Scan(&name, &average, &areaAverage, &size, &added)
			if err != nil {
				log.Fatal(err)
			}

			flat := domain.NewFlatWithDate(name, average, areaAverage, added)
			flats = append(flats, *flat)
		}

		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}

		if len(flats) > 0 {
			allFlats = append(allFlats, flats)
		}
	}

	return allFlats
}

func (f mysqlFlatRepository) GetPlacesToSearch() []domain.Place {
	db := openDB()
	defer db.Close()

	query := "SELECT id, name, center, distance, min_size, max_size, bedrooms, terrace, operation FROM idealista.places_to_search"

	log.Println("Get places to search -> ", query)
	rows, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}

	defer rows.Close()

	var id int
	var name string
	var center string
	var distance int
	var minSize int
	var maxSize int
	var bedrooms string
	var terrace bool
	var operation string
	var placeSlice []domain.Place

	for rows.Next() {
		err := rows.Scan(&id, &name, &center, &distance, &minSize, &maxSize, &bedrooms, &terrace, &operation)
		if err != nil {
			log.Fatal(err)
		}

		placeToSearch := domain.NewPlace(id, name, center, distance, minSize, maxSize, bedrooms, terrace, operation)
		placeSlice = append(placeSlice, *placeToSearch)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	if len(placeSlice) > 0 {
		return placeSlice
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
