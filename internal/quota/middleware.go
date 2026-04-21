package quota

import "context"

// SecretFetcher is the interface satisfied by vault clients that retrieve
// secret data.
type SecretFetcher interface {
	GetSecretData(ctx context.Context, path string) (map[string]interface{}, error)
}

// GuardedFetcher wraps a SecretFetcher and enforces quota limits before
// delegating to the underlying implementation.
type GuardedFetcher struct {
	fetcher SecretFetcher
	limiter *Limiter
}

// NewGuardedFetcher creates a GuardedFetcher that applies the given Limiter
// before every call to GetSecretData.
func NewGuardedFetcher(f SecretFetcher, l *Limiter) *GuardedFetcher {
	return &GuardedFetcher{fetcher: f, limiter: l}
}

// GetSecretData checks the quota for path and, if allowed, delegates to the
// underlying SecretFetcher. It returns ErrQuotaExceeded without calling the
// upstream client when the limit has been reached.
func (g *GuardedFetcher) GetSecretData(ctx context.Context, path string) (map[string]interface{}, error) {
	if err := g.limiter.Allow(path); err != nil {
		return nil, err
	}
	return g.fetcher.GetSecretData(ctx, path)
}
