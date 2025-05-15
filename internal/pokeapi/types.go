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

type PokemonResponse struct {
	Id                     int                `json:"id"`
	Name                   string             `json:"name"`
	BaseExperience         int                `json:"base_experience"`
	Height                 int                `json:"height"`
	IsDefault              bool               `json:"is_default"`
	Order                  int                `json:"order"`
	Weight                 int                `json:"weight"`
	Ablitites              []PokemonAbility   `json:"abilities"`
	Forms                  []NamedAPIResource `json:"forms"`
	GameIndices            []GameIndex        `json:"game_indices"`
	HeldItems              []string           `json:"help_items"`
	LocationAreaEncounters string             `json:"location_area_encounters"`
	Moves                  []NamedAPIResource `json:"moves"`
	Species                NamedAPIResource   `json:"species"`
	Cries                  struct {
		Latest string `json:"latest"`
		Legacy string `json:"legacy"`
	} `json:"cries"`
	Stats         []PokemonStats       `json:"stats"`
	Types         []PokemonType        `json:"types"`
	PastTypes     []PokemonPastType    `json:"past_types"`
	PastAbilities []PokemonPastAbility `json:"past_abilities"`
}

func (p PokemonResponse) FindFromStats(stat string) (int, bool) {
	for _, s := range p.Stats {
		if s.Stat.Name == stat {
			return s.BaseStat, true
		}
	}

	return 0, false
}

type PokemonAbility struct {
	IsHidden bool             `json:"is_hidden"`
	Slot     int              `json:"slot"`
	Ability  NamedAPIResource `json:"ability"`
}

type GameIndex struct {
	GameIndex int              `json:"game_index"`
	Item      NamedAPIResource `json:"item"`
}

type HeldItems struct {
	Item           NamedAPIResource `json:"item"`
	VersionDetails VersionDetails   `json:"version_details"`
}

type VersionDetails struct {
	Rarity  int              `json:"rarity"`
	Version NamedAPIResource `json:"version"`
}

type PokemonMoves struct {
	Move                NamedAPIResource      `json:"move"`
	VersionGroupDetails []VersionGroupDetails `json:"version_group_details"`
}

type VersionGroupDetails struct {
	LevelLearnedAt  int              `json:"level_learned_at"`
	VersionGroup    NamedAPIResource `json:"version_group"`
	MoveLearnMethod NamedAPIResource `json:"move_learn_method"`
	Order           int              `json:"order"`
}

type PokemonStats struct {
	BaseStat int              `json:"base_stat"`
	Effort   int              `json:"effort"`
	Stat     NamedAPIResource `json:"stat"`
}

type PokemonType struct {
	Slot int              `json:"slot"`
	Type NamedAPIResource `json:"type"`
}

type PokemonPastType struct {
	Generation NamedAPIResource `json:"generation"`
	Types      []struct {
		Slot int              `json:"slot"`
		Type NamedAPIResource `json:"type"`
	} `json:"types"`
}

type PokemonPastAbility struct {
	Generation NamedAPIResource `json:"generation"`
	Abilities  []PokemonAbility `json:"abilities"`
}
