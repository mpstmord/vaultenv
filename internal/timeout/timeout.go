// Package timeout provides context-based deadline enforcement for secret
// fetch operations, ensuring individual Vault calls do not block indefinitely.
package timeout

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// DefaultTimeout is used when no explicit duration is configured.
const DefaultTimeout = 10 * time.Second

// ErrDeadlineExceeded is returned when an operation exceeds its allowed duration.
var ErrDeadlineExceeded = errors.New("timeout: deadline exceeded")

// Fetcher is the interface satisfied by any secret-fetching component.
type Fetcher interface {
	GetSecretData(ctx context.Context, path string) (map[string]interface{}, error)
}

// GuardedFetcher wraps a Fetcher and enforces a per-call timeout.
type GuardedFetcher struct {
	upstream Fetcher
	duration time.Duration
}

// New returns a GuardedFetcher that cancels calls to upstream after d.
// If d is <= 0, DefaultTimeout is used.
func New(upstream Fetcher, d time.Duration) (*GuardedFetcher, error) {
	if upstream == nil {
		return nil, errors.New("timeout: upstream fetcher must not be nil")
	}
	if d <= 0 {
		d = DefaultTimeout
	}
	return &GuardedFetcher{upstream: upstream, duration: d}, nil
}

// GetSecretData calls the upstream fetcher with a bounded context.
// If the deadline is exceeded, ErrDeadlineExceeded is returned.
func (g *GuardedFetcher) GetSecretData(ctx context.Context, path string) (map[string]interface{}, error) {
	tctx, cancel := context.WithTimeout(ctx, g.duration)
	defer cancel()

	data, err := g.upstream.GetSecretData(tctx, path)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("%w: path %q exceeded %s", ErrDeadlineExceeded, path, g.duration)
		}
		return nil, err
	}
	return data, nil
}

// Duration returns the configured timeout duration.
func (g *GuardedFetcher) Duration() time.Duration {
	return g.duration
}
