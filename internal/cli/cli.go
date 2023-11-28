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
	Callback    func(*config.Config) error
}

func commandHelp(config *config.Config) error {
	fmt.Printf("\n")
	for name, command := range GetCommands() {
		fmt.Printf("%s: %s\n", name, command.Description)
	}
	return nil
}

func commandExit(config *config.Config) error {
	os.Exit(0)
	return nil
}

func commandMap(config *config.Config) error {
	if config.NextLocationUrl == "" {
		return errors.New("at the end of the locations list")
	}
	var locations api.LocationsData
	results, err := api.GetLocations(config.NextLocationUrl, &locations, config)
	if err != nil {
		return err
	}
	for _, result := range results {
		fmt.Println(result.Name)
	}
	return nil
}

func commandMapb(config *config.Config) error {
	if config.PrevLocationUrl == "" {
		return errors.New("at the beginning of locations list")
	}
	var locations api.LocationsData
	results, err := api.GetLocations(config.PrevLocationUrl, &locations, config)
	if err != nil {
		return err
	}
	for _, result := range results {
		fmt.Println(result.Name)
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
