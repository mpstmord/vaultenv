// Package clock provides a testable abstraction over time.
package clock

import "time"

// Clock is an interface for time-related operations, allowing tests to
// substitute a fake implementation without patching global time functions.
type Clock interface {
	// Now returns the current local time.
	Now() time.Time
	// Since returns the elapsed time since t.
	Since(t time.Time) time.Duration
	// After returns a channel that fires after the given duration.
	After(d time.Duration) <-chan time.Time
}

// Real is a Clock that delegates to the standard library.
type Real struct{}

// New returns a Real clock backed by the standard library.
func New() Clock {
	return Real{}
}

// Now returns time.Now().
func (Real) Now() time.Time { return time.Now() }

// Since returns time.Since(t).
func (Real) Since(t time.Time) time.Duration { return time.Since(t) }

// After returns time.After(d).
func (Real) After(d time.Duration) <-chan time.Time { return time.After(d) }

// Fake is a manually-controlled Clock for use in tests.
type Fake struct {
	current time.Time
}

// NewFake returns a Fake clock set to the given start time.
func NewFake(start time.Time) *Fake {
	return &Fake{current: start}
}

// Now returns the fake clock's current time.
func (f *Fake) Now() time.Time { return f.current }

// Since returns the duration elapsed since t relative to the fake clock.
func (f *Fake) Since(t time.Time) time.Duration { return f.current.Sub(t) }

// After returns a channel that fires immediately with the fake current time
// advanced by d, and also advances the clock by d.
func (f *Fake) After(d time.Duration) <-chan time.Time {
	f.current = f.current.Add(d)
	ch := make(chan time.Time, 1)
	ch <- f.current
	return ch
}

// Advance moves the fake clock forward by the given duration.
func (f *Fake) Advance(d time.Duration) {
	f.current = f.current.Add(d)
}

// Set sets the fake clock to an absolute time.
func (f *Fake) Set(t time.Time) {
	f.current = t
}
