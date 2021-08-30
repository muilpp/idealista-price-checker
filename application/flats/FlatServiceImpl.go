package application

import (
	"encoding/json"
	"idealista/adapter/out/persistance"
	"idealista/application/authentication"
	"idealista/domain"
	"idealista/domain/ports"
	"strconv"

	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const (
	flatEndpoint = "https://api.idealista.com/3.5/es/search"
	propertyType = "homes"
	country      = "es"
	maxItems     = "50"
)

type flatServiceImpl struct {
	flatRepository ports.FlatRepository
	authentication authentication.AuthenticationService
}

func NewFlatService() ports.FlatService {
	return &flatServiceImpl{persistance.NewFlatRepository(), authentication.NewAuthenticationService()}
}

func (f flatServiceImpl) AddNewFlats() bool {

	flatSlice := f.getFlatsFromIdealista()
	f.flatRepository.Add(flatSlice)

	return true
}

func (f flatServiceImpl) getFlatsFromIdealista() [][]domain.Flat {
	placesSlice := f.getPlacesToSearch()
	var allFlats [][]domain.Flat

	var bearer = "Bearer " + f.authentication.GetToken()

	for _, place := range placesSlice {
		data := url.Values{}
		data.Set("country", country)
		data.Set("operation", place.GetOperation())
		data.Set("propertyType", propertyType)
		data.Set("center", place.GetCenter())
		data.Set("distance", strconv.Itoa(place.GetDistance()))
		data.Set("minSize", strconv.Itoa(place.GetMinSize()))
		data.Set("maxSize", strconv.Itoa(place.GetMaxSize()))
		data.Set("bedrooms", place.GetBedrooms())
		data.Set("terrance", strconv.FormatBool(place.HasTerrace()))
		data.Set("maxItems", maxItems)

		req, _ := http.NewRequest("POST", flatEndpoint, strings.NewReader(data.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Authorization", bearer)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println("Error on response.\n[ERROR] -", err)
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
				flat := domain.NewFlat(place.GetId(), s.Price, s.AreaPrice, *domain.NewFlatSize(place.GetMinSize(), place.GetMaxSize()))
				flats = append(flats, *flat)
			}
		}

		allFlats = append(allFlats, flats)
	}

	return allFlats
}

func (f flatServiceImpl) getPlacesToSearch() []domain.Place {
	return f.flatRepository.GetPlacesToSearch()
}

func (f flatServiceImpl) GetFlatsFromDatabase(operation string, oncePerMonth bool) [][]domain.Flat {
	log.Println("Get flats for operation ", operation)

	return f.flatRepository.Get(operation, oncePerMonth)
}
