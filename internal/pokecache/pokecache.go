package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val 			[]byte
}

type Cache struct {
	mu sync.Mutex
	entries map[string]cacheEntry
}

func (c *Cache) Add(key string, val []byte)  {
	c.mu.Lock()
	newEntry := cacheEntry{createdAt: time.Now(), val: val}
	c.entries[key] = newEntry
	c.mu.Unlock()
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	entry, ok := c.entries[key]
	c.mu.Unlock()
	if !ok {
		return []byte{}, false
	}
	return entry.val, true
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			c.mu.Lock()
			for k, v := range c.entries {
				if time.Since(v.createdAt) > interval {
					delete(c.entries, k)
				}
			}
			c.mu.Unlock()
		}
	}()

}

func NewCache(interval time.Duration) (*Cache, error) {
	newCache := &Cache{
		entries: make(map[string]cacheEntry),
	}
	newCache.reapLoop(interval)
	return newCache, nil
}
