package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func cleanInput(text string) []string {
	lower := strings.ToLower(text)
	trimmed := strings.TrimSpace(lower)
	words := strings.Fields(trimmed)

	return words
}

func repl(config *config) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("pokedex > ")
		scanner.Scan()
		input := scanner.Text()

		args := cleanInput(input)
		cmdInput := args[0]
		if len(args) > 1 {
			config.item = args[1]
		}

		cmd, ok := getCommands()[cmdInput]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		if err := cmd.callback(config); err != nil {
			fmt.Println(err)
		}
	}
}
