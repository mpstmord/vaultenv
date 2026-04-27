// Package hedge provides a hedged-request wrapper around any secret Fetcher.
//
// Hedging is a latency-reduction technique: a second identical request is
// issued after a short configurable delay if the first has not yet returned.
// Whichever response arrives first is forwarded to the caller and the slower
// goroutine is abandoned (its result is discarded once the context is done).
//
// Usage:
//
//	h := hedge.New(vaultClient, 50*time.Millisecond)
//	data, err := h.GetSecretData(ctx, "secret/data/myapp")
package hedge
