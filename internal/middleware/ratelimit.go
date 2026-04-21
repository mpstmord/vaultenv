package middleware

import (
	"context"
	"fmt"
	"time"
)

// RateLimitError is returned when a request is rejected by the rate limiter.
type RateLimitError struct {
	Path string
	RetryAfter time.Duration
}

func (e *RateLimitError) Error() string {
	if e.RetryAfter > 0 {
		return fmt.Sprintf("rate limit exceeded for path %q: retry after %s", e.Path, e.RetryAfter)
	}
	return fmt.Sprintf("rate limit exceeded for path %q", e.Path)
}

// Limiter is the interface expected by NewRateLimitMiddleware.
// It mirrors the Allow method of internal/ratelimit.Limiter so that
// callers can pass either implementation without importing the package
// directly.
type Limiter interface {
	Allow() bool
}

// NewRateLimitMiddleware returns a Fetcher middleware that enforces the
// provided Limiter before forwarding the request to the upstream Fetcher.
//
// When the limiter rejects a call, a *RateLimitError is returned immediately
// and the upstream is never contacted.  An optional retryAfter duration can
// be supplied so that callers can back off appropriately; pass zero to omit
// it from the error.
func NewRateLimitMiddleware(limiter Limiter, retryAfter time.Duration) func(Fetcher) Fetcher {
	if limiter == nil {
		panic("middleware: NewRateLimitMiddleware requires a non-nil Limiter")
	}
	return func(next Fetcher) Fetcher {
		return FetcherFunc(func(ctx context.Context, path string) (map[string]interface{}, error) {
			if !limiter.Allow() {
				return nil, &RateLimitError{Path: path, RetryAfter: retryAfter}
			}
			return next.GetSecretData(ctx, path)
		})
	}
}
