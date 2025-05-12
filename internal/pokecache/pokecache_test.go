package internal

import (
	"testing"
	"time"
)

func TestCacheOps(t *testing.T) {
	cache := NewCache(2 * time.Second)

	value := "first value"
	cache.Add("first", []byte(value))
	if len(cache.cache) != 1 {
		t.Errorf("Expected cache len 1. Got %d", len(cache.cache))
	}

	entry, found := cache.Get("first")
	if !found {
		t.Errorf("Inserted values was not found in cache")
	}

	if string(entry) != value {
		t.Errorf("Expected entry value %s. Got %s", value, string(entry))
	}

	time.Sleep(2 * time.Second) // wait for entry to expire

	entries := len(cache.cache)

	if entries != 0 {
		t.Errorf("cache entry should have expired")
	}

}
