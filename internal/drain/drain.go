// Package drain provides graceful shutdown coordination for vaultenv,
// ensuring in-flight secret fetches complete before the process exits.
package drain

import (
	"context"
	"sync"
	"time"
)

// Drainer tracks active operations and waits for them to finish
// within a configurable deadline before allowing shutdown.
type Drainer struct {
	mu      sync.Mutex
	wg      sync.WaitGroup
	closed  bool
	timeout time.Duration
}

// New creates a Drainer with the given shutdown timeout.
// If timeout is zero or negative, DefaultTimeout is used.
const DefaultTimeout = 10 * time.Second

func New(timeout time.Duration) *Drainer {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	return &Drainer{timeout: timeout}
}

// Acquire marks the start of an operation. It returns false if the
// Drainer is already draining (i.e. shutdown has been initiated).
func (d *Drainer) Acquire() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.closed {
		return false
	}
	d.wg.Add(1)
	return true
}

// Release marks the completion of an operation previously acquired.
func (d *Drainer) Release() {
	d.wg.Done()
}

// Drain initiates shutdown: no new Acquire calls will succeed, and
// Drain blocks until all in-flight operations finish or ctx is done.
// It returns context.DeadlineExceeded if the timeout is reached.
func (d *Drainer) Drain(ctx context.Context) error {
	d.mu.Lock()
	d.closed = true
	d.mu.Unlock()

	deadline, cancel := context.WithTimeout(ctx, d.timeout)
	defer cancel()

	done := make(chan struct{})
	go func() {
		d.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-deadline.Done():
		return deadline.Err()
	}
}

// Closed reports whether Drain has been called.
func (d *Drainer) Closed() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.closed
}
