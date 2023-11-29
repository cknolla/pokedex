package cli

import (
	"errors"
	"fmt"
	"os"
	"pokedex/internal/api"
	"pokedex/internal/config"
)

type cliCommand struct {
	Name        string
	Description string
	Callback    func(*config.Config, []string) error
}

func commandHelp(config *config.Config, args []string) error {
	fmt.Printf("\n")
	for name, command := range GetCommands() {
		fmt.Printf("%s: %s\n", name, command.Description)
	}
	return nil
}

func commandExit(config *config.Config, args []string) error {
	os.Exit(0)
	return nil
}

func commandMap(config *config.Config, args []string) error {
	if config.NextLocationUrl == "" {
		return errors.New("at the end of the locations list")
	}
	var locations api.LocationsData
	results, err := api.GetLocations(config.ApiRoot+config.NextLocationUrl, &locations, config)
	if err != nil {
		return err
	}
	for _, result := range results {
		fmt.Println(result.Name)
	}
	return nil
}

func commandMapb(config *config.Config, args []string) error {
	if config.PrevLocationUrl == "" {
		return errors.New("at the beginning of locations list")
	}
	var locations api.LocationsData
	results, err := api.GetLocations(config.ApiRoot+config.PrevLocationUrl, &locations, config)
	if err != nil {
		return err
	}
	for _, result := range results {
		fmt.Println(result.Name)
	}
	return nil
}

func commandExplore(config *config.Config, args []string) error {
	if len(args) < 1 {
		return errors.New("must provide a location name to explore")
	}
	location := args[0]
	fmt.Printf("Exploring %s...\n", location)
	var data api.PokemonResponse
	pokemons, err := api.GetLocationDetails(config.ApiRoot+"/location-area/"+location, &data, config)
	if err != nil {
		return err
	}
	fmt.Println("Found Pokemon:")
	for _, pokemon := range pokemons {
		fmt.Printf(" - %s\n", pokemon.Name)
	}
	return nil
}

func GetCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			Name:        "help",
			Description: "Displays a help message",
			Callback:    commandHelp,
		},
		"exit": {
			Name:        "exit",
			Description: "Exit the Pokedex",
			Callback:    commandExit,
		},
		"explore": {
			Name:        "explore",
			Description: "Explore a particular location for its pokemon",
			Callback:    commandExplore,
		},
		"map": {
			Name:        "map",
			Description: "Get next 20 locations",
			Callback:    commandMap,
		},
		"mapb": {
			Name:        "mapb",
			Description: "Get previous 20 locations",
			Callback:    commandMapb,
		},
	}
}
