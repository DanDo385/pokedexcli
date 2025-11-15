package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"

    "github.com/DanDo385/pokedexcli/internal/pokeapi"
)

type cliCommand struct {
    name        string // name of the command
    description string // description of the command
    callback    func(*config) error // function that takes a config and returns an error
 }

var commands = map[string]cliCommand{} // map of command name to command

func registerCommand(name, description string, callback func(*config) error) {
    commands[name] = cliCommand{
        name:        name,
        description: description,
        callback:    callback,
    }
}

func init() {
    registerCommand("help", "Displays a help message", commandHelp)
    registerCommand("exit", "Exit the Pokedex", commandExit)
    registerCommand("map", "List the next 20 location areas", commandMap)
    registerCommand("mapb", "List the previous 20 location areas", commandMapb)
}

func main() {
    scanner := bufio.NewScanner(os.Stdin)
    cfg := &config{
        client: pokeapi.NewClient(),
    }

    for {
        fmt.Print("Pokedex > ")
        if !scanner.Scan() {
            break
        }

        input := strings.ToLower(strings.TrimSpace(scanner.Text()))
        if input == "" {
            continue
        }

        parts := strings.Fields(input)
        cmdName := parts[0]

        if cmd, exists := commands[cmdName]; exists {
            err := cmd.callback(cfg)
            if err != nil {
                fmt.Println("Error:", err)
            }
        } else {
            fmt.Println("Unknown command")
        }
    }
}

// ===== Commands =====

func commandExit(cfg *config) error {
    fmt.Println("Closing the Pokedex... Goodbye!")
    os.Exit(0)
    return nil
}

func commandHelp(cfg *config) error {
    fmt.Println("Welcome to the Pokedex!")
    fmt.Println("Usage:")
    for _, cmd := range commands {
        fmt.Printf("%s: %s\n", cmd.name, cmd.description)
    }
    return nil
}

func displayLocationAreas(cfg *config, url *string) error {
    data, err := cfg.client.GetLocationAreas(url)
    if err != nil {
        return err
    }

    for _, area := range data.Results {
        fmt.Println(area.Name)
    }

    cfg.nextURL = data.Next
    cfg.prevURL = data.Previous

    return nil
}

func commandMap(cfg *config) error {
    return displayLocationAreas(cfg, cfg.nextURL)
}

func commandMapb(cfg *config) error {
    if cfg.prevURL == nil {
        fmt.Println("you're on the first page")
        return nil
    }
    return displayLocationAreas(cfg, cfg.prevURL)
}
