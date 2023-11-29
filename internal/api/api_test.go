package api

import (
	"github.com/stretchr/testify/assert"
	"pokedex/internal/config"
	"testing"
)

func TestGetLocations(t *testing.T) {
	testCases := []struct {
		description  string
		queryString  string
		nextQuery    string
		prevQuery    string
		errorMessage string
	}{
		{
			description:  "get first next",
			queryString:  "offset=0&limit=20",
			nextQuery:    "offset=20&limit=20",
			prevQuery:    "",
			errorMessage: "",
		},
		{
			description:  "double-sided",
			queryString:  "offset=20&limit=20",
			nextQuery:    "offset=40&limit=20",
			prevQuery:    "offset=0&limit=20",
			errorMessage: "",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			var err error
			cfg := config.NewConfig()
			cfg.NextLocationQuery = testCase.nextQuery
			cfg.PrevLocationQuery = testCase.prevQuery
			locations, err := GetLocations(testCase.queryString, &cfg)
			if err != nil {
				assert.EqualError(t, err, testCase.errorMessage)
				return
			}
			assert.Equal(t, 20, len(locations))
			assert.Equal(t, testCase.nextQuery, cfg.NextLocationQuery)
			assert.Equal(t, testCase.prevQuery, cfg.PrevLocationQuery)
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

func TestInspectPokemon(t *testing.T) {
	cfg := config.NewConfig()
	cfg.Generator.Seed(1)
	caught, err := CatchPokemon("pikachu", &cfg)
	assert.Nil(t, err)
	assert.True(t, caught)
	outputString, err := InspectPokemon("pikachu", &cfg)
	assert.Nil(t, err)
	assert.Equal(t, "Name: pikachu\nHeight: 4\nWeight: 60\nStats:\n\thp: 35\n\tattack: 55\n\tdefense: 40\n\tspecial-attack: 50\n\tspecial-defense: 50\n\tspeed: 90\nTypes:\n\telectric\n", outputString)
}
