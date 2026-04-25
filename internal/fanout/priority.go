package fanout

import (
	"context"
	"fmt"
)

// PriorityFetcher queries fetchers in order and returns the result of the
// first one that succeeds. This differs from Fanout which queries all
// fetchers concurrently and merges results.
type PriorityFetcher struct {
	fetchers []Fetcher
}

// NewPriority returns a PriorityFetcher that tries each fetcher in order.
// It panics if no fetchers are provided.
func NewPriority(fetchers ...Fetcher) *PriorityFetcher {
	if len(fetchers) == 0 {
		panic("fanout: priority requires at least one fetcher")
	}
	return &PriorityFetcher{fetchers: fetchers}
}

// GetSecretData tries each fetcher in order, returning the first successful
// response. If all fetchers fail the last error is returned.
func (p *PriorityFetcher) GetSecretData(ctx context.Context, path string) (map[string]interface{}, error) {
	var last error
	for i, f := range p.fetchers {
		data, err := f.GetSecretData(ctx, path)
		if err == nil {
			return data, nil
		}
		last = fmt.Errorf("fanout: priority fetcher %d: %w", i, err)
		if ctx.Err() != nil {
			return nil, last
		}
	}
	return nil, last
}
