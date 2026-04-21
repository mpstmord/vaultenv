// Package backoff provides configurable exponential backoff strategies
// for use with retry loops and polling operations.
package backoff

import (
	"math"
	"math/rand"
	"time"
)

// Strategy defines how delays are calculated between attempts.
type Strategy struct {
	// Base is the initial delay duration.
	Base time.Duration
	// Max is the upper bound on the delay duration.
	Max time.Duration
	// Factor is the multiplier applied on each successive attempt.
	Factor float64
	// Jitter adds randomness to avoid thundering-herd problems.
	Jitter bool
}

// Default returns a Strategy with sensible defaults suitable for
// most Vault API retry scenarios.
func Default() Strategy {
	return Strategy{
		Base:   200 * time.Millisecond,
		Max:    30 * time.Second,
		Factor: 2.0,
		Jitter: true,
	}
}

// Delay returns the wait duration for the given attempt number (0-indexed).
// The result is capped at Strategy.Max.
func (s Strategy) Delay(attempt int) time.Duration {
	if attempt < 0 {
		attempt = 0
	}

	base := float64(s.Base)
	delay := base * math.Pow(s.Factor, float64(attempt))

	if s.Jitter {
		// Add up to 20% random jitter.
		delay += rand.Float64() * delay * 0.2 //nolint:gosec
	}

	result := time.Duration(delay)
	if result > s.Max {
		result = s.Max
	}
	return result
}

// Attempts returns a slice of delays for n attempts, useful for inspection
// and testing.
func (s Strategy) Attempts(n int) []time.Duration {
	delays := make([]time.Duration, n)
	for i := 0; i < n; i++ {
		delays[i] = s.Delay(i)
	}
	return delays
}
