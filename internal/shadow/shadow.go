// Package shadow provides a dual-read fetcher that compares results from a
// primary and a secondary secret source, logging any discrepancies without
// affecting the value returned to the caller.
package shadow

import (
	"context"
	"fmt"
	"log"
	"reflect"
)

// Fetcher is the interface satisfied by any secret backend.
type Fetcher interface {
	GetSecretData(ctx context.Context, path string) (map[string]interface{}, error)
}

// Logger is a minimal logging interface.
type Logger interface {
	Printf(format string, args ...interface{})
}

// Shadower wraps a primary Fetcher and silently mirrors requests to a secondary
// Fetcher, logging differences. The primary result is always returned.
type Shadower struct {
	primary   Fetcher
	secondary Fetcher
	logger    Logger
}

// New creates a Shadower. Both primary and secondary must be non-nil.
func New(primary, secondary Fetcher, logger Logger) *Shadower {
	if primary == nil {
		panic("shadow: primary fetcher must not be nil")
	}
	if secondary == nil {
		panic("shadow: secondary fetcher must not be nil")
	}
	if logger == nil {
		logger = log.Default()
	}
	return &Shadower{primary: primary, secondary: secondary, logger: logger}
}

// GetSecretData fetches from both sources concurrently and returns the primary
// result. Discrepancies or secondary errors are logged.
func (s *Shadower) GetSecretData(ctx context.Context, path string) (map[string]interface{}, error) {
	type result struct {
		data map[string]interface{}
		err  error
	}

	secondaryCh := make(chan result, 1)
	go func() {
		d, err := s.secondary.GetSecretData(ctx, path)
		secondaryCh <- result{data: d, err: err}
	}()

	primaryData, primaryErr := s.primary.GetSecretData(ctx, path)

	sec := <-secondaryCh
	if sec.err != nil {
		s.logger.Printf("shadow: secondary error for path %q: %v", path, sec.err)
	} else if primaryErr == nil && !reflect.DeepEqual(primaryData, sec.data) {
		s.logger.Printf("shadow: mismatch for path %q: primary=%v secondary=%v",
			path, primaryData, sec.data)
	}

	_ = fmt.Sprintf // keep fmt imported for potential future use
	return primaryData, primaryErr
}
