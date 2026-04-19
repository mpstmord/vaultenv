// Package cache provides a simple in-memory TTL cache for Vault secrets.
package cache

import (
	"sync"
	"time"
)

// Entry holds a cached secret payload and its expiry time.
type Entry struct {
	Data      map[string]interface{}
	ExpiresAt time.Time
}

// Cache is a thread-safe in-memory store for secret data.
type Cache struct {
	mu      sync.RWMutex
	items   map[string]Entry
	ttl     time.Duration
}

// New creates a Cache with the given TTL duration.
func New(ttl time.Duration) *Cache {
	return &Cache{
		items: make(map[string]Entry),
		ttl:   ttl,
	}
}

// Set stores secret data under the given key.
func (c *Cache) Set(key string, data map[string]interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = Entry{
		Data:      data,
		ExpiresAt: time.Now().Add(c.ttl),
	}
}

// Get retrieves secret data if present and not expired.
func (c *Cache) Get(key string) (map[string]interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.items[key]
	if !ok || time.Now().After(entry.ExpiresAt) {
		return nil, false
	}
	return entry.Data, true
}

// Delete removes an entry from the cache.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Purge removes all expired entries from the cache.
func (c *Cache) Purge() {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	for k, v := range c.items {
		if now.After(v.ExpiresAt) {
			delete(c.items, k)
		}
	}
}
