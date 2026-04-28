// Package resolve provides a multi-step secret resolution pipeline
// that walks a prioritised list of fetchers and returns the first
// successful result, recording which source satisfied the request.
package resolve

import (
	"context"
	"errors"
	"fmt"
)

// Fetcher is the common interface for any secret backend.
type Fetcher interface {
	GetSecretData(ctx context.Context, path string) (map[string]interface{}, error)
}

// Source pairs a human-readable name with a Fetcher.
type Source struct {
	Name    string
	Fetcher Fetcher
}

// Result holds the resolved secret data together with the name of the
// source that provided it.
type Result struct {
	Data   map[string]interface{}
	Source string
}

// Resolver walks an ordered list of Sources and returns the first
// successful response.
type Resolver struct {
	sources []Source
}

// New returns a Resolver that will consult sources in the order given.
// At least one source must be supplied.
func New(sources ...Source) (*Resolver, error) {
	if len(sources) == 0 {
		return nil, errors.New("resolve: at least one source is required")
	}
	for i, s := range sources {
		if s.Name == "" {
			return nil, fmt.Errorf("resolve: source at index %d has an empty name", i)
		}
		if s.Fetcher == nil {
			return nil, fmt.Errorf("resolve: source %q has a nil fetcher", s.Name)
		}
	}
	return &Resolver{sources: sources}, nil
}

// Resolve queries each source in priority order and returns the first
// successful Result. If all sources fail the last error is returned.
func (r *Resolver) Resolve(ctx context.Context, path string) (Result, error) {
	var lastErr error
	for _, src := range r.sources {
		data, err := src.Fetcher.GetSecretData(ctx, path)
		if err != nil {
			lastErr = err
			continue
		}
		return Result{Data: data, Source: src.Name}, nil
	}
	return Result{}, fmt.Errorf("resolve: all sources exhausted for %q: %w", path, lastErr)
}
