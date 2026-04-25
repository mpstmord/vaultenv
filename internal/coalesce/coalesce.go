// Package coalesce provides a fetcher that returns the first non-error
// result from an ordered list of secret fetchers, merging their data.
package coalesce

import (
	"context"
	"errors"
	"fmt"
)

// Fetcher is the interface for retrieving secret data.
type Fetcher interface {
	GetSecretData(ctx context.Context, path string) (map[string]interface{}, error)
}

// Coalescer wraps multiple fetchers and returns the first successful result.
type Coalescer struct {
	fetchers []Fetcher
}

// New creates a Coalescer from the provided fetchers.
// Panics if no fetchers are supplied.
func New(fetchers ...Fetcher) *Coalescer {
	if len(fetchers) == 0 {
		panic("coalesce: at least one fetcher is required")
	}
	return &Coalescer{fetchers: fetchers}
}

// GetSecretData queries each fetcher in order and returns the data from the
// first one that succeeds. If all fetchers fail, a combined error is returned.
func (c *Coalescer) GetSecretData(ctx context.Context, path string) (map[string]interface{}, error) {
	var errs []error
	for i, f := range c.fetchers {
		data, err := f.GetSecretData(ctx, path)
		if err == nil {
			return data, nil
		}
		errs = append(errs, fmt.Errorf("fetcher[%d]: %w", i, err))
	}
	return nil, errors.Join(errs...)
}
