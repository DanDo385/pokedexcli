package pokeapi

import (
    "fmt"
    "io"
    "net/http"
    "time"

    "github.com/DanDo385/pokedexcli/internal/pokecache"
)

const baseURL = "https://pokeapi.co/api/v2"
type Client struct {
    httpClient *http.Client
    cache      *pokecache.Cache
}

func NewClient() *Client {
    return &Client{
        httpClient: &http.Client{},
        cache:      pokecache.NewCache(5 * time.Second),
    }
}

func (c *Client) doGet(url string) ([]byte, error) {
    // Check cache first
    if val, ok := c.cache.Get(url); ok {
        return val, nil
    }

    // Make HTTP request
    resp, err := c.httpClient.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("bad status: %s", resp.Status)
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    // Add to cache
    c.cache.Add(url, body)

    return body, nil
}
