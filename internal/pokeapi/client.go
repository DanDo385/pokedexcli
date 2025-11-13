package pokeapi

import (
    "fmt"
    "io"
    "net/http"
)

const baseURL = "https://pokeapi.co/api/v2"

type Client struct {
    httpClient *http.Client
}

func NewClient() *Client {
    return &Client{
        httpClient: &http.Client{},
    }
}

func (c *Client) doGet(url string) ([]byte, error) {
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

    return body, nil
}
