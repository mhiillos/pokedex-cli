package pokeapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mhiillos/pokedex-cli/internal/pokecache"
)

func TestAPIGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"results":[{"name":"test-area","url":"test-url"}],"next":"next-url","previous":"previous-url"}`))
	}))
	defer server.Close()

	cache, err := pokecache.NewCache(5000 * time.Millisecond)
	if err != nil {
    t.Fatalf("failed to create cache: %v", err)
	}

	client := Client{
		HTTP: server.Client(),
		Cache: cache,
	}

	resp, err := client.Get(server.URL)
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
		t.Errorf("Expected previous url 'previous-url', got %q", resp.Previous)
	}
	if resp.Next != "next-url" {
		t.Errorf("Expected previous url 'next-url', got %q", resp.Next)
	}

}
