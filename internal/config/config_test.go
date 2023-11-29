package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewConfig(t *testing.T) {
	config := NewConfig()
	expectedNextLocationQuery := "offset=0&limit=20"
	expectedPrevLocationQuery := ""
	assert.Equal(t, expectedNextLocationQuery, config.NextLocationQuery)
	assert.Equal(t, expectedPrevLocationQuery, config.PrevLocationQuery)
}
