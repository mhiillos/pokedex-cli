package pokeapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mhiillos/pokedex-cli/internal/pokecache"
)

func CreateMockServer(body []byte) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
    w.Write(body)
	}))
	return server
}

func TestGetLocationAreas(t *testing.T) {
	server := CreateMockServer([]byte(`{"results":[{"name":"test-area","url":"test-url"}],"next":"next-url","previous":"previous-url"}`))
	defer server.Close()

	cache, err := pokecache.NewCache(5000 * time.Millisecond)
	if err != nil {
    t.Fatalf("failed to create cache: %v", err)
	}

	client := Client{
		HTTP: server.Client(),
		Cache: cache,
	}

	resp, err := client.GetLocationAreas(server.URL)
	if err != nil {
		t.Fatalf("GET failed: %v", err)
	}

	if len(resp.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(resp.Results))
	}
	if resp.Results[0].Name != "test-area" {
		t.Errorf("expected Name 'test-area', got %q", resp.Results[0].Name)
	}
	if resp.Results[0].URL != "test-url" {
		t.Errorf("expected URL 'test-url', got %q", resp.Results[0].URL)
	}
	if resp.Previous != "previous-url" {
		t.Errorf("expected previous url 'previous-url', got %q", resp.Previous)
	}
	if resp.Next != "next-url" {
		t.Errorf("expected previous url 'next-url', got %q", resp.Next)
	}
}

func TestExploreLocationArea(t *testing.T) {
	server := CreateMockServer([]byte(`{
		"name": "test_area",
		"pokemon_encounters": [
			{
				"pokemon": {
					"name": "pokemon_1",
					"url": "pokemon_1_url"
				}
			},
			{
				"pokemon": {
					"name": "pokemon_2",
					"url": "pokemon_2_url"
				}
			}
		]
	}`))
	defer server.Close()

	cache, err := pokecache.NewCache(5000 * time.Millisecond)
	if err != nil {
    t.Fatalf("failed to create cache: %v", err)
	}

	client := Client{
		HTTP: server.Client(),
		Cache: cache,
	}

	resp, err := client.ExploreLocationArea(server.URL)
	if err != nil {
		t.Fatalf("GET failed: %v", err)
	}

	if resp.Name != "test_area" {
		t.Errorf("expected name `test_area`, got %q", resp.Name)
	}
	if len(resp.PokemonEncounters) != 2 {
		t.Fatal("amount of pokemon do not match")
	}
	if resp.PokemonEncounters[0].Pokemon.Name != "pokemon_1" {
		t.Errorf("expected name `pokemon_1`, got %q", resp.PokemonEncounters[0].Pokemon.Name)
	}
}
