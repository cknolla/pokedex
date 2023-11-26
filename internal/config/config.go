package config

type Config struct {
	NextLocationUrl string
	PrevLocationUrl string
}

func NewConfig() Config {
	return Config{
		NextLocationUrl: "https://pokeapi.co/api/v2/location-area?offset=0&limit=20",
		PrevLocationUrl: "",
	}
}
