package main

import (
	"fmt"
	"strings"
)

func cleanInput(text string) []string {
	cleanTextSplit := strings.Split(strings.ToLower(strings.TrimSpace(text)), " ")

	fmt.Println(cleanTextSplit)

	return []string{}
}

func main() {
	fmt.Println(cleanInput("Hello, World!"))
}
