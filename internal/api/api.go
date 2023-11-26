package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"pokedex/internal/config"
)

type locationResult struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type LocationsData struct {
	Count    int              `json:"count"`
	Next     string           `json:"next"`
	Previous string           `json:"previous"`
	Results  []locationResult `json:"results"`
}

func GetLocations(url string, data *LocationsData, config *config.Config) error {
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(response.Body)
	response.Body.Close()
	if response.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and \nbody: %s\n", response.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(body, data)
	if err != nil {
		log.Fatal(err)
	}
	config.PrevLocationUrl = data.Previous
	config.NextLocationUrl = data.Next
	for _, result := range data.Results {
		fmt.Println(result.Name)
	}
	return nil
}
