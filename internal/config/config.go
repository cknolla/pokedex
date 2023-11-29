package config

import (
	"pokedex/internal/pokecache"
	"time"
)

type Config struct {
	Cache           pokecache.Cache
	ApiRoot         string
	NextLocationUrl string
	PrevLocationUrl string
}

func NewConfig() Config {
	return Config{
		Cache:           pokecache.NewCache(5 * time.Minute),
		ApiRoot:         "https://pokeapi.co/api/v2",
		NextLocationUrl: "/location-area?offset=0&limit=20",
		PrevLocationUrl: "",
	}
}
