package prefetch

import (
	"context"
	"sync"
	"time"
)

// entry is a cached secret result with an expiry time.
type entry struct {
	data    map[string]interface{}
	expires time.Time
}

// CachingFetcher wraps a Fetcher and caches results for ttl duration.
type CachingFetcher struct {
	upstream Fetcher
	ttl      time.Duration
	mu       sync.RWMutex
	cache    map[string]entry
}

// NewCachingFetcher returns a CachingFetcher backed by upstream.
// If ttl <= 0 caching is effectively disabled (entries expire immediately).
func NewCachingFetcher(upstream Fetcher, ttl time.Duration) *CachingFetcher {
	if upstream == nil {
		panic("prefetch: nil upstream fetcher")
	}
	return &CachingFetcher{
		upstream: upstream,
		ttl:      ttl,
		cache:    make(map[string]entry),
	}
}

// GetSecretData returns cached data when available and not expired,
// otherwise delegates to the upstream fetcher and stores the result.
func (c *CachingFetcher) GetSecretData(ctx context.Context, path string) (map[string]interface{}, error) {
	c.mu.RLock()
	e, ok := c.cache[path]
	c.mu.RUnlock()

	if ok && time.Now().Before(e.expires) {
		return e.data, nil
	}

	data, err := c.upstream.GetSecretData(ctx, path)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	c.cache[path] = entry{data: data, expires: time.Now().Add(c.ttl)}
	c.mu.Unlock()

	return data, nil
}

// Invalidate removes a single path from the cache.
func (c *CachingFetcher) Invalidate(path string) {
	c.mu.Lock()
	delete(c.cache, path)
	c.mu.Unlock()
}

// Flush clears the entire cache.
func (c *CachingFetcher) Flush() {
	c.mu.Lock()
	c.cache = make(map[string]entry)
	c.mu.Unlock()
}
