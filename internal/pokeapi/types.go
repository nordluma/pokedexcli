package pokeapi

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
