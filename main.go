package main

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
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
	callback		func(c *config, args []string) error
}

func commandExit(c *config, args []string) error {
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
func nextAreas(c *config, client *pokeapi.Client) error {
	endpoint := "https://pokeapi.co/api/v2/location-area/"
	if c.next != "" {
		endpoint = c.next
	}
	response, err := client.GetLocationAreas(endpoint)
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
func prevAreas(c *config, client *pokeapi.Client) error {
	if c.previous == "" {
		return errors.New("You're on the first page")
	}
	endpoint := c.previous
	response, err := client.GetLocationAreas(endpoint)
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

func exploreArea(client *pokeapi.Client, args []string) error {
	endpoint := "https://pokeapi.co/api/v2/location-area/"
	if len(args) != 1 {
		return errors.New("Please provide one location to explore")
	}
	fmt.Printf("Exploring %s\n", args[0])
	endpoint += args[0]
	response, err := client.ExploreLocationArea(endpoint)
	if err != nil {
		return err
	}
	fmt.Println("Found Pokemon:")
	for _, pokemonEncounter := range response.PokemonEncounters {
		fmt.Printf("  - %s\n", pokemonEncounter.Pokemon.Name)
	}
	return nil
}

func catchPokemon(client *pokeapi.Client, pokedex map[string]pokeapi.Pokemon, args []string) error {
	endpoint := "https://pokeapi.co/api/v2/pokemon/"
	if len(args) != 1 {
		return errors.New("Please provide one Pokemon to catch")
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", args[0])
	endpoint += args[0]
	pokemon, err := client.RollPokemon(endpoint)
	if err != nil {
		return err
	}
	if pokemon.Name == ("") {
		fmt.Printf("%s escaped!\n", args[0])
		return nil
	}
	fmt.Printf("%s was caught!\n", args[0])
	pokedex[args[0]] = pokemon
	fmt.Println("You may now inspect it with the inspect command.")
	return nil
}

func inspectPokemon(pokedex map[string]pokeapi.Pokemon, args []string) error {
	if len(args) != 1 {
		return errors.New("Please provide one Pokemon to inspect")
	}
	pokemon, ok := pokedex[args[0]]
	if !ok {
		return fmt.Errorf("Pokemon %q not in Pokedex", args[0])
	}
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("-%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, pokemonType := range pokemon.Types {
		fmt.Printf("- %s\n", pokemonType.Type.Name)

	}
	return nil
}

func listPokemon(pokedex map[string]pokeapi.Pokemon) error {
	fmt.Println("Your Pokedex:")
	for name := range(pokedex) {
		fmt.Printf("- %s\n", name)
	}
	return nil
}

func main() {
	s := bufio.NewScanner(os.Stdin)
	cfg := &config{}
	pokedex := make(map[string]pokeapi.Pokemon)
	cache, err := pokecache.NewCache(5000 * time.Millisecond)
	client := &pokeapi.Client{HTTP: &http.Client{}, Cache: cache}
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
		description: "Lists the available commands",
		callback:		 func(cfg *config, args []string) error { return printHelp(commands) },
	}
	commands["map"] = cliCommand {
		name: 			 "map",
		description: "Displays the next 20 map areas",
		callback:    func (cfg *config, args[]string) error { return nextAreas(cfg, client) },
	}
	commands["mapb"] = cliCommand {
		name: 			 "mapb",
		description: "Displays the previous 20 map areas",
		callback:    func (cfg *config, args []string) error { return prevAreas(cfg, client) },
	}
	commands["explore"] = cliCommand {
		name: "explore <location_area>",
		description: "Displays pokemon in the location area",
		callback: func (cfg *config, args []string) error { return exploreArea(client, args) },
	}
	commands["catch"] = cliCommand {
		name: "catch <pokemon_name>",
		description: "Attempt to catch the specified Pokemon",
		callback: func (cfg *config, args []string) error { return catchPokemon(client, pokedex, args) },
	}
	commands["inspect"] = cliCommand {
		name: "inspect <pokemon_name>",
		description: "View a caught Pokemon's stats",
		callback: func (cfg *config, args []string) error { return inspectPokemon(pokedex, args) },
	}
	commands["pokedex"] = cliCommand {
		name: "pokedex",
		description: "View all Pokemon in your Pokedex",
		callback: func (cfg *config, args []string) error { return listPokemon(pokedex) },
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
			err := cmdStruct.callback(cfg, cleanedInput[1:])
			if err != nil {
					fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}
