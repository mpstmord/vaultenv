package observe

import (
	"context"
	"sync/atomic"
)

// Counters holds cumulative counts updated by MetricsObserver.
type Counters struct {
	Total   atomic.Int64
	Errors  atomic.Int64
	Cached  atomic.Int64
}

// MetricsObserver returns an Observer that increments the provided
// Counters on every fetch event. The Counters pointer must not be nil.
func MetricsObserver(c *Counters) Observer {
	if c == nil {
		panic("observe: Counters must not be nil")
	}
	return func(_ context.Context, e Event) {
		c.Total.Add(1)
		if e.Err != nil {
			c.Errors.Add(1)
		}
		if e.Cached {
			c.Cached.Add(1)
		}
	}
}
