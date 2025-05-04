package pokeapi

import (
	"fmt"
	"encoding/json"
	"net/http"
	"io"
  "github.com/mhiillos/pokedex-cli/internal/pokecache"
)

type LocationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Response struct {
	Results  []LocationArea `json:"results"`
	Next 		 string `json:"next"`
	Previous string `json:"previous"`
}

func Get(c *pokecache.Cache, endpoint string) (Response, error) {
	// Try to load from cache
	cached, ok := c.Get(endpoint)
	if ok {
		response, err := decodeResponse(cached)
		if err != nil {
			return Response{}, fmt.Errorf("Error unmarshaling cached data: %w", err)
		}
		return response, nil
	}

	// Otherwise, make a new request
	res, err := http.Get(endpoint)
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
		c.Add(endpoint, body)
		return response, nil
}

func decodeResponse(data []byte) (Response, error) {
	response := Response{}
	err := json.Unmarshal(data, &response)
	return response, err
}
