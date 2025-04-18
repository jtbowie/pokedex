package main

import (
	"errors"
	"fmt"
	"math/rand"
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

func commandInspect(args ...string) error {
	if len(args) < 1 {
		return errors.New("inspect: no pokemon given")
	}

	targetPokemon := args[0]
	found := false

	for key := range pokeDex {
		if pokeJSON, ok := pokeDex[key]; ok {
			if pokeJSON.Name == targetPokemon {
				found = true
			}
		}
	}

	if !found {
		fmt.Printf("You have not caught %s yet!!\n", args[0])
		return nil
	}

	pokemon := pokeDex[targetPokemon]

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Stats:")

	for idx := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n",
			pokemon.Stats[idx].Stat.Name,
			pokemon.Stats[idx].BaseStat)
	}

	fmt.Println("Types:")
	for idx := range pokemon.Types {
		fmt.Printf("  - %s\n",
			pokemon.Types[idx].Type.Name)
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

	targetPokemon := args[0]
	fmt.Printf("Throwing a Pokeball at %s...\n", targetPokemon)
	targetPokemonURL, err := pokemonInArea(targetPokemon)
	if err != nil {
		fmt.Printf("I'm sorry, %s isn't found in %s\n", targetPokemon, currentEncounter.Name)
		return errors.New("pokemon not found")
	}

	var pokemonJSON PokemonJSON
	err = pokemonJSON.fill(targetPokemonURL)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	xp := pokemonJSON.BaseExperience

	toughness := 100 - (xp / 4)

	fmt.Printf("Toughness: %d\n", toughness)

	if rand.Intn(165) <= toughness {
		fmt.Printf("You caught %s!!!\n", targetPokemon)
		fmt.Printf("Inspect this pokemon to see its attributes! 'inspect %s'\n", targetPokemon)
		pokeDex[targetPokemon] = pokemonJSON
		return nil
	}

	fmt.Printf("%s 'scaped :( :(\n", targetPokemon)

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
		return fmt.Errorf("error: %w", err)
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
