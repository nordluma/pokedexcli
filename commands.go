package main

import (
	"encoding/json"
	"errors"
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
		"inspect": {
			name:        "inspect <pokemon>",
			description: "Inspect a pokemon in pokedex",
			callback:    inspectCommand,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Print all pokemons in pokedex",
			callback:    printPokemonsInPokedexCommand,
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
		fmt.Println("you may now inspect in with the inspect command.")
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

func inspectCommand(config *config) error {
	pokemon, ok := config.pokedex[config.item]
	if !ok {
		return errors.New("you have not caught that pokemon")
	}
	printPokemon(pokemon)

	return nil
}

func printPokemon(pokemon pokeapi.PokemonResponse) {
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, s := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", s.Stat.Name, s.BaseStat)
	}

	fmt.Println("Types:")
	for _, kind := range pokemon.Types {
		fmt.Printf("  - %s\n", kind.Type.Name)
	}
}

func printPokemonsInPokedexCommand(config *config) error {
	if len(config.pokedex) == 0 {
		return errors.New("No Pokemons caught")
	}

	fmt.Println("Your Pokedex:")
	for _, pokemon := range config.pokedex {
		fmt.Printf("  - %s\n", pokemon.Name)
	}

	return nil
}

func commandExit(config *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)

	return nil
}
