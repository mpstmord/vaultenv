// Package middleware provides a composable fetcher middleware chain
// for wrapping secret retrieval with cross-cutting concerns.
package middleware

import "context"

// Fetcher retrieves secret data for a given path.
type Fetcher interface {
	GetSecretData(ctx context.Context, path string) (map[string]interface{}, error)
}

// Middleware wraps a Fetcher with additional behaviour.
type Middleware func(Fetcher) Fetcher

// Chain applies a list of middlewares to a base Fetcher, outermost first.
// Chain(base, m1, m2) produces m1(m2(base)).
func Chain(base Fetcher, middlewares ...Middleware) Fetcher {
	for i := len(middlewares) - 1; i >= 0; i-- {
		base = middlewares[i](base)
	}
	return base
}
