package main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strings"
)

const PLACE_CURSOR string = "\033[H\033[3J\033[80;1H"

func clearScreen() {
	fmt.Fprintf(os.Stdout, PLACE_CURSOR)
}

func cleanInput(text string) []string {
	cleanTextSplit := strings.Split(strings.ToLower(strings.TrimSpace(text)), " ")

	fmt.Println(cleanTextSplit)

	return []string{}
}

func replLoop() {
	var line, word string
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("PODEX9001 > ")
		scanner.Scan()
		line = scanner.Text()
		word = strings.Split(line, " ")[0]
		fmt.Printf("Your command was: %s\n", strings.ToLower(word))
	}
}

func main() {
	clearScreen()
	fmt.Printf("%s\n", reflect.TypeOf(os.Stdin))
	replLoop()
}
