package cli

import (
	"errors"
	"fmt"
	"os"
	"pokedex/internal/api"
	"pokedex/internal/config"
	"strings"
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
	if config.NextLocationQuery == "" {
		return errors.New("at the end of the locations list")
	}
	results, err := api.GetLocations(config.NextLocationQuery, config)
	if err != nil {
		return err
	}
	for _, result := range results {
		fmt.Println(result.Name)
	}
	return nil
}

func commandMapb(config *config.Config, args []string) error {
	if config.PrevLocationQuery == "" {
		return errors.New("at the beginning of locations list")
	}
	results, err := api.GetLocations(config.PrevLocationQuery, config)
	if err != nil {
		return err
	}
	for _, result := range results {
		fmt.Println(result.Name)
	}
	return nil
}

func commandExplore(config *config.Config, args []string) error {
	if len(args) != 1 {
		return errors.New("must provide a location name to explore")
	}
	location := args[0]
	fmt.Printf("Exploring %s...\n", location)
	pokemons, err := api.GetLocationDetails("/location-area/"+location, config)
	if err != nil {
		return err
	}
	fmt.Println("Found Pokemon:")
	for _, pokemon := range pokemons {
		fmt.Printf(" - %s\n", pokemon.Name)
	}
	return nil
}

func commandCatch(config *config.Config, args []string) error {
	if len(args) != 1 {
		return errors.New("must provide a pokemon name to catch")
	}
	pokemonName := args[0]
	fmt.Printf("Throwing a pokeball at %s...\n", pokemonName)
	caught, err := api.CatchPokemon(pokemonName, config)
	if err != nil {
		return err
	}
	if caught {
		fmt.Printf("%s was caught!\n", pokemonName)
	} else {
		fmt.Printf("%s escaped!\n", pokemonName)
	}
	return nil
}

func commandInspect(config *config.Config, args []string) error {
	if len(args) != 1 {
		return errors.New("must provide a pokemon name to inspect")
	}
	pokemonName := args[0]
	outputString, err := api.InspectPokemon(pokemonName, config)
	if err != nil {
		return err
	}
	fmt.Printf("%s", outputString)
	return nil
}

func commandPokedex(cfg *config.Config, args []string) error {
	fmt.Println("Your Pokedex:")
	for _, key := range cfg.Cache.GetKeys() {
		name, found := strings.CutPrefix(key, "pokemon/")
		if found {
			fmt.Printf("\t%s\n", name)
		}
	}
	return nil
}

func GetCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"catch": {
			Name:        "catch",
			Description: "Attempt to catch a Pokemon by name",
			Callback:    commandCatch,
		},
		"exit": {
			Name:        "exit",
			Description: "Exit the Pokedex",
			Callback:    commandExit,
		},
		"explore": {
			Name:        "explore",
			Description: "Explore a particular location for its Pokemon",
			Callback:    commandExplore,
		},
		"help": {
			Name:        "help",
			Description: "Displays a help message",
			Callback:    commandHelp,
		},
		"inspect": {
			Name:        "inspect",
			Description: "Displays information about Pokemon that have been caught",
			Callback:    commandInspect,
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
		"pokedex": {
			Name:        "pokedex",
			Description: "Print all caught Pokemon",
			Callback:    commandPokedex,
		},
	}
}
