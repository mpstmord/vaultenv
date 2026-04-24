// Package batch provides concurrent secret fetching for multiple paths.
package batch

import (
	"context"
	"fmt"
	"sync"
)

// Fetcher retrieves secret data by path and field.
type Fetcher interface {
	GetSecretData(ctx context.Context, path string) (map[string]interface{}, error)
}

// Request describes a single secret lookup.
type Request struct {
	Path  string
	Field string
}

// Result holds the outcome of a single Request.
type Result struct {
	Request Request
	Value   string
	Err     error
}

// Fetcher concurrently resolves a slice of Requests.
type BatchFetcher struct {
	upstream Fetcher
	workers  int
}

// New returns a BatchFetcher with the given upstream and worker count.
// If workers < 1 it defaults to 4.
func New(upstream Fetcher, workers int) *BatchFetcher {
	if upstream == nil {
		panic("batch: upstream fetcher must not be nil")
	}
	if workers < 1 {
		workers = 4
	}
	return &BatchFetcher{upstream: upstream, workers: workers}
}

// FetchAll resolves all requests concurrently and returns one Result per Request.
// Results are returned in the same order as the input slice.
func (b *BatchFetcher) FetchAll(ctx context.Context, reqs []Request) []Result {
	results := make([]Result, len(reqs))

	type indexedReq struct {
		idx int
		req Request
	}

	work := make(chan indexedReq, len(reqs))
	for i, r := range reqs {
		work <- indexedReq{idx: i, req: r}
	}
	close(work)

	var wg sync.WaitGroup
	for w := 0; w < b.workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for item := range work {
				results[item.idx] = b.fetch(ctx, item.req)
			}
		}()
	}
	wg.Wait()
	return results
}

func (b *BatchFetcher) fetch(ctx context.Context, req Request) Result {
	data, err := b.upstream.GetSecretData(ctx, req.Path)
	if err != nil {
		return Result{Request: req, Err: err}
	}
	raw, ok := data[req.Field]
	if !ok {
		return Result{Request: req, Err: fmt.Errorf("batch: field %q not found at path %q", req.Field, req.Path)}
	}
	val, ok := raw.(string)
	if !ok {
		return Result{Request: req, Err: fmt.Errorf("batch: field %q at path %q is not a string", req.Field, req.Path)}
	}
	return Result{Request: req, Value: val}
}
