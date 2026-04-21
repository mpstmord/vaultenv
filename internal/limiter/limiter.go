// Package limiter provides a concurrency limiter that caps the number of
// simultaneous in-flight secret fetch operations.
package limiter

import (
	"context"
	"errors"
	"fmt"
)

// ErrLimitExceeded is returned when the concurrency limit is reached and the
// caller's context expires before a slot becomes available.
var ErrLimitExceeded = errors.New("limiter: concurrency limit exceeded")

// Limiter caps the number of goroutines that may execute concurrently.
type Limiter struct {
	sem chan struct{}
}

// New creates a Limiter that allows at most n concurrent operations.
// n must be >= 1.
func New(n int) (*Limiter, error) {
	if n < 1 {
		return nil, fmt.Errorf("limiter: concurrency limit must be >= 1, got %d", n)
	}
	return &Limiter{sem: make(chan struct{}, n)}, nil
}

// Acquire blocks until a concurrency slot is available or ctx is done.
// Returns ErrLimitExceeded if the context expires while waiting.
func (l *Limiter) Acquire(ctx context.Context) error {
	select {
	case l.sem <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ErrLimitExceeded
	}
}

// Release frees a previously acquired concurrency slot.
// Callers must call Release exactly once for each successful Acquire.
func (l *Limiter) Release() {
	<-l.sem
}

// Capacity returns the maximum number of concurrent operations.
func (l *Limiter) Capacity() int {
	return cap(l.sem)
}

// InFlight returns the number of operations currently holding a slot.
func (l *Limiter) InFlight() int {
	return len(l.sem)
}
