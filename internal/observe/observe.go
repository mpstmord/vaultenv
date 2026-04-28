// Package observe provides a lightweight hook system for attaching
// observers to secret fetch operations, enabling metrics, logging,
// and tracing integrations without coupling the core fetcher.
package observe

import (
	"context"
	"time"
)

// Event describes the outcome of a single secret fetch attempt.
type Event struct {
	// Path is the Vault secret path that was requested.
	Path string
	// Field is the specific key requested within the secret.
	Field string
	// Duration is how long the fetch took.
	Duration time.Duration
	// Err is non-nil when the fetch failed.
	Err error
	// Cached is true when the result was served from a local cache.
	Cached bool
}

// Observer is called after each fetch attempt with the resulting Event.
type Observer func(ctx context.Context, e Event)

// Fetcher is the interface satisfied by all secret providers.
type Fetcher interface {
	GetSecretData(ctx context.Context, path string) (map[string]interface{}, error)
}

// observedFetcher wraps a Fetcher and notifies registered Observers.
type observedFetcher struct {
	upstream  Fetcher
	observers []Observer
}

// New returns a Fetcher that calls each Observer after every fetch.
// Observers are invoked synchronously in registration order.
func New(upstream Fetcher, observers ...Observer) Fetcher {
	if upstream == nil {
		panic("observe: upstream fetcher must not be nil")
	}
	obs := make([]Observer, len(observers))
	copy(obs, observers)
	return &observedFetcher{upstream: upstream, observers: obs}
}

func (o *observedFetcher) GetSecretData(ctx context.Context, path string) (map[string]interface{}, error) {
	start := time.Now()
	data, err := o.upstream.GetSecretData(ctx, path)
	e := Event{
		Path:     path,
		Duration: time.Since(start),
		Err:      err,
	}
	for _, obs := range o.observers {
		obs(ctx, e)
	}
	return data, err
}
