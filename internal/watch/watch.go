package watch

import (
	"context"
	"sync"
	"time"
)

// SecretFetcher retrieves a map of secret key/value pairs.
type SecretFetcher interface {
	FetchSecrets(ctx context.Context) (map[string]string, error)
}

// ChangeHandler is called when a change in secrets is detected.
type ChangeHandler func(old, new map[string]string) error

// Watcher polls a SecretFetcher at a fixed interval and invokes
// ChangeHandler whenever the secret values differ from the previous fetch.
type Watcher struct {
	fetcher  SecretFetcher
	handler  ChangeHandler
	interval time.Duration
	mu       sync.Mutex
	last     map[string]string
}

// New creates a Watcher with the given fetcher, handler, and poll interval.
// If interval is zero or negative it defaults to 30 seconds.
func New(fetcher SecretFetcher, handler ChangeHandler, interval time.Duration) *Watcher {
	if interval <= 0 {
		interval = 30 * time.Second
	}
	return &Watcher{
		fetcher:  fetcher,
		handler:  handler,
		interval: interval,
	}
}

// Start begins polling until ctx is cancelled.
func (w *Watcher) Start(ctx context.Context) error {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := w.poll(ctx); err != nil {
				return err
			}
		}
	}
}

func (w *Watcher) poll(ctx context.Context) error {
	current, err := w.fetcher.FetchSecrets(ctx)
	if err != nil {
		return err
	}

	w.mu.Lock()
	old := w.last
	w.mu.Unlock()

	if !equal(old, current) {
		if err := w.handler(old, current); err != nil {
			return err
		}
		w.mu.Lock()
		w.last = current
		w.mu.Unlock()
	}
	return nil
}

func equal(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}
