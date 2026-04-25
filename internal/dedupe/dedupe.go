// Package dedupe provides a single-flight deduplication layer for secret
// fetches. Concurrent requests for the same path are collapsed into a single
// upstream call; all waiters receive the same result.
package dedupe

import (
	"context"
	"sync"
)

// Fetcher is the interface satisfied by any upstream secret client.
type Fetcher interface {
	GetSecretData(ctx context.Context, path string) (map[string]interface{}, error)
}

// call tracks an in-flight or completed fetch.
type call struct {
	wg  sync.WaitGroup
	val map[string]interface{}
	err error
}

// Deduplicator wraps a Fetcher and collapses concurrent identical requests.
type Deduplicator struct {
	upstream Fetcher
	mu       sync.Mutex
	inflight map[string]*call
}

// New returns a Deduplicator backed by upstream.
// It panics if upstream is nil.
func New(upstream Fetcher) *Deduplicator {
	if upstream == nil {
		panic("dedupe: upstream Fetcher must not be nil")
	}
	return &Deduplicator{
		upstream: upstream,
		inflight: make(map[string]*call),
	}
}

// GetSecretData fetches the secret at path. If another goroutine is already
// fetching the same path the caller blocks and shares the result.
func (d *Deduplicator) GetSecretData(ctx context.Context, path string) (map[string]interface{}, error) {
	d.mu.Lock()
	if c, ok := d.inflight[path]; ok {
		d.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}

	c := &call{}
	c.wg.Add(1)
	d.inflight[path] = c
	d.mu.Unlock()

	c.val, c.err = d.upstream.GetSecretData(ctx, path)
	c.wg.Done()

	d.mu.Lock()
	delete(d.inflight, path)
	d.mu.Unlock()

	return c.val, c.err
}
