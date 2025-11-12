package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)
type cliCommand struct { // Define a struct for the CLI commands
	name        string
	description string
	callback    func() error
}

func commandExit() error { // Command to exit the Pokedex
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil // unreachable, but required
}

func commandHelp() error { // Command to display the help message
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")

	for _, cmd := range commands { // Loop through the commands and print the name and description
		fmt.Printf("%s: %s\n", cmd.name, cmd.description) // Print the name and description of the command
	}
	return nil // Return nil to indicate success
}

var commands = map[string]cliCommand{}

func init() {
	commands["exit"] = cliCommand{
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	}
	commands["help"] = cliCommand{
		name:        "help",
		description: "Displays a help message",
		callback:    commandHelp,
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")

		if !scanner.Scan() {
			break
		}

		input := strings.ToLower(strings.TrimSpace(scanner.Text()))
		if input == "" {
			continue // If the input is empty, continue to the next iteration of the loop
		}

		parts := strings.Fields(input) // Split the input into parts
		cmdName := parts[0] // Get the command name

		if cmd, exists := commands[cmdName]; exists {
			err := cmd.callback()
			if err != nil {
				fmt.Println("Error:", err)
			}
		} else { // If the command does not exist, print an error message
			fmt.Println("Unknown command")
		}
	}
}
