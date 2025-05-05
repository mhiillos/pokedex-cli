package pokeapi

import (
	"fmt"
	"errors"
	"encoding/json"
	"math/rand"
	"net/http"
	"io"

  "github.com/mhiillos/pokedex-cli/internal/pokecache"
)

type HTTPError struct {
	StatusCode  int
	Body      	string
}

func (e *HTTPError) Error() string {
    return fmt.Sprintf("Response failed with status code: %d and body: %s", e.StatusCode, e.Body)
}

type Client struct {
    HTTP *http.Client
    Cache *pokecache.Cache
}

type LocationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type LocationAreasResponse struct {
	Results  []LocationArea `json:"results"`
	Next 		 string 				`json:"next"`
	Previous string 				`json:"previous"`
}

type BasicPokemon struct {
	Name string `json:"name"`
	URL string  `json:"url"`
}

type PokemonEncounter struct {
	Pokemon BasicPokemon `json:"pokemon"`
}

type Pokemon struct {
	Name string        `json:"name"`
	BaseExperience int `json:"base_experience"`
}

type LocationAreaResponse struct {
	Name string 								  			 `json:"name"`
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}


// Does a raw get, returns raw bytestream and caches the result
func (client *Client) Get(endpoint string) ([]byte, error) {
	// Try to load from cache
	cached, ok := client.Cache.Get(endpoint)
	if ok {
		return cached, nil
	}

	// Otherwise, make a new request
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return []byte{}, fmt.Errorf("Error creating new request: %w", err)
	}
	res, err := client.HTTP.Do(req)
	if err != nil {
		return []byte{}, fmt.Errorf("Error fetching endpoint: %w", err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if res.StatusCode > 299 {
		return []byte{}, &HTTPError{StatusCode: res.StatusCode, Body: string(body)}
	}
	if err != nil {
		return []byte{}, &HTTPError{StatusCode: res.StatusCode, Body: string(body)}
	}

	// Cache the data
	client.Cache.Add(endpoint, body)
	return body, nil
}

// This function fetches the 20 location areas for the endpoint
func (client *Client) GetLocationAreas(endpoint string) (LocationAreasResponse, error) {
	body, err := client.Get(endpoint)
	if err != nil {
		return LocationAreasResponse{}, err
	}

	response := LocationAreasResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return LocationAreasResponse{}, err
	}
	return response, nil
}

// This function gets the list of pokemon in an area
func (client *Client) ExploreLocationArea(endpoint string) (LocationAreaResponse, error) {
	body, err := client.Get(endpoint)
	if err != nil {
		var httpErr *HTTPError
		if errors.As(err, &httpErr) {
			if httpErr.StatusCode == 404 {
				return LocationAreaResponse{}, fmt.Errorf("Location area not found")
			}
		}
		return LocationAreaResponse{}, err
	}

	response := LocationAreaResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return LocationAreaResponse{}, err
	}
	return response, nil
}

func (client *Client) RollPokemon(endpoint string) (Pokemon, error) {
	body, err := client.Get(endpoint)
	if err != nil {
		var httpErr *HTTPError
		if errors.As(err, &httpErr) {
			if httpErr.StatusCode == 404 {
				return Pokemon{}, fmt.Errorf("Pokemon not found")
			}
		}
	}

	pokemon := Pokemon{}
	err = json.Unmarshal(body, &pokemon)
	if err != nil {
		return Pokemon{}, err
	}
	// Roll for catch
	baseXP := pokemon.BaseExperience
	roll := rand.Intn(baseXP)
	if roll > (baseXP * baseXP + 2)/(2 * baseXP) {
		return pokemon, nil
	}
	return Pokemon{}, nil
}
