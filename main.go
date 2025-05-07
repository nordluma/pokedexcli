package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/nordluma/pokedexcli/internal"
)

func cleanInput(text string) []string {
	lower := strings.ToLower(text)
	trimmed := strings.TrimSpace(lower)
	words := strings.Fields(trimmed)

	return words
}

type config struct {
	cache    *internal.Cache
	next     string
	previous string
	item     string
}

var commands map[string]clicommand

type clicommand struct {
	name        string
	description string
	callback    func(*config) error
}

func commandExit(config *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)

	return nil
}

func printUsage(config *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("Usage:\n\n")

	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}

	return nil
}

type LocationAreaResponse struct {
	Count    int
	Next     string
	Previous string
	Results  []Area
}

type Area struct {
	Name string
	Url  string
}

func mapCommand(config *config) error {
	url := "https://pokeapi.co/api/v2/location-area/"
	if config.next != "" {
		url = config.next
	}

	var areaResponse LocationAreaResponse
	data, err := get(url, config.cache)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(data, &areaResponse); err != nil {
		return err
	}

	printAreas(areaResponse.Results)

	config.next = areaResponse.Next
	config.previous = areaResponse.Previous

	return nil
}

func mapbCommand(config *config) error {
	if config.previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	var areaResponse LocationAreaResponse
	data, err := get(config.previous, config.cache)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(data, &areaResponse); err != nil {
		return err
	}

	printAreas(areaResponse.Results)

	config.next = areaResponse.Next
	config.previous = areaResponse.Previous

	return nil
}

func printAreas(areas []Area) {
	for _, area := range areas {
		fmt.Println(area.Name)
	}
}

type ExploreResponse struct {
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

type NamedAPIResource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type PokemonEncounter struct {
	Pokemon NamedAPIResource `json:"pokemon"`
}

func exploreCommand(config *config) error {
	if config.item == "" {
		return fmt.Errorf("no location given")
	}

	url := "https://pokeapi.co/api/v2/location-area/" + config.item

	var exploreRes ExploreResponse
	data, err := get(url, config.cache)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(data, &exploreRes); err != nil {
		return err
	}

	printPokemonsInArea(exploreRes.PokemonEncounters)

	return nil
}

func printPokemonsInArea(encounters []PokemonEncounter) {
	fmt.Println("Found Pokemon:")
	for _, pokemon := range encounters {
		fmt.Printf(" - %s\n", pokemon.Pokemon.Name)
	}
}

func get(url string, cache *internal.Cache) ([]byte, error) {
	var areaResponse LocationAreaResponse

	// check cache
	entry, found := cache.Get(url)
	if found {
		if err := json.Unmarshal(entry, &areaResponse); err != nil {
			return nil, err
		}

		return entry, nil
	}

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non ok response: %d", res.StatusCode)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	cache.Add(url, data)

	return data, nil
}

func init() {
	commands = map[string]clicommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    printUsage,
		},
		"map": {
			name:        "map",
			description: "Get the next 20 areas",
			callback:    mapCommand,
		},
		"mapb": {
			name:        "mapb",
			description: "Get the previous 20 areas",
			callback:    mapbCommand,
		},
		"explore": {
			name:        "explore <location_name>",
			description: "Explore a location",
			callback:    exploreCommand,
		},
	}
}

func repl(config *config) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("pokedex > ")
		scanner.Scan()
		input := scanner.Text()

		args := cleanInput(input)
		cmdInput := args[0]
		if len(args) > 1 {
			config.item = args[1]
		}

		cmd, ok := commands[cmdInput]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		if err := cmd.callback(config); err != nil {
			fmt.Println(err)
		}
	}
}

func main() {
	config := &config{cache: internal.NewCache(2 * time.Minute)}

	repl(config)
}
