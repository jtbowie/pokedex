package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	pc "github.com/jtbowie/pokedex/internal/pokecache"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type cliCommand struct {
	name        string
	description string
	callBack    func(args ...string) error
}

type locationArea struct {
	Count   int    `json:"count"`
	Next    string `json:"next"`
	Prev    string `json:"previous"`
	Results []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

var currentMapUrl = ""
var currentLocationArea locationArea
var pokeCache *pc.Cache

const PLACE_CURSOR string = "\033[H\033[3J\033[80;1H"
const BASE_URL string = "https://pokeapi.co/api/v2/location-area/"

var commandHooks map[string]cliCommand = make(map[string]cliCommand)

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

func commandExplore(args ...string) error {
	if len(args) < 1 {
		return errors.New("explore: Need location argument")
	}

	fmt.Printf("Exploring %s...\n", args[0])

	pokeJSON, err := fillPokemon(BASE_URL + args[0])
	if err != nil {
		return fmt.Errorf("Error: %w\n", err)
	}

	if len(pokeJSON.PokemonEncounters) > 0 {
		fmt.Println("Found Pokemon:")
		for _, result := range pokeJSON.PokemonEncounters {
			fmt.Printf(" - %s\n", result.Pokemon.Name)
		}
	} else {
		fmt.Println("No Pokemon found/Bad location")
	}

	return nil
}

func commandMapB(args ...string) error {
	if currentMapUrl == "null" {
		fmt.Println("Final page reached")
		return nil
	}

	currentLocationArea, err := fillLocationArea(currentLocationArea.Prev)

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
	if currentMapUrl == "null" {
		fmt.Println("Final page reached")
		return nil
	}

	currentLocationArea, err := fillLocationArea(currentMapUrl)

	if err != nil {
		fmt.Println("Error encountered retrieving locations... Exiting gracefully.")
		commandExit()
	}

	for _, result := range currentLocationArea.Results {
		fmt.Println(result.Name)
	}

	return nil
}

func parseLocationJSON(data []byte) (locationArea, error) {

	var locationAreas locationArea
	err := json.Unmarshal(data, &locationAreas)
	if err != nil {
		return locationArea{}, errors.New("Invalid JSON returned")
	}

	return locationAreas, nil
}

func parsePokemonJSON(data []byte) (Pokemon, error) {

	var pokeJSON Pokemon
	err := json.Unmarshal(data, &pokeJSON)
	if err != nil {
		return Pokemon{}, errors.New("Invalid JSON returned")
	}

	return pokeJSON, nil
}

func fillPokemon(url string) (Pokemon, error) {
	var data []byte

	if url == "" {
		return Pokemon{}, errors.New("Enter a url dude.")
	}
	if cacheItem, ok := pokeCache.Get(url); ok {
		data = cacheItem
	} else {
		res, err := http.Get(url)
		if err != nil {
			return Pokemon{}, errors.New("Error connecting to endpoint.")
		}
		defer res.Body.Close()
		data, err = io.ReadAll(res.Body)
		if err != nil {
			return Pokemon{}, err
		}
		pokeCache.Add(url, data)
	}
	pokeJSON, err := parsePokemonJSON(data)
	if err != nil {
		return Pokemon{}, errors.New("JSON Parsing failed")
	}

	return pokeJSON, nil
}

func fillLocationArea(url string) (locationArea, error) {
	var data []byte

	if url == "" {
		url = BASE_URL
	}

	if cacheItem, ok := pokeCache.Get(url); ok {
		data = cacheItem
	} else {
		res, err := http.Get(url)
		if err != nil {
			return locationArea{}, errors.New("Error connecting to endpoint.")
		}
		defer res.Body.Close()
		data, err = io.ReadAll(res.Body)
		if err != nil {
			return locationArea{}, err
		}
		pokeCache.Add(url, data)
	}
	locationAreas, err := parseLocationJSON(data)
	if err != nil {
		return locationArea{}, errors.New("JSON Parsing failed")
	}

	currentMapUrl = locationAreas.Next

	return locationAreas, nil
}

func clearScreen() {
	fmt.Fprintf(os.Stdout, PLACE_CURSOR)
}

func cleanInput(text string) []string {
	cleanTextSplit := strings.Split(strings.ToLower(strings.TrimSpace(text)), " ")

	fmt.Println(cleanTextSplit)

	return []string{}
}

func buildCommandHooks(rawHooks map[string]cliCommand) {
	if rawHooks == nil {
		return
	}
	rawHooks["exit"] = cliCommand{name: "exit", description: "Exit the Pokedex", callBack: commandExit}
	rawHooks["help"] = cliCommand{name: "help", description: "Displays a help message", callBack: commandHelp}
	rawHooks["map"] = cliCommand{name: "map", description: "Display the next location areas", callBack: commandMap}
	rawHooks["mapb"] = cliCommand{name: "mapb", description: "Display the prev location areas", callBack: commandMapB}
	rawHooks["explore"] = cliCommand{name: "explore", description: "Return the pokemon in a given area!", callBack: commandExplore}
}

func replLoop() {
	var line string
	scanner := bufio.NewScanner(os.Stdin)

	buildCommandHooks(commandHooks)

	for {
		fmt.Printf("PODEX9001 > ")
		scanner.Scan()
		line = scanner.Text()
		args := strings.Split(line, " ")
		arg_count := len(args)
		input_command := args[0]

		if command, ok := commandHooks[input_command]; ok {
			switch input_command {
			case "explore":
				if arg_count > 1 {
					command.callBack(args[1:]...)
				}
			default:
				command.callBack()
			}
		} else {
			fmt.Printf("Invalid command: %s\n%s", input_command, PLACE_CURSOR)
		}
	}
}

func main() {
	clearScreen()
	pokeCache = pc.NewCache(5 * time.Minute)
	replLoop()
}
