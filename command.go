package main

import (
	"errors"
	"fmt"
	"os"
)

func commandExit(args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return errors.New("WTF, os.Exit failed...")
}

func commandHelp(args ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	for command := range commandHooks {
		fmt.Printf("%s: %s\n%s", commandHooks[command].name, commandHooks[command].description, PLACE_CURSOR)
	}
	return nil
}

func commandCatch(args ...string) error {
	if len(args) < 1 {
		return errors.New("catch: must supply pokemon to catch")
	}

	if currentEncounter.ID == 0 {
		fmt.Println("Please explore an area first!")
	}

	return nil
}

func commandExplore(args ...string) error {
	if len(args) < 1 {
		return errors.New("explore: Need location argument")
	}

	fmt.Printf("Exploring %s...\n", args[0])

	var pokeEncJSON PokemonEncounterJSON
	err := pokeEncJSON.fill(BASE_URL + args[0])
	if err != nil {
		fmt.Println("Explore failed.")
		return fmt.Errorf("Error: %w", err)
	}

	if len(pokeEncJSON.PokemonEncounters) > 0 {
		fmt.Println("Found Pokemon:")
		for _, result := range pokeEncJSON.PokemonEncounters {
			fmt.Printf(" - %s\n", result.Pokemon.Name)
		}
	} else {
		fmt.Println("No Pokemon found/Bad location")
	}

	return nil
}

func commandMapB(args ...string) error {
	err := currentLocationArea.fill(currentLocationArea.Prev)

	if err != nil {
		fmt.Println("Error encountered retrieving locations... Exiting gracefully.")
		commandExit()
	}

	for _, result := range currentLocationArea.Results {
		fmt.Println(result.Name)
	}

	return nil
}

func commandMap(args ...string) error {
	if currentMapUrl == "" {
		currentMapUrl = BASE_URL
	}

	err := currentLocationArea.fill(currentMapUrl)

	if err != nil {
		fmt.Println("Error encountered retrieving locations... Exiting gracefully.")
		commandExit()
	}

	for _, result := range currentLocationArea.Results {
		fmt.Println(result.Name)
	}

	return nil
}
