package main

import (
	"encoding/json"
	"errors"
)

func parseLocationJSON(data []byte) (locationAreaJSON, error) {

	var locationAreas locationAreaJSON
	err := json.Unmarshal(data, &locationAreas)
	if err != nil {
		return locationAreaJSON{}, errors.New("Invalid JSON returned")
	}

	return locationAreas, nil
}

func parsePokemonEncounterJSON(data []byte) (PokemonEncounter, error) {

	var pokeJSON PokemonEncounter
	err := json.Unmarshal(data, &pokeJSON)
	if err != nil {
		return PokemonEncounter{}, errors.New("Invalid JSON returned")
	}

	return pokeJSON, nil
}
