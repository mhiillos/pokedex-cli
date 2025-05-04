package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mhiillos/pokedex-cli/internal/pokeapi"
	"github.com/mhiillos/pokedex-cli/internal/pokecache"
)

type config struct {
	previous string
	next 		 string
}

type cliCommand struct {
	name				string
	description string
	callback		func(c *config) error
}

func commandExit(c *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func printHelp(commands map[string]cliCommand) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	for _, v := range commands {
		fmt.Printf("%s: %s\n", v.name, v.description)
	}

	return nil
}

func cleanInput(text string) []string {
	cleanedWords := strings.ToLower(text)
	return strings.Fields(cleanedWords)
}

// Prints the next 20 areas
func nextAreas(c *config, cache *pokecache.Cache) error {
	endpoint := "https://pokeapi.co/api/v2/location-area/"
	if c.next != "" {
		endpoint = c.next
	}
	response, err := pokeapi.Get(cache, endpoint)
	if err != nil {
		return err
	}
	for _, area := range response.Results {
		fmt.Println(area.Name)
	}
	c.next = response.Next
	c.previous = response.Previous

	return nil
}

// Prints the previous 20 areas
func prevAreas(c *config, cache *pokecache.Cache) error {
	if c.previous == "" {
		return errors.New("You're on the first page")
	}
	endpoint := c.previous
	response, err := pokeapi.Get(cache, endpoint)
	if err != nil {
		return err
	}
	for _, area := range response.Results {
		fmt.Println(area.Name)
	}
	c.next = response.Next
	c.previous = response.Previous

	return nil
}

func main() {
	s := bufio.NewScanner(os.Stdin)
	cfg := &config{}
	cache, err := pokecache.NewCache(5000 * time.Millisecond)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	commands := map[string]cliCommand{
    "exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
    },
	}
	commands["help"] = cliCommand {
		name: 			 "help",
		description: "Displays a help message",
		callback:		 func(cfg *config) error { return printHelp(commands) },
	}
	commands["map"] = cliCommand {
		name: 			 "map",
		description: "Displays the next 20 map areas",
		callback:    func (cfg *config) error { return nextAreas(cfg, cache) },
	}
	commands["mapb"] = cliCommand {
		name: 			 "mapb",
		description: "Displays the previous 20 map areas",
		callback:    func (cfg *config) error { return prevAreas(cfg, cache) },
	}

	for {
		fmt.Print("Pokedex > ")
		s.Scan()
		input := s.Text()
		cleanedInput := cleanInput(input)
		if len(cleanedInput) == 0 {
			continue
		}
		cmd := cleanedInput[0]

		if cmdStruct, ok := commands[cmd]; ok {
			err := cmdStruct.callback(cfg)
			if err != nil {
					fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}
