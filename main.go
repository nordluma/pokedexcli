package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
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

type AreaResponse struct {
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

	areaResponse, err := getMaps(url, config.cache)
	if err != nil {
		return err
	}

	config.next = areaResponse.Next
	config.previous = areaResponse.Previous

	return nil
}

func mapbCommand(config *config) error {
	if config.previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	areaResponse, err := getMaps(config.previous, config.cache)
	if err != nil {
		return err
	}

	config.next = areaResponse.Next
	config.previous = areaResponse.Previous

	return nil
}

func getMaps(url string, cache *internal.Cache) (AreaResponse, error) {
	var areaResponse AreaResponse

	// check cache
	entry, found := cache.Get(url)
	if found {
		if err := json.Unmarshal(entry, &areaResponse); err != nil {
			return areaResponse, err
		}

		return areaResponse, nil
	}

	res, err := http.Get(url)
	if err != nil {
		return areaResponse, err
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return areaResponse, errors.New("request for areas failed")
	}

	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(&areaResponse); err != nil {
		return areaResponse, err
	}

	for _, area := range areaResponse.Results {
		fmt.Println(area.Name)
	}

	return areaResponse, nil
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
			description: "Get the next 20 maps",
			callback:    mapCommand,
		},
		"mapb": {
			name:        "mapb",
			description: "Get the previous 20 maps",
			callback:    mapbCommand,
		},
	}
}

func main() {
	config := &config{cache: internal.NewCache(2 * time.Minute)}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		args := cleanInput(input)
		cmdInput := args[0]
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
