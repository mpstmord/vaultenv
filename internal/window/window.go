// Package window provides a sliding time-window counter for tracking
// event rates over a rolling duration.
package window

import (
	"fmt"
	"sync"
	"time"
)

// Counter is a thread-safe sliding-window event counter.
type Counter struct {
	mu      sync.Mutex
	window  time.Duration
	buckets int
	counts  []int64
	times   []time.Time
	cursor  int
}

// New creates a Counter that tracks events over the given window, divided
// into the specified number of buckets. buckets must be >= 1.
func New(window time.Duration, buckets int) (*Counter, error) {
	if window <= 0 {
		return nil, fmt.Errorf("window must be positive, got %s", window)
	}
	if buckets < 1 {
		return nil, fmt.Errorf("buckets must be >= 1, got %d", buckets)
	}
	return &Counter{
		window:  window,
		buckets: buckets,
		counts:  make([]int64, buckets),
		times:   make([]time.Time, buckets),
	}, nil
}

// Add records n events at the current time.
func (c *Counter) Add(n int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	c.advance(now)
	c.counts[c.cursor] += n
}

// Total returns the sum of all events recorded within the sliding window.
func (c *Counter) Total() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	c.advance(now)
	cutoff := now.Add(-c.window)
	var total int64
	for i := 0; i < c.buckets; i++ {
		if c.times[i].After(cutoff) {
			total += c.counts[i]
		}
	}
	return total
}

// Reset clears all recorded events.
func (c *Counter) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for i := range c.counts {
		c.counts[i] = 0
		c.times[i] = time.Time{}
	}
	c.cursor = 0
}

// advance moves the cursor forward when the current bucket has aged past
// the per-bucket duration, clearing the new bucket before use.
func (c *Counter) advance(now time.Time) {
	if c.times[c.cursor].IsZero() {
		c.times[c.cursor] = now
		return
	}
	bucketDur := c.window / time.Duration(c.buckets)
	if now.Sub(c.times[c.cursor]) >= bucketDur {
		next := (c.cursor + 1) % c.buckets
		c.cursor = next
		c.counts[c.cursor] = 0
		c.times[c.cursor] = now
	}
}
