package rotate

import (
	"context"
	"fmt"
	"time"
)

// SecretFetcher retrieves secret data by path.
type SecretFetcher interface {
	GetSecretData(ctx context.Context, path string) (map[string]interface{}, error)
}

// Handler is called when a secret is rotated.
type Handler func(path string, data map[string]interface{}) error

// Rotator polls Vault for secret changes and invokes a handler on updates.
type Rotator struct {
	fetcher  SecretFetcher
	interval time.Duration
	handler  Handler
	paths    []string
	cache    map[string]string
}

// New creates a Rotator that watches the given paths.
func New(fetcher SecretFetcher, interval time.Duration, handler Handler, paths []string) *Rotator {
	if interval <= 0 {
		interval = 30 * time.Second
	}
	return &Rotator{
		fetcher:  fetcher,
		interval: interval,
		handler:  handler,
		paths:    paths,
		cache:    make(map[string]string),
	}
}

// Start begins polling until ctx is cancelled.
func (r *Rotator) Start(ctx context.Context) error {
	if err := r.poll(ctx); err != nil {
		return err
	}
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			_ = r.poll(ctx)
		}
	}
}

func (r *Rotator) poll(ctx context.Context) error {
	for _, path := range r.paths {
		data, err := r.fetcher.GetSecretData(ctx, path)
		if err != nil {
			return fmt.Errorf("rotate: fetch %q: %w", path, err)
		}
		sig := signature(data)
		if prev, ok := r.cache[path]; !ok || prev != sig {
			r.cache[path] = sig
			if err := r.handler(path, data); err != nil {
				return fmt.Errorf("rotate: handler %q: %w", path, err)
			}
		}
	}
	return nil
}

func signature(data map[string]interface{}) string {
	h := fmt.Sprintf("%v", data)
	return h
}
