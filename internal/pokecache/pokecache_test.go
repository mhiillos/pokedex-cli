package pokecache

import (
	"bytes"
	"testing"
	"time"
)

func newTestCache(t *testing.T) (*Cache) {
	cache, err := NewCache(100 * time.Millisecond)
	if err != nil {
		t.Fatalf("Failed creating cache: %q", err)
	}
	return cache
}


func TestCacheAddGet(t *testing.T) {
	cache := newTestCache(t)
	cache.Add("testKey", []byte("testVal1"))
	val, ok := cache.Get("testKey")
	if !ok {
		t.Fatalf("key `testKey` not found")
	}
	if !bytes.Equal(val, []byte("testVal1")) {
		t.Errorf("Values do not match. (Expected 'testVal1', got %q)", val)
	}
}

func TestCacheDelete(t *testing.T) {
	cache := newTestCache(t)
	cache.Add("testKey", []byte("testVal1"))
	time.Sleep(200 * time.Millisecond)
	_, ok := cache.Get("testKey")
	if ok {
		t.Errorf("Key `testKey` did not get removed after expiry.")
	}
}

