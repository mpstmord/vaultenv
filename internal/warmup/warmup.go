// Package warmup provides pre-flight secret fetching to populate a store
// before the child process starts. It fetches all configured secret paths
// concurrently and returns an aggregated error if any fetch fails.
package warmup

import (
	"context"
	"fmt"
	"sync"
)

// Fetcher retrieves secret data for a given path and field.
type Fetcher interface {
	GetSecretData(ctx context.Context, path string) (map[string]interface{}, error)
}

// Writer stores a key/value pair into the environment store.
type Writer interface {
	Set(key, value string)
}

// Mapping describes a single secret field to environment variable binding.
type Mapping struct {
	Path  string
	Field string
	EnvKey string
}

// Runner executes the warmup phase by resolving all mappings.
type Runner struct {
	fetcher Fetcher
	writer  Writer
}

// New creates a new warmup Runner.
func New(f Fetcher, w Writer) *Runner {
	if f == nil {
		panic("warmup: fetcher must not be nil")
	}
	if w == nil {
		panic("warmup: writer must not be nil")
	}
	return &Runner{fetcher: f, writer: w}
}

// Run fetches all mappings concurrently and writes resolved values to the writer.
// It returns an error if any mapping fails to resolve.
func (r *Runner) Run(ctx context.Context, mappings []Mapping) error {
	type result struct {
		key string
		val string
		err error
	}

	results := make([]result, len(mappings))
	var wg sync.WaitGroup
	wg.Add(len(mappings))

	for i, m := range mappings {
		i, m := i, m
		go func() {
			defer wg.Done()
			data, err := r.fetcher.GetSecretData(ctx, m.Path)
			if err != nil {
				results[i] = result{err: fmt.Errorf("warmup: fetch %q: %w", m.Path, err)}
				return
			}
			v, ok := data[m.Field]
			if !ok {
				results[i] = result{err: fmt.Errorf("warmup: field %q not found in %q", m.Field, m.Path)}
				return
			}
			s, ok := v.(string)
			if !ok {
				results[i] = result{err: fmt.Errorf("warmup: field %q in %q is not a string", m.Field, m.Path)}
				return
			}
			results[i] = result{key: m.EnvKey, val: s}
		}()
	}
	wg.Wait()

	var errs []error
	for _, res := range results {
		if res.err != nil {
			errs = append(errs, res.err)
			continue
		}
		r.writer.Set(res.key, res.val)
	}
	if len(errs) > 0 {
		return joinErrors(errs)
	}
	return nil
}

func joinErrors(errs []error) error {
	msg := errs[0].Error()
	for _, e := range errs[1:] {
		msg += "; " + e.Error()
	}
	return fmt.Errorf("%s", msg)
}
