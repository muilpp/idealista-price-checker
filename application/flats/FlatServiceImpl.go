package application

import (
	"encoding/json"
	"idealista/application/authentication"
	"idealista/domain"
	"idealista/domain/ports"
	"idealista/infrastructure/persistance"
	"strconv"

	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const (
	flatEndpoint   = "https://api.idealista.com/3.5/es/search"
	geolocation    = "41.6053142,2.2851149"
	marginInMeters = "4000"
	propertyType   = "homes"
	country        = "es"
	rentType       = "rent"
	saleType       = "sale"
)

var (
	smallFlatSize = domain.NewFlatSize(75, 90)
	bigFlatSize   = domain.NewFlatSize(90, 110)
)

type flatServiceImpl struct {
	flatRepository ports.FlatRepository
	authentication authentication.AuthenticationService
}

func NewFlatService() ports.FlatService {
	return &flatServiceImpl{persistance.NewFlatRepository(), authentication.NewAuthenticationService()}
}

func (f flatServiceImpl) AddNewFlats() bool {

	rentalSmallFlatsSlice := f.getFlatsFromIdealista(rentType, smallFlatSize)
	saleSmallFlatsSlice := f.getFlatsFromIdealista(saleType, smallFlatSize)

	rentalBigFlatsSlice := f.getFlatsFromIdealista(rentType, bigFlatSize)
	saleBigFlatsSlice := f.getFlatsFromIdealista(saleType, bigFlatSize)

	f.flatRepository.Add(rentalSmallFlatsSlice, rentType, smallFlatSize.GetMinSize())
	f.flatRepository.Add(saleSmallFlatsSlice, saleType, smallFlatSize.GetMinSize())

	f.flatRepository.Add(rentalBigFlatsSlice, rentType, bigFlatSize.GetMinSize())
	f.flatRepository.Add(saleBigFlatsSlice, saleType, bigFlatSize.GetMinSize())

	return true
}

func (f flatServiceImpl) getFlatsFromIdealista(operation string, flatSize *domain.FlatSize) []domain.Flat {
	data := url.Values{}
	data.Set("country", country)
	data.Set("operation", operation)
	data.Set("propertyType", propertyType)
	data.Set("center", geolocation)
	data.Set("distance", marginInMeters)
	data.Set("minSize", strconv.Itoa(flatSize.GetMinSize()))
	data.Set("maxSize", strconv.Itoa(flatSize.GetMaxSize()))

	if flatSize.GetMaxSize() == bigFlatSize.GetMaxSize() {
		data.Set("bedrooms", "3")
		data.Set("terrance", "1")
		data.Set("maxItems", "50")
	}

	var bearer = "Bearer " + f.authentication.GetToken()

	req, _ := http.NewRequest("POST", flatEndpoint, strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	log.Println(string([]byte(body)))

	var flatList domain.FlatList
	var flats []domain.Flat
	err2 := json.Unmarshal([]byte(body), &flatList)

	if err2 != nil {
		log.Println(err2)
	} else {
		for _, s := range flatList.Flats {
			flat := domain.NewFlat(s.Price, s.AreaPrice)
			flats = append(flats, *flat)
		}
	}

	return flats
}

func (f flatServiceImpl) GetFlatsFromDatabase(operation string, oncePerMonth bool, isFormatDate bool) ([]domain.Flat, []domain.Flat) {
	log.Println("Get flats for operation ", operation)

	return f.flatRepository.Get(operation, oncePerMonth, isFormatDate, smallFlatSize.GetMinSize()), f.flatRepository.Get(operation, oncePerMonth, isFormatDate, bigFlatSize.GetMinSize())
}
