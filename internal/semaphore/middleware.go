package semaphore

import (
	"context"
	"fmt"
)

// Fetcher retrieves secret data by path and field.
type Fetcher interface {
	GetSecretData(ctx context.Context, path string) (map[string]interface{}, error)
}

// GuardedFetcher wraps a Fetcher and limits concurrent calls via a Semaphore.
type GuardedFetcher struct {
	upstream Fetcher
	sem      *Semaphore
}

// NewGuardedFetcher creates a GuardedFetcher that allows at most concurrency
// simultaneous calls to the upstream Fetcher.
// It panics if upstream is nil or concurrency < 1.
func NewGuardedFetcher(upstream Fetcher, concurrency int) *GuardedFetcher {
	if upstream == nil {
		panic("semaphore: upstream fetcher must not be nil")
	}
	return &GuardedFetcher{
		upstream: upstream,
		sem:      New(concurrency),
	}
}

// GetSecretData acquires a semaphore slot before forwarding to upstream.
// Returns an error if the context is cancelled while waiting.
func (g *GuardedFetcher) GetSecretData(ctx context.Context, path string) (map[string]interface{}, error) {
	if err := g.sem.Acquire(ctx); err != nil {
		return nil, fmt.Errorf("semaphore acquire: %w", err)
	}
	defer g.sem.Release()
	return g.upstream.GetSecretData(ctx, path)
}

// Available returns the number of free concurrency slots.
func (g *GuardedFetcher) Available() int {
	return g.sem.Available()
}
