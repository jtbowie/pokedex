package main

import (
	"bufio"
	"errors"
	"fmt"
	pc "github.com/jtbowie/pokedex/internal/pokecache"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

var currentMapUrl = ""
var currentLocationArea locationAreaJSON
var currentEncounter PokemonEncounterJSON
var pokeCache *pc.Cache
var pokeDex map[string]PokemonJSON = make(map[string]PokemonJSON)

const PLACE_CURSOR string = "\033[H\033[3J\033[80;1H"
const BASE_URL string = "https://pokeapi.co/api/v2/location-area/"

var commandHooks map[string]cliCommand = make(map[string]cliCommand)

func (pEJ *PokemonJSON) fill(url string) error {
	var data []byte

	if url == "" {
		return errors.New("enter a url dude")
	}
	if cacheItem, ok := pokeCache.Get(url); ok {
		data = cacheItem
	} else {
		res, err := http.Get(url)
		if err != nil {
			return errors.New("error connecting to endpoint")
		}
		defer res.Body.Close()
		data, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		pokeCache.Add(url, data)
	}

	var pokeJSONObj JSONObj = pEJ

	err := pokeJSONObj.parseJSON(data)
	if err != nil {
		return errors.New("JSON Parsing failed")
	}

	return nil
}

func (pEJ *PokemonEncounterJSON) fill(url string) error {
	var data []byte

	if url == "" {
		return errors.New("enter a url dude")
	}
	if cacheItem, ok := pokeCache.Get(url); ok {
		data = cacheItem
	} else {
		res, err := http.Get(url)
		if err != nil {
			return errors.New("error connecting to endpoint")
		}
		defer res.Body.Close()
		data, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		pokeCache.Add(url, data)
	}

	var pokeJSONObj JSONObj = pEJ

	err := pokeJSONObj.parseJSON(data)
	if err != nil {
		return errors.New("JSON Parsing failed")
	}
	currentEncounter = *pEJ

	return nil
}

func (lAJ *locationAreaJSON) fill(url string) error {
	var data []byte

	if url == "" {
		url = BASE_URL
	}

	if cacheItem, ok := pokeCache.Get(url); ok {
		data = cacheItem
	} else {
		res, err := http.Get(url)
		if err != nil {
			return errors.New("error connecting to endpoint")
		}
		defer res.Body.Close()
		data, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		pokeCache.Add(url, data)
	}

	var localjsonObj JSONObj = lAJ

	err := localjsonObj.parseJSON(data)
	if err != nil {
		return errors.New("JSON Parsing failed")
	}

	currentMapUrl = lAJ.Next

	return nil
}

func clearScreen() {
	_, err := fmt.Fprint(os.Stdout, PLACE_CURSOR)
	if err != nil {
		fmt.Println("Error clearing screen")
	}
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
	rawHooks["catch"] = cliCommand{name: "catch", description: "Catch a pokemon in a given area!", callBack: commandCatch}
	rawHooks["inspect"] = cliCommand{name: "inspect", description: "Inspect a pokemon!", callBack: commandInspect}

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
			case "catch":
				if arg_count > 1 {
					command.callBack(args[1:]...)
				}
			case "inspect":
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
