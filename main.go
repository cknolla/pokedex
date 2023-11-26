package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"pokedex/internal/api"
	"pokedex/internal/config"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config.Config) error
}

func commandHelp(config *config.Config) error {
	fmt.Printf("\n")
	for name, command := range getCommands() {
		fmt.Printf("%s: %s\n", name, command.description)
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
	return api.GetLocations(config.NextLocationUrl, &locations, config)
}

func commandMapb(config *config.Config) error {
	if config.PrevLocationUrl == "" {
		return errors.New("at the beginning of locations list")
	}
	var locations api.LocationsData
	return api.GetLocations(config.PrevLocationUrl, &locations, config)
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Get next 20 locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Get previous 20 locations",
			callback:    commandMapb,
		},
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	commands := getCommands()
	conf := config.Config{
		NextLocationUrl: "https://pokeapi.co/api/v2/location-area?offset=0&limit=20",
		PrevLocationUrl: "",
	}
	fmt.Printf("pokedex > ")
	for scanner.Scan() {
		line := scanner.Text()
		//fmt.Printf("%s\n", line)
		for _, command := range commands {
			if line == command.name {
				err := command.callback(&conf)
				if err != nil {
					fmt.Printf("%s\n", err)
				}
			}
		}
		fmt.Printf("pokedex > ")
	}
}
