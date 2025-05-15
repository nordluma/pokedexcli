package main

import (
	"time"

	"github.com/nordluma/pokedexcli/internal/pokeapi"
)

type config struct {
	client   pokeapi.Client
	next     string
	previous string
	item     string
	pokedex  map[string]pokeapi.PokemonResponse
}

func main() {
	client := pokeapi.NewClient(5*time.Second, 5*time.Minute)
	cfg := config{
		client:  client,
		pokedex: map[string]pokeapi.PokemonResponse{},
	}

	repl(&cfg)
}
