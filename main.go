package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)
type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil // unreachable, but required
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")

	
	scanner := bufio.NewScanner(os.Stdin)

	for {
		// Print the prompt without a newline
		fmt.Print("Pokedex > ")
		// Read the input from the user
		if !scanner.Scan() {
			break // user hit ctrl+c or input ended
		}
		// Get the input from the user
		input := scanner.Text() 
		// Clean the input
		input = strings.ToLower(strings.TrimSpace(input))
		// If input is empty, skip back to the start of the loop
		if input == "" {
			continue
		}

		// Split into words
		parts := strings.Fields(input)

		// First word is the command
		first := parts[0]

	}
}