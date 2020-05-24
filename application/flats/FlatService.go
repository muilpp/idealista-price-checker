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

type flatService interface {
	GetAllFlats() bool
}

type flatServiceImpl struct {
	flatRepository persistance.FlatRepository
	authentication authentication.AuthenticationService
}

func NewFlatService() flatService {
	return &flatServiceImpl{persistance.NewFlatRepository(), authentication.NewAuthenticationService()}
}

func (f flatServiceImpl) GetAllFlats() bool {
	rentalFlatsSlice := f.getFlats("rent")
	saleFlatsSlice := f.getFlats("sale")

	f.flatRepository.Add(rentalFlatsSlice, "rent")
	f.flatRepository.Add(saleFlatsSlice, "sale")

	return true
}

func (f flatServiceImpl) getFlats(operation string) []domain.Flat {
	data := url.Values{}
	data.Set("country", "es")
	data.Set("operation", operation)
	data.Set("propertyType", "homes")
	data.Set("center", "41.6061846,2.2703413")
	data.Set("distance", "4000")
	data.Set("minSize", "75")
	data.Set("maxSize", "90")

	var bearer = "Bearer " + f.authentication.GetToken()

	req, err := http.NewRequest("POST", "https://api.idealista.com/3.5/es/search", strings.NewReader(data.Encode()))
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
