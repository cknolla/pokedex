package cli

import (
	"github.com/stretchr/testify/assert"
	"pokedex/internal/config"
	"testing"
)

func TestCommandMap(t *testing.T) {
	cfg := config.NewConfig()
	assert.Nil(t, commandMap(&cfg))
}

func TestCommandMapErrorsAtEnd(t *testing.T) {
	cfg := config.NewConfig()
	cfg.NextLocationUrl = ""
	err := commandMap(&cfg)
	assert.EqualError(t, err, "at the end of the locations list")
}

func TestCommandMapb(t *testing.T) {
	cfg := config.NewConfig()
	cfg.PrevLocationUrl = "https://pokeapi.co/api/v2/location-area?offset=0&limit=20"
	assert.Nil(t, commandMapb(&cfg))
}

func TestCommandMapbErrorsAtBeginning(t *testing.T) {
	cfg := config.NewConfig()
	cfg.PrevLocationUrl = ""
	err := commandMapb(&cfg)
	assert.EqualError(t, err, "at the beginning of locations list")
}
