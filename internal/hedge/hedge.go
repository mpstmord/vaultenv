// Package hedge implements a hedged request pattern for secret fetching.
// A hedged request fires a second parallel request after a short delay if
// the first has not yet returned, returning whichever response arrives first.
package hedge

import (
	"context"
	"time"

	"github.com/your-org/vaultenv/internal/vault"
)

// DefaultDelay is the hedging delay used when none is specified.
const DefaultDelay = 50 * time.Millisecond

// Fetcher is the interface satisfied by any upstream secret source.
type Fetcher interface {
	GetSecretData(ctx context.Context, path string) (map[string]interface{}, error)
}

// Hedger wraps an upstream Fetcher and issues a second request after Delay
// if the first has not completed, returning whichever result arrives first.
type Hedger struct {
	upstream Fetcher
	delay    time.Duration
}

type result struct {
	data map[string]interface{}
	err  error
}

// New returns a Hedger wrapping upstream. If delay is <= 0 DefaultDelay is used.
func New(upstream Fetcher, delay time.Duration) *Hedger {
	if upstream == nil {
		panic("hedge: upstream fetcher must not be nil")
	}
	if delay <= 0 {
		delay = DefaultDelay
	}
	return &Hedger{upstream: upstream, delay: delay}
}

// GetSecretData satisfies the Fetcher interface. It launches the first request
// immediately and, after h.delay, launches a second. The first result to
// arrive (success or error) is returned to the caller.
func (h *Hedger) GetSecretData(ctx context.Context, path string) (map[string]interface{}, error) {
	ch := make(chan result, 2)

	fire := func() {
		data, err := h.upstream.GetSecretData(ctx, path)
		ch <- result{data: data, err: err}
	}

	go fire()

	select {
	case res := <-ch:
		return res.data, res.err
	case <-time.After(h.delay):
		go fire()
	}

	select {
	case res := <-ch:
		return res.data, res.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// compile-time interface check
var _ Fetcher = (*Hedger)(nil)
var _ vault.Fetcher = (*Hedger)(nil)
