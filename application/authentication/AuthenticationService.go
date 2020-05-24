package authentication

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type AuthenticationService interface {
	GetToken() string
}

func NewAuthenticationService() AuthenticationService {
	return &authentication{}
}

type authentication struct {
	Token string `json:"access_token"`
}

func (a authentication) GetToken() string {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("scope", "read")

	req, err := http.NewRequest("POST", "https://api.idealista.com/oauth/token", strings.NewReader(data.Encode()))
	req.SetBasicAuth(os.Getenv("AUTH_USER"), os.Getenv("AUTH_PWD"))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		panic(err)
	}
	// fmt.Println(req.Header.Get("Authorization"))

	if err != nil {
		log.Fatalln(err)
	}

	client := &http.Client{}
	resp, _ := client.Do(req)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println("response Body:", string(body))

	err2 := json.Unmarshal(body, &a)
	if err2 != nil {
		log.Println(err2)
	}

	return a.Token
}
