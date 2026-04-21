package drain

import (
	"context"
	"errors"
)

// ErrDraining is returned when a secret fetch is attempted during shutdown.
var ErrDraining = errors.New("drain: shutdown in progress, request rejected")

// Fetcher is the interface satisfied by secret-fetching middleware.
type Fetcher interface {
	GetSecretData(ctx context.Context, path string) (map[string]interface{}, error)
}

// GuardedFetcher wraps an upstream Fetcher with drain protection.
// Requests are rejected once Drain has been called on the Drainer.
type GuardedFetcher struct {
	upstream Fetcher
	drainer  *Drainer
}

// NewGuardedFetcher returns a GuardedFetcher that rejects fetches
// when the provided Drainer is in draining state.
func NewGuardedFetcher(upstream Fetcher, d *Drainer) *GuardedFetcher {
	return &GuardedFetcher{upstream: upstream, drainer: d}
}

// GetSecretData forwards to the upstream fetcher only if the Drainer
// has not yet been closed. It holds an acquire slot for the lifetime
// of the upstream call so Drain waits for it to complete.
func (g *GuardedFetcher) GetSecretData(ctx context.Context, path string) (map[string]interface{}, error) {
	if !g.drainer.Acquire() {
		return nil, ErrDraining
	}
	defer g.drainer.Release()
	return g.upstream.GetSecretData(ctx, path)
}
