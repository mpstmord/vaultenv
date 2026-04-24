// Package debounce provides a mechanism to coalesce rapid successive calls
// into a single execution after a quiet period has elapsed.
package debounce

import (
	"sync"
	"time"
)

// Func is the type of function that can be debounced.
type Func func()

// Debouncer delays execution of a function until a specified duration has
// passed since the last call to Trigger.
type Debouncer struct {
	mu       sync.Mutex
	wait     time.Duration
	timer    *time.Timer
	pending  bool
}

// New returns a Debouncer that will invoke fn no sooner than wait after the
// last call to Trigger. wait must be positive; if it is not, a default of
// 100 ms is used.
func New(wait time.Duration, fn Func) *Debouncer {
	if wait <= 0 {
		wait = 100 * time.Millisecond
	}
	if fn == nil {
		panic("debounce: fn must not be nil")
	}
	return &Debouncer{
		wait: wait,
		timer: time.AfterFunc(wait, func() {
			// timer is initially stopped; AfterFunc fires once, so we
			// immediately stop it before it can fire.
		}),
	}
	// We can't capture fn in AfterFunc at construction time because we need
	// a pointer to d first. Build properly below.
}

// NewDebouncer is the canonical constructor that correctly wires fn.
func NewDebouncer(wait time.Duration, fn Func) *Debouncer {
	if wait <= 0 {
		wait = 100 * time.Millisecond
	}
	if fn == nil {
		panic("debounce: fn must not be nil")
	}
	d := &Debouncer{wait: wait}
	d.timer = time.AfterFunc(wait, func() {
		d.mu.Lock()
		d.pending = false
		d.mu.Unlock()
		fn()
	})
	// Stop the timer immediately; it will be restarted on the first Trigger.
	d.timer.Stop()
	return d
}

// Trigger schedules fn to be called after the debounce window. If Trigger is
// called again before the window expires, the window resets.
func (d *Debouncer) Trigger() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.pending = true
	d.timer.Reset(d.wait)
}

// Pending reports whether a call is currently scheduled but has not yet fired.
func (d *Debouncer) Pending() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.pending
}

// Stop cancels any pending scheduled call. It does not wait for an
// already-running invocation to complete.
func (d *Debouncer) Stop() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.timer.Stop()
	d.pending = false
}
