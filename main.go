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
	callBack    func() error
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

var commandHooks map[string]cliCommand = make(map[string]cliCommand)

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return errors.New("WTF, os.Exit failed...")
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:\n")

	for command := range commandHooks {
		fmt.Printf("%s: %s\n%s", commandHooks[command].name, commandHooks[command].description, PLACE_CURSOR)
	}
	return nil
}

func commandMapB() error {
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

func commandMap() error {
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

func fillLocationArea(url string) (locationArea, error) {
	var data []byte

	if url == "" {
		url = "https://pokeapi.co/api/v2/location-area/"
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
}

func replLoop() {
	var line, word string
	scanner := bufio.NewScanner(os.Stdin)

	buildCommandHooks(commandHooks)

	for {
		fmt.Printf("PODEX9001 > ")
		scanner.Scan()
		line = scanner.Text()
		word = strings.Split(line, " ")[0]
		if command, ok := commandHooks[word]; ok {
			command.callBack()
		} else {
			fmt.Printf("Invalid command: %s\n%s", word, PLACE_CURSOR)
		}
	}
}

func main() {
	clearScreen()
	pokeCache = pc.NewCache(5 * time.Minute)
	replLoop()
}
