package metrics

import (
	"sync"
	"sync/atomic"
	"time"
)

// Counter tracks a monotonically increasing value.
type Counter struct {
	value uint64
}

func (c *Counter) Inc() { atomic.AddUint64(&c.value, 1) }
func (c *Counter) Get() uint64 { return atomic.LoadUint64(&c.value) }

// Collector holds runtime metrics for vaultenv.
type Collector struct {
	mu sync.RWMutex
	start time.Time

	SecretsResolved Counter
	CacheHits       Counter
	CacheMisses     Counter
	Renewals        Counter
	Errors          Counter
}

// New returns a new Collector with the start time set to now.
func New() *Collector {
	return &Collector{start: time.Now()}
}

// Uptime returns the duration since the Collector was created.
func (c *Collector) Uptime() time.Duration {
	return time.Since(c.start)
}

// Snapshot returns a point-in-time copy of all metric values.
func (c *Collector) Snapshot() map[string]uint64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return map[string]uint64{
		"secrets_resolved": c.SecretsResolved.Get(),
		"cache_hits":       c.CacheHits.Get(),
		"cache_misses":     c.CacheMisses.Get(),
		"renewals":         c.Renewals.Get(),
		"errors":           c.Errors.Get(),
	}
}
