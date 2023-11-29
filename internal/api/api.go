package api

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"pokedex/internal/config"
	"strings"
)

type namedApiResource struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

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

type pokemonEncounters struct {
	Pokemon namedApiResource `json:"pokemon"`
}

type pokemonStat struct {
	Stat     namedApiResource `json:"stat"`
	BaseStat int              `json:"base_stat"`
	Effort   int              `json:"effort"`
}

type pokemonType struct {
	Slot int              `json:"slot"`
	Type namedApiResource `json:"type"`
}

type Pokemon struct {
	Name           string        `json:"name"`
	BaseExperience int           `json:"base_experience"`
	Height         int           `json:"height"`
	Weight         int           `json:"weight"`
	Stats          []pokemonStat `json:"stats"`
	Types          []pokemonType `json:"types"`
}

func decodePokemon(pokemonBytes []byte, pokemon *Pokemon) error {
	buffer := bytes.NewBuffer(pokemonBytes)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(&pokemon); err != nil {
		return err
	}
	return nil
}

func encodePokemon(pokemon *Pokemon) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(pokemon); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// getResponseBody performs the generic http request and caches the result
func getResponseBody(path string, cfg *config.Config) ([]byte, int, error) {
	if path == "" {
		return nil, 0, errors.New("path must not be empty")
	}
	body, found := cfg.Cache.Get(path)
	if !found {
		//log.Println("path not found in cache")
		response, err := http.Get(cfg.ApiRoot + path)
		if err != nil {
			return nil, response.StatusCode, err
		}
		body, err = io.ReadAll(response.Body)
		response.Body.Close()
		if response.StatusCode > 299 {
			return nil, response.StatusCode, errors.New(
				fmt.Sprintf("Response failed with status code: %d and \nbody: %s\n", response.StatusCode, body),
			)
		}
		if err != nil {
			return nil, response.StatusCode, err
		}
		cfg.Cache.Add(path, body, false)
	}
	return body, http.StatusOK, nil
}

// GetLocations returns a paginated list of all locations
func GetLocations(queryString string, cfg *config.Config) ([]locationResult, error) {
	body, _, err := getResponseBody("/location-area?"+queryString, cfg)
	if err != nil {
		return nil, err
	}
	var data locationsResponse
	if err = json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	_, prevQuery, _ := strings.Cut(data.Previous, "?")
	cfg.PrevLocationQuery = prevQuery
	_, nextQuery, _ := strings.Cut(data.Next, "?")
	cfg.NextLocationQuery = nextQuery
	return data.Results, nil
}

// GetLocationDetails returns a list of Pokemon which can be found at a particular location
func GetLocationDetails(path string, cfg *config.Config) ([]namedApiResource, error) {
	body, _, err := getResponseBody(path, cfg)
	if err != nil {
		return nil, err
	}
	var data locationResponse
	if err = json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	var pokemons []namedApiResource
	for _, encounter := range data.PokemonEncounters {
		pokemons = append(pokemons, encounter.Pokemon)
	}
	return pokemons, nil
}

// CatchPokemon attempts to catch a pokemon based on base_experience and returns true if caught
func CatchPokemon(name string, cfg *config.Config) (bool, error) {
	var pokemon Pokemon
	pokemonBytes, found := cfg.Cache.Get("pokemon/" + name)
	if found {
		log.Println("pokemon found in cache. decoding")
		if err := decodePokemon(pokemonBytes, &pokemon); err != nil {
			return false, err
		}
	} else {
		log.Println("pokemon not found in cache. retrieving via API")
		body, statusCode, err := getResponseBody("/pokemon/"+name, cfg)
		if err != nil {
			if statusCode == http.StatusNotFound {
				err = errors.New("pokemon name not found")
			}
			return false, err
		}
		if err = json.Unmarshal(body, &pokemon); err != nil {
			return false, err
		}
	}
	threshold := cfg.Generator.Float32()
	difficulty := 1 / float32(pokemon.BaseExperience) * 100
	if difficulty > threshold {
		// only store the pokemon if it was caught
		if !found {
			pokemonBytes, err := encodePokemon(&pokemon)
			if err != nil {
				return false, nil
			}
			cfg.Cache.Add("pokemon/"+name, pokemonBytes, true)
		}
		return true, nil
	}
	return false, nil
}

// InspectPokemon retrieves a caught Pokemon from cache and prints its details
func InspectPokemon(name string, cfg *config.Config) (string, error) {
	pokemonBytes, found := cfg.Cache.Get("pokemon/" + name)
	if !found {
		return "", errors.New(name + " has not been caught so cannot be inspected.")
	}
	var pokemon Pokemon
	if err := decodePokemon(pokemonBytes, &pokemon); err != nil {
		return "", err
	}
	var statsStrings []string
	for _, stat := range pokemon.Stats {
		statsStrings = append(statsStrings, fmt.Sprintf("\t%s: %d", stat.Stat.Name, stat.BaseStat))
	}
	statsString := strings.Join(statsStrings, "\n")
	var typesStrings []string
	for _, pokemonType := range pokemon.Types {
		typesStrings = append(typesStrings, fmt.Sprintf("\t%s", pokemonType.Type.Name))
	}
	typesString := strings.Join(typesStrings, "\n")
	return fmt.Sprintf("Name: %s\nHeight: %d\nWeight: %d\nStats:\n%s\nTypes:\n%s\n", pokemon.Name, pokemon.Height, pokemon.Weight, statsString, typesString), nil
}
