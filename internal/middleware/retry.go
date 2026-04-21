package middleware

import (
	"context"

	"github.com/your-org/vaultenv/internal/retry"
)

// RetryMiddleware wraps a Fetcher with automatic retry logic using the
// provided retry policy. Transient errors are retried according to the
// policy; the first successful result is returned immediately.
type RetryMiddleware struct {
	policy retry.Policy
}

// NewRetryMiddleware returns a new RetryMiddleware with the given policy.
// If policy is the zero value, retry.DefaultPolicy is used.
func NewRetryMiddleware(policy retry.Policy) *RetryMiddleware {
	if policy.Attempts == 0 {
		policy = retry.DefaultPolicy
	}
	return &RetryMiddleware{policy: policy}
}

// Wrap returns a Fetcher that retries the upstream on error.
func (m *RetryMiddleware) Wrap(next Fetcher) Fetcher {
	return FetcherFunc(func(ctx context.Context, path string) (map[string]interface{}, error) {
		var result map[string]interface{}
		err := retry.Do(ctx, m.policy, func() error {
			var fetchErr error
			result, fetchErr = next.GetSecretData(ctx, path)
			return fetchErr
		})
		if err != nil {
			return nil, err
		}
		return result, nil
	})
}
