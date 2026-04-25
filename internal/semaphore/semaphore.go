package semaphore

import (
	"context"
	"errors"
	"time"
)

// ErrTimeout is returned when Acquire times out waiting for a slot.
var ErrTimeout = errors.New("semaphore: acquire timed out")

// Semaphore is a counting semaphore that limits concurrent access.
type Semaphore struct {
	slots chan struct{}
}

// New creates a Semaphore with the given capacity.
// It panics if capacity is less than 1.
func New(capacity int) *Semaphore {
	if capacity < 1 {
		panic("semaphore: capacity must be at least 1")
	}
	return &Semaphore{slots: make(chan struct{}, capacity)}
}

// Acquire blocks until a slot is available or ctx is cancelled.
// Returns an error if the context is cancelled.
func (s *Semaphore) Acquire(ctx context.Context) error {
	select {
	case s.slots <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// AcquireTimeout blocks until a slot is available or the timeout elapses.
// Returns ErrTimeout if the deadline is exceeded.
func (s *Semaphore) AcquireTimeout(d time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()
	if err := s.Acquire(ctx); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return ErrTimeout
		}
		return err
	}
	return nil
}

// Release frees one slot. It panics if called more times than Acquire.
func (s *Semaphore) Release() {
	select {
	case <-s.slots:
	default:
		panic("semaphore: Release called without matching Acquire")
	}
}

// Available returns the number of free slots.
func (s *Semaphore) Available() int {
	return cap(s.slots) - len(s.slots)
}

// Capacity returns the total capacity of the semaphore.
func (s *Semaphore) Capacity() int {
	return cap(s.slots)
}
