package config

import (
	"math/rand"
	"pokedex/internal/pokecache"
	"time"
)

type Config struct {
	Cache            pokecache.Cache
	Generator        *rand.Rand
	ApiRoot          string
	NextLocationPath string
	PrevLocationPath string
}

func NewConfig() Config {
	generator := rand.New(rand.NewSource(time.Now().UnixNano()))
	return Config{
		Cache:            pokecache.NewCache(5 * time.Minute),
		Generator:        generator,
		ApiRoot:          "https://pokeapi.co/api/v2",
		NextLocationPath: "/location-area?offset=0&limit=20",
		PrevLocationPath: "",
	}
}
