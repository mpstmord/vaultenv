// Package jitter provides utilities for adding randomised jitter to
// durations, which helps spread out retry and renewal bursts.
package jitter

import (
	"math/rand"
	"time"
)

// Source is a function that returns a pseudo-random float64 in [0.0, 1.0).
// It can be replaced in tests to produce deterministic results.
type Source func() float64

// Jitter adds a random fraction of base to base and returns the result.
// The fraction is drawn from src and clamped to [0, factor], where factor
// must be in (0, 1]. If factor is outside that range it is clamped to 1.
//
//	result = base + base * factor * rand()
func Jitter(base time.Duration, factor float64, src Source) time.Duration {
	if factor <= 0 {
		return base
	}
	if factor > 1 {
		factor = 1
	}
	delta := float64(base) * factor * src()
	return base + time.Duration(delta)
}

// Default returns a Jitter call using the global math/rand source.
func Default(base time.Duration, factor float64) time.Duration {
	return Jitter(base, factor, rand.Float64)
}

// Full returns a duration uniformly distributed in [0, 2*base], which gives
// a mean equal to base ("full jitter" strategy).
func Full(base time.Duration, src Source) time.Duration {
	if base <= 0 {
		return 0
	}
	return time.Duration(float64(2*base) * src())
}

// Equal returns a duration uniformly distributed in [base/2, base]
// ("equal jitter" strategy).
func Equal(base time.Duration, src Source) time.Duration {
	if base <= 0 {
		return 0
	}
	half := base / 2
	return half + time.Duration(float64(half)*src())
}
