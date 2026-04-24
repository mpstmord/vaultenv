package prefetch

import (
	"context"
	"sync"
)

// Fetcher retrieves secret data by path and field.
type Fetcher interface {
	GetSecretData(ctx context.Context, path string) (map[string]interface{}, error)
}

// Result holds the outcome of a prefetch operation for a single path.
type Result struct {
	Path   string
	Data   map[string]interface{}
	Err    error
}

// Prefetcher concurrently fetches a set of secret paths.
type Prefetcher struct {
	fetcher Fetcher
	workers int
}

// New returns a Prefetcher that uses up to workers goroutines.
// If workers < 1 it defaults to 4.
func New(fetcher Fetcher, workers int) *Prefetcher {
	if fetcher == nil {
		panic("prefetch: nil fetcher")
	}
	if workers < 1 {
		workers = 4
	}
	return &Prefetcher{fetcher: fetcher, workers: workers}
}

// FetchAll fetches all paths concurrently and returns one Result per path.
// The order of results matches the order of paths.
func (p *Prefetcher) FetchAll(ctx context.Context, paths []string) []Result {
	results := make([]Result, len(paths))

	type job struct {
		idx  int
		path string
	}

	jobs := make(chan job, len(paths))
	for i, path := range paths {
		jobs <- job{idx: i, path: path}
	}
	close(jobs)

	var wg sync.WaitGroup
	for w := 0; w < p.workers && w < len(paths); w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				data, err := p.fetcher.GetSecretData(ctx, j.path)
				results[j.idx] = Result{Path: j.path, Data: data, Err: err}
			}
		}()
	}
	wg.Wait()
	return results
}
