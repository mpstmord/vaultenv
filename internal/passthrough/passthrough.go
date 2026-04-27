// Package passthrough provides a middleware that selectively bypasses
// secret fetching for keys that are already present in the environment.
package passthrough

import (
	"context"
	"os"
	"strings"
)

// Fetcher is the interface for retrieving secret data from a backend.
type Fetcher interface {
	GetSecretData(ctx context.Context, path string) (map[string]interface{}, error)
}

// Passthrough wraps a Fetcher and short-circuits lookups for paths whose
// environment variable override is already set.
type Passthrough struct {
	upstream Fetcher
	prefix   string
}

// New returns a Passthrough that delegates to upstream when no environment
// override is found. prefix is prepended to the uppercased secret path when
// constructing the env var name (e.g. "VAULTENV_").
func New(upstream Fetcher, prefix string) *Passthrough {
	if upstream == nil {
		panic("passthrough: upstream fetcher must not be nil")
	}
	return &Passthrough{upstream: upstream, prefix: prefix}
}

// GetSecretData returns a synthetic single-field map from the environment when
// an override variable is present, otherwise delegates to the upstream fetcher.
func (p *Passthrough) GetSecretData(ctx context.Context, path string) (map[string]interface{}, error) {
	envKey := p.envKey(path)
	if val, ok := os.LookupEnv(envKey); ok {
		return map[string]interface{}{"value": val}, nil
	}
	return p.upstream.GetSecretData(ctx, path)
}

// envKey converts a secret path to an environment variable name.
// Slashes and hyphens are replaced with underscores and the result is
// uppercased before the prefix is applied.
func (p *Passthrough) envKey(path string) string {
	r := strings.NewReplacer("/", "_", "-", "_")
	norm := strings.ToUpper(r.Replace(path))
	if p.prefix == "" {
		return norm
	}
	return p.prefix + norm
}
