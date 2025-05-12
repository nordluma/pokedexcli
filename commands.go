package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"

	"github.com/nordluma/pokedexcli/internal/pokeapi"
)

type clicommand struct {
	name        string
	description string
	callback    func(*config) error
}

func getCommands() map[string]clicommand {
	return map[string]clicommand{
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
		"catch": {
			name:        "catch <pokemon>",
			description: "Catch a pokemon",
			callback:    catchCommand,
		},
	}

}

func printUsage(config *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("Usage:\n\n")

	for _, cmd := range getCommands() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}

	return nil
}

func mapCommand(config *config) error {
	url := "https://pokeapi.co/api/v2/location-area/"
	if config.next != "" {
		url = config.next
	}

	var areaResponse pokeapi.LocationAreaResponse
	data, err := config.client.Get(url)
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

	var areaResponse pokeapi.LocationAreaResponse
	data, err := config.client.Get(config.previous)
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

func printAreas(areas []pokeapi.Area) {
	for _, area := range areas {
		fmt.Println(area.Name)
	}
}

func exploreCommand(config *config) error {
	if config.item == "" {
		return fmt.Errorf("no location given")
	}

	url := "https://pokeapi.co/api/v2/location-area/" + config.item

	var exploreRes pokeapi.ExploreResponse
	data, err := config.client.Get(url)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(data, &exploreRes); err != nil {
		return err
	}

	printPokemonsInArea(exploreRes.PokemonEncounters)

	return nil
}

func printPokemonsInArea(encounters []pokeapi.PokemonEncounter) {
	fmt.Println("Found Pokemon:")
	for _, pokemon := range encounters {
		fmt.Printf(" - %s\n", pokemon.Pokemon.Name)
	}
}

func catchCommand(config *config) error {
	fmt.Printf("Throwing a Pokeball at %s...\n", config.item)

	url := "https://pokeapi.co/api/v2/pokemon/" + config.item

	var pokemonResponse pokeapi.PokemonResponse
	data, err := config.client.Get(url)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(data, &pokemonResponse); err != nil {
		return err
	}

	pokemonCaught := tryCatching(pokemonResponse.BaseExperience)
	if pokemonCaught {
		fmt.Printf("%s was caught!\n", config.item)
		config.pokedex[config.item] = pokemonResponse
	} else {
		fmt.Printf("%s escaped!\n", config.item)
	}

	return nil
}

func tryCatching(baseStat int) bool {
	roll := rand.Intn(100) + 1
	threshold := 100 - baseStat

	return roll <= threshold
}

func commandExit(config *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)

	return nil
}
