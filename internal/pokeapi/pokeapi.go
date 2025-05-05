package pokeapi

import (
	"fmt"
	"encoding/json"
	"net/http"
	"io"
  "github.com/mhiillos/pokedex-cli/internal/pokecache"
)

type Client struct {
    HTTP *http.Client
    Cache *pokecache.Cache
}

type LocationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Response struct {
	Results  []LocationArea `json:"results"`
	Next 		 string `json:"next"`
	Previous string `json:"previous"`
}

func (client *Client) Get(endpoint string) (Response, error) {
	// Try to load from cache
	cached, ok := client.Cache.Get(endpoint)
	if ok {
		response, err := decodeResponse(cached)
		if err != nil {
			return Response{}, fmt.Errorf("Error unmarshaling cached data: %w", err)
		}
		return response, nil
	}

	// Otherwise, make a new request
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return Response{}, fmt.Errorf("Error creating new request: %w", err)
	}
	res, err := client.HTTP.Do(req)
	if err != nil {
		return Response{}, fmt.Errorf("Error fetching endpoint: %w", err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if res.StatusCode > 299 {
		return Response{}, fmt.Errorf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		return Response{}, fmt.Errorf("%w", err)
	}
	response, err := decodeResponse(body)
	if err != nil {
		return Response{}, fmt.Errorf("Error unmarshaling body: %w", err)
	}

	// Cache the data
	client.Cache.Add(endpoint, body)
	return response, nil
}

func decodeResponse(data []byte) (Response, error) {
	response := Response{}
	err := json.Unmarshal(data, &response)
	return response, err
}
