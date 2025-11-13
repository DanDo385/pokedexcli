package pokeapi

import (
    "encoding/json"
    "fmt"
)

type LocationAreasResponse struct {
    Count    int `json:"count"`
    Next     *string `json:"next"`
    Previous *string `json:"previous"`
    Results  []struct {
        Name string `json:"name"`
        URL  string `json:"url"`
    } `json:"results"`
}

// GetLocationAreas fetches a page of 20 location areas.
func (c *Client) GetLocationAreas(pageURL *string) (LocationAreasResponse, error) {
    var url string
    if pageURL != nil {
        url = *pageURL
    } else {
        url = fmt.Sprintf("%s/location-area", baseURL)
    }

    body, err := c.doGet(url)
    if err != nil {
        return LocationAreasResponse{}, err
    }

    var data LocationAreasResponse
    err = json.Unmarshal(body, &data)
    if err != nil {
        return LocationAreasResponse{}, err
    }

    return data, nil
}
