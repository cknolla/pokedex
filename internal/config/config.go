package config

import (
	"math/rand"
	"pokedex/internal/pokecache"
	"time"
)

type Config struct {
	Cache             pokecache.Cache
	Generator         *rand.Rand
	ApiRoot           string
	NextLocationQuery string
	PrevLocationQuery string
}

func NewConfig() Config {
	return Config{
		Cache:             pokecache.NewCache(5 * time.Minute),
		Generator:         rand.New(rand.NewSource(time.Now().UnixNano())),
		ApiRoot:           "https://pokeapi.co/api/v2",
		NextLocationQuery: "offset=0&limit=20",
		PrevLocationQuery: "",
	}
}
