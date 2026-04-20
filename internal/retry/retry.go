// Package retry provides configurable retry logic with exponential backoff
// for transient errors encountered when communicating with Vault.
package retry

import (
	"context"
	"errors"
	"math"
	"time"
)

// Policy defines the retry behaviour.
type Policy struct {
	// MaxAttempts is the total number of attempts (including the first call).
	MaxAttempts int
	// BaseDelay is the wait time before the second attempt.
	BaseDelay time.Duration
	// MaxDelay caps the exponential back-off.
	MaxDelay time.Duration
}

// DefaultPolicy returns a sensible default retry policy.
func DefaultPolicy() Policy {
	return Policy{
		MaxAttempts: 4,
		BaseDelay:   250 * time.Millisecond,
		MaxDelay:    10 * time.Second,
	}
}

// ErrMaxAttempts is returned when all attempts are exhausted.
var ErrMaxAttempts = errors.New("retry: max attempts reached")

// Do calls fn repeatedly according to p until fn returns nil, the context is
// cancelled, or the maximum number of attempts is reached.  The last non-nil
// error returned by fn is wrapped with ErrMaxAttempts and returned.
func Do(ctx context.Context, p Policy, fn func() error) error {
	if p.MaxAttempts <= 0 {
		p.MaxAttempts = 1
	}

	var lastErr error
	for attempt := 0; attempt < p.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		lastErr = fn()
		if lastErr == nil {
			return nil
		}

		if attempt < p.MaxAttempts-1 {
			delay := backoff(p.BaseDelay, p.MaxDelay, attempt)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
		}
	}

	return errors.Join(ErrMaxAttempts, lastErr)
}

// backoff computes exponential back-off: BaseDelay * 2^attempt, capped at MaxDelay.
func backoff(base, max time.Duration, attempt int) time.Duration {
	exp := time.Duration(math.Pow(2, float64(attempt))) * base
	if exp > max || exp <= 0 {
		return max
	}
	return exp
}
