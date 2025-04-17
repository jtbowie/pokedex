package main

import (
	"encoding/json"
	"errors"
)

func (lAJ locationAreaJSON) parseJSON(data []byte) (locationAreaJSON, error) {

	var locationAreas locationAreaJSON
	err := json.Unmarshal(data, &locationAreas)
	if err != nil {
		return locationAreaJSON{}, errors.New("Invalid JSON returned")
	}

	return locationAreas, nil
}

func (pEJ PokemonEncounterJSON) parseJSON(data []byte) (PokemonEncounterJSON, error) {

	var pokeJSON PokemonEncounterJSON
	err := json.Unmarshal(data, &pokeJSON)
	if err != nil {
		return PokemonEncounterJSON{}, errors.New("Invalid JSON returned")
	}

	return pokeJSON, nil
}
