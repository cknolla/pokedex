package config

import (
	"pokedex/internal/test"
	"testing"
)

func TestNewConfig(t *testing.T) {
	t.Run("default values", func(t *testing.T) {
		config := NewConfig()
		expectedNextLocationUrl := "https://pokeapi.co/api/v2/location-area?offset=0&limit=10"
		expectedPrevLocationUrl := ""
		if config.NextLocationUrl != expectedNextLocationUrl {
			t.Error(test.WriteStringDiff(expectedNextLocationUrl, config.NextLocationUrl))
		}
		if config.PrevLocationUrl != expectedPrevLocationUrl {
			t.Error(test.WriteStringDiff(expectedPrevLocationUrl, config.PrevLocationUrl))
		}
	})
}
