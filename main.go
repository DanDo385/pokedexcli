package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
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

		fmt.Printf("Your command was: %s\n", first)
	}
}