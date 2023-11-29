package cli

import (
	"github.com/stretchr/testify/assert"
	"pokedex/internal/config"
	"testing"
)

func TestCommandMap(t *testing.T) {
	cfg := config.NewConfig()
	assert.Nil(t, commandMap(&cfg, nil))
}

func TestCommandMapErrorsAtEnd(t *testing.T) {
	cfg := config.NewConfig()
	cfg.NextLocationUrl = ""
	err := commandMap(&cfg, nil)
	assert.EqualError(t, err, "at the end of the locations list")
}

func TestCommandMapb(t *testing.T) {
	cfg := config.NewConfig()
	cfg.PrevLocationUrl = "/location-area?offset=0&limit=20"
	assert.Nil(t, commandMapb(&cfg, nil))
}

func TestCommandMapbErrorsAtBeginning(t *testing.T) {
	cfg := config.NewConfig()
	cfg.PrevLocationUrl = ""
	err := commandMapb(&cfg, nil)
	assert.EqualError(t, err, "at the beginning of locations list")
}
