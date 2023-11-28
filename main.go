package main

import (
	"bufio"
	"fmt"
	"os"
	"pokedex/internal/cli"
	"pokedex/internal/config"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	commands := cli.GetCommands()
	conf := config.NewConfig()
	fmt.Printf("pokedex > ")
	for scanner.Scan() {
		line := scanner.Text()
		//fmt.Printf("%s\n", line)
		for _, command := range commands {
			if line == command.Name {
				err := command.Callback(&conf)
				if err != nil {
					fmt.Printf("%s\n", err)
				}
			}
		}
		fmt.Printf("pokedex > ")
	}
}
