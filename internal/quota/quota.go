// Package quota enforces per-path secret fetch limits within a time window.
package quota

import (
	"fmt"
	"sync"
	"time"
)

// ErrQuotaExceeded is returned when a path has exceeded its allowed requests.
type ErrQuotaExceeded struct {
	Path  string
	Limit int
}

func (e *ErrQuotaExceeded) Error() string {
	return fmt.Sprintf("quota exceeded for path %q: limit %d", e.Path, e.Limit)
}

// entry tracks request count and window start for a single path.
type entry struct {
	count     int
	windowEnd time.Time
}

// Limiter enforces a maximum number of requests per path per window duration.
type Limiter struct {
	mu       sync.Mutex
	entries  map[string]*entry
	limit    int
	window   time.Duration
	nowFunc  func() time.Time
}

// New creates a Limiter allowing at most limit requests per path per window.
func New(limit int, window time.Duration) (*Limiter, error) {
	if limit <= 0 {
		return nil, fmt.Errorf("quota: limit must be positive, got %d", limit)
	}
	if window <= 0 {
		return nil, fmt.Errorf("quota: window must be positive, got %s", window)
	}
	return &Limiter{
		entries: make(map[string]*entry),
		limit:   limit,
		window:  window,
		nowFunc: time.Now,
	}, nil
}

// Allow checks whether a request for path is within quota.
// It returns ErrQuotaExceeded if the limit has been reached for the current window.
func (l *Limiter) Allow(path string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.nowFunc()
	e, ok := l.entries[path]
	if !ok || now.After(e.windowEnd) {
		l.entries[path] = &entry{count: 1, windowEnd: now.Add(l.window)}
		return nil
	}
	if e.count >= l.limit {
		return &ErrQuotaExceeded{Path: path, Limit: l.limit}
	}
	e.count++
	return nil
}

// Reset clears quota state for a specific path.
func (l *Limiter) Reset(path string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.entries, path)
}

// Stats returns the current request count and window end for a path.
// Returns zero values if no requests have been recorded.
func (l *Limiter) Stats(path string) (count int, windowEnd time.Time) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if e, ok := l.entries[path]; ok {
		return e.count, e.windowEnd
	}
	return 0, time.Time{}
}
