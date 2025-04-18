package main

import "errors"

func pokemonInArea(targetPokemon string) (string, error) {

	for idx := range currentEncounter.PokemonEncounters {
		if currentEncounter.PokemonEncounters[idx].Pokemon.Name == targetPokemon {
			return currentEncounter.PokemonEncounters[idx].Pokemon.URL, nil
		}

	}
	return "", errors.New("your princess is in another castle...")
}
