package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewConfig(t *testing.T) {
	config := NewConfig()
	expectedNextLocationUrl := "/location-area?offset=0&limit=20"
	expectedPrevLocationUrl := ""
	assert.Equal(t, expectedNextLocationUrl, config.NextLocationPath)
	assert.Equal(t, expectedPrevLocationUrl, config.PrevLocationPath)
}
