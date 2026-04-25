// Package expire provides utilities for tracking and enforcing
// expiration deadlines on fetched secrets.
package expire

import (
	"context"
	"errors"
	"sync"
	"time"
)

// ErrExpired is returned when a secret has passed its expiration deadline.
var ErrExpired = errors.New("expire: secret has expired")

// Entry holds a value together with its expiration time.
type Entry struct {
	Value     map[string]any
	ExpiresAt time.Time
}

// Expired reports whether the entry has passed its deadline.
func (e Entry) Expired() bool {
	return time.Now().After(e.ExpiresAt)
}

// TTL returns the remaining lifetime of the entry. It may be negative
// if the entry has already expired.
func (e Entry) TTL() time.Duration {
	return time.Until(e.ExpiresAt)
}

// Tracker maps secret paths to their expiration entries.
type Tracker struct {
	mu      sync.RWMutex
	entries map[string]Entry
}

// New returns an initialised Tracker.
func New() *Tracker {
	return &Tracker{entries: make(map[string]Entry)}
}

// Set records value for path with a TTL measured from now.
func (t *Tracker) Set(path string, value map[string]any, ttl time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries[path] = Entry{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
	}
}

// Get returns the entry for path. It returns ErrExpired if the entry
// exists but has passed its deadline, and a false ok if not found.
func (t *Tracker) Get(path string) (Entry, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	e, ok := t.entries[path]
	if !ok {
		return Entry{}, nil
	}
	if e.Expired() {
		return e, ErrExpired
	}
	return e, nil
}

// Delete removes the entry for path.
func (t *Tracker) Delete(path string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.entries, path)
}

// Purge removes all entries that have expired by now.
func (t *Tracker) Purge(_ context.Context) int {
	t.mu.Lock()
	defer t.mu.Unlock()
	n := 0
	for k, e := range t.entries {
		if e.Expired() {
			delete(t.entries, k)
			n++
		}
	}
	return n
}
