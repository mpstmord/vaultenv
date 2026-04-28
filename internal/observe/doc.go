// Package observe wraps secret Fetchers with a non-intrusive observer
// hook. Observers receive an Event after each fetch attempt and can be
// used to collect metrics, emit structured logs, or feed distributed
// traces without modifying the underlying fetcher implementation.
//
// Usage:
//
//	counters := &observe.Counters{}
//	f := observe.New(upstream,
//	    observe.MetricsObserver(counters),
//	    observe.LogObserver(os.Stdout),
//	)
package observe
