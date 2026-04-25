// Package fanout provides concurrent secret fetching that dispatches
// a single logical request to multiple upstream fetchers and merges results.
package fanout

import (
	"context"
	"fmt"
	"sync"
)

// Fetcher is the interface for retrieving secret data.
type Fetcher interface {
	GetSecretData(ctx context.Context, path string) (map[string]interface{}, error)
}

// Fanout dispatches GetSecretData to multiple fetchers concurrently and
// merges their results. Later fetchers override keys from earlier ones.
type Fanout struct {
	fetchers []Fetcher
}

// New returns a Fanout that queries each fetcher in parallel.
// It panics if fetchers is empty.
func New(fetchers ...Fetcher) *Fanout {
	if len(fetchers) == 0 {
		panic("fanout: at least one fetcher is required")
	}
	return &Fanout{fetchers: fetchers}
}

type result struct {
	data map[string]interface{}
	err  error
	idx  int
}

// GetSecretData queries all upstream fetchers concurrently for path and
// merges their responses. If any fetcher returns an error the first error
// encountered (by original order) is returned.
func (f *Fanout) GetSecretData(ctx context.Context, path string) (map[string]interface{}, error) {
	results := make([]result, len(f.fetchers))
	var wg sync.WaitGroup

	for i, fetcher := range f.fetchers {
		wg.Add(1)
		go func(idx int, ft Fetcher) {
			defer wg.Done()
			data, err := ft.GetSecretData(ctx, path)
			results[idx] = result{data: data, err: err, idx: idx}
		}(i, fetcher)
	}

	wg.Wait()

	merged := make(map[string]interface{})
	for _, r := range results {
		if r.err != nil {
			return nil, fmt.Errorf("fanout: fetcher %d: %w", r.idx, r.err)
		}
		for k, v := range r.data {
			merged[k] = v
		}
	}
	return merged, nil
}
