package internal

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	data      []byte
}

type Cache struct {
	interval time.Duration
	cache    map[string]cacheEntry
	mu       sync.Mutex
}

func NewCache(interval time.Duration) *Cache {
	self := &Cache{
		interval: interval,
		cache:    make(map[string]cacheEntry),
		mu:       sync.Mutex{},
	}
	go self.reapLoop()

	return self
}

func (c *Cache) Add(key string, data []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[key] = cacheEntry{
		data:      data,
		createdAt: time.Now(),
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, found := c.cache[key]
	if !found {
		return nil, found
	}

	return entry.data, found
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	for t := range ticker.C {
		c.mu.Lock()

		for k, v := range c.cache {
			if v.createdAt.Compare(t) == -1 {
				delete(c.cache, k)
			}

		}

		c.mu.Unlock()
	}
}
