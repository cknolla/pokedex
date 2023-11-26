package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
	nextLocationUrl string
	prevLocationUrl string
}

type locationResult struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type locationsData struct {
	Count    int              `json:"count"`
	Next     string           `json:"next"`
	Previous string           `json:"previous"`
	Results  []locationResult `json:"results"`
}

func commandHelp(config *config) error {
	fmt.Printf("\n")
	for name, command := range getCommands() {
		fmt.Printf("%s: %s\n", name, command.description)
	}
	return nil
}

func commandExit(config *config) error {
	os.Exit(0)
	return nil
}

func commandMap(config *config) error {
	response, err := http.Get(config.nextLocationUrl)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(response.Body)
	response.Body.Close()
	if response.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and \nbody: %s\n", response.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}
	var locations locationsData
	err = json.Unmarshal(body, &locations)
	if err != nil {
		log.Fatal(err)
	}
	config.prevLocationUrl = locations.Previous
	config.nextLocationUrl = locations.Next
	for _, result := range locations.Results {
		fmt.Println(result.Name)
	}
	return nil
}

func commandMapb(config *config) error {
	if config.prevLocationUrl == "" {
		return errors.New("at the beginning of locations list")
	}
	response, err := http.Get(config.prevLocationUrl)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(response.Body)
	response.Body.Close()
	if response.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and \nbody: %s\n", response.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}
	var locations locationsData
	err = json.Unmarshal(body, &locations)
	if err != nil {
		log.Fatal(err)
	}
	config.prevLocationUrl = locations.Previous
	config.nextLocationUrl = locations.Next
	for _, result := range locations.Results {
		fmt.Println(result.Name)
	}
	return nil
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
	conf := config{
		nextLocationUrl: "https://pokeapi.co/api/v2/location-area?offset=0&limit=20",
		prevLocationUrl: "",
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
