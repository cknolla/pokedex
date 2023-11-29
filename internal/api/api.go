package api

import (
	"encoding/json"
	"errors"
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

type PokemonResponse struct {
	Id                int                 `json:"id"`
	Name              string              `json:"name"`
	PokemonEncounters []pokemonEncounters `json:"pokemon_encounters"`
}

type pokemonResult struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type pokemonEncounters struct {
	Pokemon pokemonResult `json:"pokemon"`
}

func getResponseBody(url string, config *config.Config) ([]byte, error) {
	body, found := config.Cache.Get(url)
	if !found {
		log.Println("url not found in cache")
		response, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		body, err = io.ReadAll(response.Body)
		response.Body.Close()
		if response.StatusCode > 299 {
			return nil, errors.New(
				fmt.Sprintf("Response failed with status code: %d and \nbody: %s\n", response.StatusCode, body),
			)
		}
		if err != nil {
			return nil, err
		}
		config.Cache.Add(url, body)
	}
	return body, nil
}

func GetLocations(url string, data *LocationsData, config *config.Config) ([]locationResult, error) {
	body, err := getResponseBody(url, config)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, data)
	if err != nil {
		return nil, err
	}
	config.PrevLocationUrl = data.Previous
	config.NextLocationUrl = data.Next
	return data.Results, nil
}

func GetLocationDetails(url string, data *PokemonResponse, config *config.Config) ([]pokemonResult, error) {
	body, err := getResponseBody(url, config)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, data)
	if err != nil {
		return nil, err
	}
	var pokemons []pokemonResult
	for _, encounter := range data.PokemonEncounters {
		pokemons = append(pokemons, encounter.Pokemon)
	}
	return pokemons, nil
}
