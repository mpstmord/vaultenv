package middleware

import "context"

// FetcherFunc is a function type that implements Fetcher.
// It allows inline construction of middleware without defining a named type.
type FetcherFunc func(ctx context.Context, path string) (map[string]interface{}, error)

// GetSecretData calls the underlying function.
func (f FetcherFunc) GetSecretData(ctx context.Context, path string) (map[string]interface{}, error) {
	return f(ctx, path)
}
