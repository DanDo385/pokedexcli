package main

import "github.com/DanDo385/pokedexcli/internal/pokeapi"

type config struct {
	nextURL *string
	prevURL *string
	client  *pokeapi.Client
}