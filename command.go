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
	printHelpBanner()
	printHelpDescriptions()
	return nil
}

func printHelpDescriptions() {
	for command := range commandHooks {
		fmt.Printf("%s: %s\n%s", commandHooks[command].name, commandHooks[command].description, PLACE_CURSOR)
	}
}

func printHelpBanner() {
	fmt.Print(HELP_BANNER)
}

func commandPokedex(args ...string) error {
	if err := checkPokedexCommandLengthAndPrintCodex(); err != nil {
		return err
	}
	printPokedexContents()
	return nil
}

func checkPokedexCommandLengthAndPrintCodex() error {
	if err := checkPokedexCommandLength(); err != nil {
		return err
	}
	printPokedexContents()
	return nil
}

func printPokedexContents() {
	fmt.Print(POKEDEX_CONTENTS_BANNER)
	for name := range pokeDex {
		fmt.Printf("  - %s\n", name)
	}
}

func checkPokedexCommandLength() error {
	if len(pokeDex) < 1 {
		fmt.Print(POKEDEX_COMMAND_LENGTH_ERROR_MSG)
		return errors.New("pokedex size error (empty)")
	}

	return nil
}

func checkForPokemonInPokedex(key string, targetPokemon string) bool {
	pokeJSON, ok := pokeDex[key]
	if !ok {
		return false
	}
	return checkPokemonNameAgainstTarget(pokeJSON, targetPokemon)
}

func checkPokemonNameAgainstTarget(pokeJSON PokemonJSON, targetPokemon string) bool {
	return pokeJSON.Name == targetPokemon
}

func checkForPokemonInArea(targetPokemon string) bool {
	found := false
	for key := range pokeDex {
		found = checkForPokemonInPokedex(key, targetPokemon)
	}

	if !found {
		fmt.Printf("You have not caught %s yet!!\n", targetPokemon)
		return found
	}

	return found
}

func commandInspect(args ...string) error {
	if len(args) < 1 {
		return errors.New("inspect: no pokemon given")
	}

	targetPokemon := args[0]

	if !checkForPokemonInArea(targetPokemon) {
		return nil
	}

	pokemon := pokeDex[targetPokemon]
	fmt.Printf(POKEDEX_INSPECT_TEMPLATE, pokemon.Name, pokemon.Height, pokemon.Weight)

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
