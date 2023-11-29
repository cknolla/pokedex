package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"pokedex/internal/config"
	"strings"
)

type locationResult struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type locationsResponse struct {
	Count    int              `json:"count"`
	Next     string           `json:"next"`
	Previous string           `json:"previous"`
	Results  []locationResult `json:"results"`
}

type locationResponse struct {
	Id                int                 `json:"id"`
	Name              string              `json:"name"`
	PokemonEncounters []pokemonEncounters `json:"pokemon_encounters"`
}

type locationPokemonResult struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type pokemonEncounters struct {
	Pokemon locationPokemonResult `json:"pokemon"`
}

type pokemonResponse struct {
	BaseExperience int `json:"base_experience"`
}

func getResponseBody(path string, config *config.Config) ([]byte, error) {
	if path == "" {
		return nil, errors.New("path must not be empty")
	}
	body, found := config.Cache.Get(path)
	if !found {
		//log.Println("path not found in cache")
		response, err := http.Get(config.ApiRoot + path)
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
		config.Cache.Add(path, body)
	}
	return body, nil
}

func GetLocations(path string, config *config.Config) ([]locationResult, error) {
	body, err := getResponseBody(path, config)
	if err != nil {
		return nil, err
	}
	var data locationsResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	config.PrevLocationUrl = strings.ReplaceAll(data.Previous, config.ApiRoot, "")
	config.NextLocationUrl = strings.ReplaceAll(data.Next, config.ApiRoot, "")
	return data.Results, nil
}

func GetLocationDetails(path string, config *config.Config) ([]locationPokemonResult, error) {
	body, err := getResponseBody(path, config)
	if err != nil {
		return nil, err
	}
	var data locationResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	var pokemons []locationPokemonResult
	for _, encounter := range data.PokemonEncounters {
		pokemons = append(pokemons, encounter.Pokemon)
	}
	return pokemons, nil
}

func CatchPokemon(name string, config *config.Config) (bool, error) {
	body, err := getResponseBody("/pokemon/"+name, config)
	if err != nil {
		if strings.Contains(err.Error(), "status code: 404") {
			err = errors.New("pokemon name not found")
		}
		return false, err
	}
	var data pokemonResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return false, err
	}
	threshold := config.Generator.Float32()
	difficulty := 1 / float32(data.BaseExperience) * 100
	if difficulty > threshold {
		return true, nil
	}
	return false, nil
}
