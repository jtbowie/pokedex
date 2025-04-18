package main

import (
	"encoding/json"
	"errors"
)

func (lAJ *locationAreaJSON) parseJSON(data []byte) error {

	err := json.Unmarshal(data, lAJ)
	if err != nil {
		return errors.New("invalid JSON returned")
	}

	return nil
}

func (pEJ *PokemonEncounterJSON) parseJSON(data []byte) error {

	err := json.Unmarshal(data, pEJ)
	if err != nil {
		return errors.New("invalid JSON returned")
	}

	return nil
}

func (pJ *PokemonJSON) parseJSON(data []byte) error {

	err := json.Unmarshal(data, pJ)
	if err != nil {
		return errors.New("invalid JSON returned")
	}

	return nil
}
