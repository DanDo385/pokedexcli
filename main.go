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

var commands = map[string]cliCommand{}

func init() {
    commands["help"] = cliCommand{
        name:        "help",
        description: "Displays a help message",
        callback:    commandHelp,
    }

    commands["exit"] = cliCommand{
        name:        "exit",
        description: "Exit the Pokedex",
        callback:    commandExit,
    }

    commands["map"] = cliCommand{
        name:        "map",
        description: "List the next 20 location areas",
        callback:    commandMap,
    }

    commands["mapb"] = cliCommand{
        name:        "mapb",
        description: "List the previous 20 location areas",
        callback:    commandMapb,
    }
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

func commandMap(cfg *config) error {
    data, err := cfg.client.GetLocationAreas(cfg.nextURL)
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

func commandMapb(cfg *config) error {
    if cfg.prevURL == nil {
        fmt.Println("you're on the first page")
        return nil
    }

    data, err := cfg.client.GetLocationAreas(cfg.prevURL)
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
