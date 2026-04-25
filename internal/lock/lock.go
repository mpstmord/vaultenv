// Package lock provides distributed-style advisory locking for secret
// fetch operations, preventing concurrent duplicate writes to the same key.
package lock

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ErrTimeout is returned when a lock cannot be acquired within the deadline.
var ErrTimeout = fmt.Errorf("lock: timed out waiting to acquire")

// ErrCancelled is returned when the context is cancelled while waiting.
var ErrCancelled = fmt.Errorf("lock: context cancelled while waiting")

const defaultTimeout = 5 * time.Second

// Locker manages per-key mutexes.
type Locker struct {
	mu      sync.Mutex
	locks   map[string]*entry
	timeout time.Duration
}

type entry struct {
	mu      sync.Mutex
	holders int
}

// New creates a Locker with an optional acquisition timeout.
// If timeout is zero, defaultTimeout is used.
func New(timeout time.Duration) *Locker {
	if timeout <= 0 {
		timeout = defaultTimeout
	}
	return &Locker{
		locks:   make(map[string]*entry),
		timeout: timeout,
	}
}

// Lock acquires the advisory lock for key, blocking until it is available,
// the context is cancelled, or the configured timeout elapses.
func (l *Locker) Lock(ctx context.Context, key string) error {
	e := l.getOrCreate(key)

	acquired := make(chan struct{}, 1)
	go func() {
		e.mu.Lock()
		acquired <- struct{}{}
	}()

	timer := time.NewTimer(l.timeout)
	defer timer.Stop()

	select {
	case <-acquired:
		return nil
	case <-ctx.Done():
		return ErrCancelled
	case <-timer.C:
		return ErrTimeout
	}
}

// Unlock releases the advisory lock for key.
func (l *Locker) Unlock(key string) {
	l.mu.Lock()
	e, ok := l.locks[key]
	if !ok {
		l.mu.Unlock()
		return
	}
	e.holders--
	if e.holders == 0 {
		delete(l.locks, key)
	}
	l.mu.Unlock()
	e.mu.Unlock()
}

func (l *Locker) getOrCreate(key string) *entry {
	l.mu.Lock()
	defer l.mu.Unlock()
	e, ok := l.locks[key]
	if !ok {
		e = &entry{}
		l.locks[key] = e
	}
	e.holders++
	return e
}
