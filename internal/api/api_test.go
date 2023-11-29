package api

import (
	"github.com/stretchr/testify/assert"
	"pokedex/internal/config"
	"testing"
)

func TestGetLocations(t *testing.T) {
	testCases := []struct {
		description  string
		url          string
		nextUrl      string
		prevUrl      string
		errorMessage string
	}{
		{
			description:  "get first next",
			url:          "https://pokeapi.co/api/v2/location-area?offset=0&limit=20",
			nextUrl:      "https://pokeapi.co/api/v2/location-area?offset=20&limit=20",
			prevUrl:      "",
			errorMessage: "",
		},
		{
			description:  "error if empty url",
			url:          "",
			nextUrl:      "",
			prevUrl:      "",
			errorMessage: "Get \"\": unsupported protocol scheme \"\"",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			var err error
			cfg := config.NewConfig()
			cfg.NextLocationUrl = testCase.nextUrl
			cfg.PrevLocationUrl = testCase.prevUrl
			var data LocationsData
			locations, err := GetLocations(testCase.url, &data, &cfg)
			if err != nil {
				assert.EqualError(t, err, testCase.errorMessage)
				return
			}
			assert.Equal(t, 20, len(locations), "locations not of expected length")
			assert.Equal(t, testCase.nextUrl, cfg.NextLocationUrl, "resulting nextUrl not as expected")
			assert.Equal(t, testCase.prevUrl, cfg.PrevLocationUrl, "resulting prevUrl not as expected")
		})
	}
}

func TestGetLocationDetails(t *testing.T) {
	cfg := config.NewConfig()
	var data PokemonResponse
	pokemons, err := GetLocationDetails(cfg.ApiRoot+"/location-area/canalave-city-area", &data, &cfg)
	assert.Nil(t, err)
	for index, pokemonName := range []string{
		"tentacool",
		"tentacruel",
		"staryu",
		"magikarp",
		"gyarados",
		"wingull",
		"pelipper",
		"shellos",
		"gastrodon",
		"finneon",
		"lumineon",
	} {
		assert.Equal(t, pokemonName, pokemons[index].Name)
	}
}
