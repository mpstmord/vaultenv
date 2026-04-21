package middleware_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/your-org/vaultenv/internal/middleware"
	"github.com/your-org/vaultenv/internal/retry"
)

type countingFetcher struct {
	calls   atomic.Int32
	failFor int32
	data    map[string]interface{}
}

func (c *countingFetcher) GetSecretData(_ context.Context, _ string) (map[string]interface{}, error) {
	n := c.calls.Add(1)
	if n <= c.failFor {
		return nil, errors.New("transient error")
	}
	return c.data, nil
}

func fastRetryPolicy(attempts int) retry.Policy {
	return retry.Policy{
		Attempts:   attempts,
		MinBackoff: time.Millisecond,
		MaxBackoff: time.Millisecond,
	}
}

func TestRetryMiddleware_SuccessOnFirstAttempt(t *testing.T) {
	upstream := &countingFetcher{data: map[string]interface{}{"k": "v"}}
	m := middleware.NewRetryMiddleware(fastRetryPolicy(3))
	f := m.Wrap(upstream)

	data, err := f.GetSecretData(context.Background(), "secret/foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data["k"] != "v" {
		t.Errorf("expected k=v, got %v", data["k"])
	}
	if upstream.calls.Load() != 1 {
		t.Errorf("expected 1 call, got %d", upstream.calls.Load())
	}
}

func TestRetryMiddleware_RetriesAndSucceeds(t *testing.T) {
	upstream := &countingFetcher{failFor: 2, data: map[string]interface{}{"token": "abc"}}
	m := middleware.NewRetryMiddleware(fastRetryPolicy(5))
	f := m.Wrap(upstream)

	_, err := f.GetSecretData(context.Background(), "secret/bar")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if upstream.calls.Load() != 3 {
		t.Errorf("expected 3 calls, got %d", upstream.calls.Load())
	}
}

func TestRetryMiddleware_ExhaustsAttempts(t *testing.T) {
	upstream := &countingFetcher{failFor: 10}
	m := middleware.NewRetryMiddleware(fastRetryPolicy(3))
	f := m.Wrap(upstream)

	_, err := f.GetSecretData(context.Background(), "secret/baz")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if upstream.calls.Load() != 3 {
		t.Errorf("expected 3 calls, got %d", upstream.calls.Load())
	}
}

func TestRetryMiddleware_DefaultPolicy(t *testing.T) {
	// Zero-value policy should fall back to DefaultPolicy without panicking.
	upstream := &countingFetcher{data: map[string]interface{}{}}
	m := middleware.NewRetryMiddleware(retry.Policy{})
	f := m.Wrap(upstream)

	_, err := f.GetSecretData(context.Background(), "secret/qux")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRetryMiddleware_ContextCancelled(t *testing.T) {
	upstream := &countingFetcher{failFor: 100}
	m := middleware.NewRetryMiddleware(fastRetryPolicy(10))
	f := m.Wrap(upstream)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := f.GetSecretData(ctx, "secret/cancel")
	if err == nil {
		t.Fatal("expected error due to cancelled context")
	}
}
