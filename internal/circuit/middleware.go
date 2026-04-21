package circuit

import (
	"context"
	"fmt"
)

// SecretFetcher fetches a secret value by path and field.
type SecretFetcher interface {
	GetSecretData(ctx context.Context, path string) (map[string]interface{}, error)
}

// GuardedFetcher wraps a SecretFetcher with circuit-breaker protection.
type GuardedFetcher struct {
	upstream SecretFetcher
	breaker  *Breaker
}

// NewGuardedFetcher returns a GuardedFetcher that protects upstream with b.
func NewGuardedFetcher(upstream SecretFetcher, b *Breaker) *GuardedFetcher {
	return &GuardedFetcher{upstream: upstream, breaker: b}
}

// GetSecretData calls the upstream fetcher if the circuit allows it.
// On success the breaker is reset; on failure the breaker records the error.
func (g *GuardedFetcher) GetSecretData(ctx context.Context, path string) (map[string]interface{}, error) {
	if err := g.breaker.Allow(); err != nil {
		return nil, fmt.Errorf("circuit: %w", err)
	}
	data, err := g.upstream.GetSecretData(ctx, path)
	if err != nil {
		g.breaker.RecordFailure()
		return nil, err
	}
	g.breaker.RecordSuccess()
	return data, nil
}
