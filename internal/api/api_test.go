package api

import (
	"github.com/stretchr/testify/assert"
	"pokedex/internal/config"
	"testing"
)

func TestGetLocations(t *testing.T) {
	testCases := []struct {
		description  string
		path         string
		nextPath     string
		prevPath     string
		errorMessage string
	}{
		{
			description:  "get first next",
			path:         "/location-area?offset=0&limit=20",
			nextPath:     "/location-area?offset=20&limit=20",
			prevPath:     "",
			errorMessage: "",
		},
		{
			description:  "error if empty path",
			path:         "",
			nextPath:     "",
			prevPath:     "",
			errorMessage: "path must not be empty",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			var err error
			cfg := config.NewConfig()
			cfg.NextLocationUrl = testCase.nextPath
			cfg.PrevLocationUrl = testCase.prevPath
			locations, err := GetLocations(testCase.path, &cfg)
			if err != nil {
				assert.EqualError(t, err, testCase.errorMessage)
				return
			}
			assert.Equal(t, 20, len(locations), "locations not of expected length")
			assert.Equal(t, testCase.nextPath, cfg.NextLocationUrl, "resulting nextPath not as expected")
			assert.Equal(t, testCase.prevPath, cfg.PrevLocationUrl, "resulting prevPath not as expected")
		})
	}
}

func TestGetLocationDetails(t *testing.T) {
	cfg := config.NewConfig()
	pokemons, err := GetLocationDetails("/location-area/canalave-city-area", &cfg)
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

func TestCatchPokemon(t *testing.T) {
	testCases := []struct {
		description    string
		name           string
		expectedCaught bool
		expectedError  string
	}{
		{
			description:    "catch pikachu",
			name:           "pikachu",
			expectedCaught: true,
			expectedError:  "",
		},
		{
			description:    "fail to catch zapdos",
			name:           "zapdos",
			expectedCaught: false,
			expectedError:  "",
		},
		{
			description:    "bad name",
			name:           "NaN",
			expectedCaught: false,
			expectedError:  "pokemon name not found",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			cfg := config.NewConfig()
			cfg.Generator.Seed(1)
			caught, err := CatchPokemon(testCase.name, &cfg)
			if testCase.expectedError != "" {
				assert.EqualError(t, err, testCase.expectedError)
				return
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, testCase.expectedCaught, caught)
		})
	}
}
