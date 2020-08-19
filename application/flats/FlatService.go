package application

import (
	"encoding/json"
	"idealista/application/authentication"
	"idealista/domain"
	"idealista/persistance"

	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const (
	flatEndpoint   = "https://api.idealista.com/3.5/es/search"
	geolocation    = "41.6061846,2.2703413"
	marginInMeters = "4000"
	flatMinSize    = "75"
	flatMaxSize    = "90"
	propertyType   = "homes"
	country        = "es"
	rentType       = "rent"
	saleType       = "sale"
)

type flatService interface {
	AddNewFlats() bool
	GetFlatsFromDatabase(string) []domain.Flat
}

type flatServiceImpl struct {
	flatRepository persistance.FlatRepository
	authentication authentication.AuthenticationService
}

func NewFlatService() flatService {
	return &flatServiceImpl{persistance.NewFlatRepository(), authentication.NewAuthenticationService()}
}

func (f flatServiceImpl) AddNewFlats() bool {
	rentalFlatsSlice := f.getFlatsFromIdealista(rentType)
	saleFlatsSlice := f.getFlatsFromIdealista(saleType)

	f.flatRepository.Add(rentalFlatsSlice, rentType)
	f.flatRepository.Add(saleFlatsSlice, saleType)

	return true
}

func (f flatServiceImpl) getFlatsFromIdealista(operation string) []domain.Flat {
	data := url.Values{}
	data.Set("country", country)
	data.Set("operation", operation)
	data.Set("propertyType", propertyType)
	data.Set("center", geolocation)
	data.Set("distance", marginInMeters)
	data.Set("minSize", flatMinSize)
	data.Set("maxSize", flatMaxSize)

	var bearer = "Bearer " + f.authentication.GetToken()

	req, err := http.NewRequest("POST", flatEndpoint, strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	// log.Println(string([]byte(body)))

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

func (f flatServiceImpl) GetFlatsFromDatabase(operation string) []domain.Flat {
	log.Println("Get flats for operation ", operation)
	flats := f.flatRepository.Get(operation)

	return flats
}
