package dedupe_test

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/your-org/vaultenv/internal/dedupe"
)

// staticFetcher returns a fixed payload after an optional delay.
type staticFetcher struct {
	calls int64
	delay time.Duration
	data  map[string]interface{}
	err   error
}

func (f *staticFetcher) GetSecretData(_ context.Context, _ string) (map[string]interface{}, error) {
	atomic.AddInt64(&f.calls, 1)
	if f.delay > 0 {
		time.Sleep(f.delay)
	}
	return f.data, f.err
}

func TestNew_NilUpstreamPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil upstream")
		}
	}()
	dedupe.New(nil)
}

func TestGetSecretData_SingleCall(t *testing.T) {
	upstream := &staticFetcher{data: map[string]interface{}{"k": "v"}}
	d := dedupe.New(upstream)

	got, err := d.GetSecretData(context.Background(), "secret/foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["k"] != "v" {
		t.Fatalf("unexpected value: %v", got["k"])
	}
	if atomic.LoadInt64(&upstream.calls) != 1 {
		t.Fatalf("expected 1 upstream call, got %d", upstream.calls)
	}
}

func TestGetSecretData_DeduplicatesConcurrent(t *testing.T) {
	upstream := &staticFetcher{
		delay: 40 * time.Millisecond,
		data:  map[string]interface{}{"token": "abc"},
	}
	d := dedupe.New(upstream)

	const goroutines = 10
	var wg sync.WaitGroup
	results := make([]map[string]interface{}, goroutines)
	errs := make([]error, goroutines)

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			results[idx], errs[idx] = d.GetSecretData(context.Background(), "secret/shared")
		}(i)
	}
	wg.Wait()

	for i, err := range errs {
		if err != nil {
			t.Errorf("goroutine %d: unexpected error: %v", i, err)
		}
	}
	if calls := atomic.LoadInt64(&upstream.calls); calls > 3 {
		t.Errorf("expected deduplication: got %d upstream calls for %d goroutines", calls, goroutines)
	}
}

func TestGetSecretData_PropagatesError(t *testing.T) {
	expected := errors.New("vault unavailable")
	upstream := &staticFetcher{err: expected}
	d := dedupe.New(upstream)

	_, err := d.GetSecretData(context.Background(), "secret/bad")
	if !errors.Is(err, expected) {
		t.Fatalf("expected %v, got %v", expected, err)
	}
}

func TestGetSecretData_DifferentPathsAreIndependent(t *testing.T) {
	upstream := &staticFetcher{data: map[string]interface{}{"x": "1"}}
	d := dedupe.New(upstream)

	d.GetSecretData(context.Background(), "secret/a") //nolint:errcheck
	d.GetSecretData(context.Background(), "secret/b") //nolint:errcheck

	if calls := atomic.LoadInt64(&upstream.calls); calls != 2 {
		t.Fatalf("expected 2 upstream calls for different paths, got %d", calls)
	}
}
