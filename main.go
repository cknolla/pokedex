package main

import (
	"bufio"
	"fmt"
	"os"
	"pokedex/internal/cli"
	"pokedex/internal/config"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	commands := cli.GetCommands()
	conf := config.NewConfig()
	fmt.Printf("pokedex > ")
	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Split(line, " ")
		cmd, args := words[0], words[1:]
		//fmt.Printf("%s\n", line)
		for _, command := range commands {
			if cmd == command.Name {
				err := command.Callback(&conf, args)
				if err != nil {
					fmt.Printf("%s\n", err)
				}
			}
		}
		fmt.Printf("pokedex > ")
	}
}
