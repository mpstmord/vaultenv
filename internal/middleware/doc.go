// Package middleware provides composable Fetcher wrappers for the vaultenv
// secret retrieval pipeline.
//
// Each middleware implements a cross-cutting concern — logging, rate-limiting,
// retry, etc. — and can be composed using Chain to build a layered Fetcher
// without modifying the underlying Vault client.
//
// Usage:
//
//	f := middleware.Chain(
//		base,
//		middleware.NewLoggingMiddleware(logger),
//		middleware.NewRateLimitMiddleware(rl),
//		middleware.NewRetryMiddleware(retry.DefaultPolicy),
//	)
package middleware
