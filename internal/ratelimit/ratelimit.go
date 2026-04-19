// Package ratelimit provides a simple token-bucket rate limiter for
// controlling the frequency of Vault API requests.
package ratelimit

import (
	"fmt"
	"sync"
	"time"
)

// Limiter controls the rate of operations using a token-bucket algorithm.
type Limiter struct {
	mu       sync.Mutex
	tokens   float64
	max      float64
	rate     float64 // tokens per second
	lastTick time.Time
	clock    func() time.Time
}

// New creates a Limiter that allows up to burst operations and refills at
// rate operations per second.
func New(ratePerSec float64, burst int) (*Limiter, error) {
	if ratePerSec <= 0 {
		return nil, fmt.Errorf("ratelimit: rate must be positive, got %v", ratePerSec)
	}
	if burst <= 0 {
		return nil, fmt.Errorf("ratelimit: burst must be positive, got %d", burst)
	}
	return &Limiter{
		tokens:   float64(burst),
		max:      float64(burst),
		rate:     ratePerSec,
		lastTick: time.Now(),
		clock:    time.Now,
	}, nil
}

// Allow reports whether an operation may proceed. It is safe for concurrent use.
func (l *Limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.clock()
	elapsed := now.Sub(l.lastTick).Seconds()
	l.lastTick = now

	l.tokens += elapsed * l.rate
	if l.tokens > l.max {
		l.tokens = l.max
	}

	if l.tokens < 1 {
		return false
	}
	l.tokens--
	return true
}

// Wait blocks until an operation is allowed or the context deadline is reached.
func (l *Limiter) Wait() error {
	for {
		if l.Allow() {
			return nil
		}
		time.Sleep(time.Duration(float64(time.Second) / l.rate))
	}
}
