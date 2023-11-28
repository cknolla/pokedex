package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewConfig(t *testing.T) {
	config := NewConfig()
	expectedNextLocationUrl := "https://pokeapi.co/api/v2/location-area?offset=0&limit=20"
	expectedPrevLocationUrl := ""
	assert.Equal(t, expectedNextLocationUrl, config.NextLocationUrl)
	assert.Equal(t, expectedPrevLocationUrl, config.PrevLocationUrl)
}
