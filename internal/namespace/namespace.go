// Package namespace provides utilities for scoping secret paths
// under a configurable prefix, allowing multi-tenant or environment-aware
// secret isolation within a shared Vault instance.
package namespace

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

// Fetcher is the interface for retrieving secret data.
type Fetcher interface {
	GetSecretData(ctx context.Context, path string) (map[string]interface{}, error)
}

// Namespace scopes all secret paths under a fixed prefix.
type Namespace struct {
	prefix  string
	upstream Fetcher
}

// New returns a Namespace that prepends prefix to every path passed to
// GetSecretData. prefix must be non-empty and must not contain leading or
// trailing slashes (they are trimmed automatically).
func New(prefix string, upstream Fetcher) (*Namespace, error) {
	if upstream == nil {
		panic("namespace: upstream fetcher must not be nil")
	}
	prefix = strings.Trim(prefix, "/")
	if prefix == "" {
		return nil, errors.New("namespace: prefix must not be empty")
	}
	return &Namespace{prefix: prefix, upstream: upstream}, nil
}

// Prefix returns the configured namespace prefix.
func (n *Namespace) Prefix() string {
	return n.prefix
}

// GetSecretData prepends the namespace prefix to path and delegates to the
// upstream Fetcher. The path is joined with a single "/" separator.
func (n *Namespace) GetSecretData(ctx context.Context, path string) (map[string]interface{}, error) {
	path = strings.TrimLeft(path, "/")
	scoped := fmt.Sprintf("%s/%s", n.prefix, path)
	return n.upstream.GetSecretData(ctx, scoped)
}
