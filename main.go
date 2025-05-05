package main

import (
	"fmt"
	"strings"
)

func cleanInput(text string) []string {
	lower := strings.ToLower(text)
	trimmed := strings.TrimSpace(lower)
	words := strings.Fields(trimmed)

	return words
}

func main() {
	fmt.Println("Hello, World!")
}
