package middleware_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/yourusername/vaultenv/internal/middleware"
)

// stubFetcher returns a fixed result for use in middleware tests.
type rateLimitStubFetcher struct {
	data map[string]interface{}
	err  error
	calls int
}

func (s *rateLimitStubFetcher) GetSecretData(ctx context.Context, path string) (map[string]interface{}, error) {
	s.calls++
	return s.data, s.err
}

func TestRateLimitMiddleware_AllowsWithinBurst(t *testing.T) {
	upstream := &rateLimitStubFetcher{
		data: map[string]interface{}{"key": "value"},
	}

	// rate=10/s, burst=5 — first 5 calls should succeed immediately
	mw := middleware.NewRateLimitMiddleware(10, 5)
	wrapped := mw(upstream)

	ctx := context.Background()
	for i := 0; i < 5; i++ {
		data, err := wrapped.GetSecretData(ctx, "secret/test")
		if err != nil {
			t.Fatalf("call %d: unexpected error: %v", i+1, err)
		}
		if data["key"] != "value" {
			t.Fatalf("call %d: unexpected data: %v", i+1, data)
		}
	}

	if upstream.calls != 5 {
		t.Errorf("expected 5 upstream calls, got %d", upstream.calls)
	}
}

func TestRateLimitMiddleware_BlocksWhenExhausted(t *testing.T) {
	upstream := &rateLimitStubFetcher{
		data: map[string]interface{}{"key": "value"},
	}

	// rate=1/s, burst=1 — second immediate call should be rate-limited
	mw := middleware.NewRateLimitMiddleware(1, 1)
	wrapped := mw(upstream)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	// First call consumes the burst token.
	_, err := wrapped.GetSecretData(ctx, "secret/test")
	if err != nil {
		t.Fatalf("first call: unexpected error: %v", err)
	}

	// Second call should fail because context times out before a token is available.
	_, err = wrapped.GetSecretData(ctx, "secret/test")
	if err == nil {
		t.Fatal("expected rate-limit error on second call, got nil")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected DeadlineExceeded, got: %v", err)
	}
}

func TestRateLimitMiddleware_PropagatesUpstreamError(t *testing.T) {
	sentinel := errors.New("upstream failure")
	upstream := &rateLimitStubFetcher{err: sentinel}

	mw := middleware.NewRateLimitMiddleware(100, 10)
	wrapped := mw(upstream)

	_, err := wrapped.GetSecretData(context.Background(), "secret/test")
	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got: %v", err)
	}
}

func TestRateLimitMiddleware_ContextCancelledBeforeToken(t *testing.T) {
	upstream := &rateLimitStubFetcher{
		data: map[string]interface{}{"x": "y"},
	}

	// Very slow refill: 1 token per 10 seconds, burst=1
	mw := middleware.NewRateLimitMiddleware(0.1, 1)
	wrapped := mw(upstream)

	ctx, cancel := context.WithCancel(context.Background())

	// Consume the single burst token.
	_, _ = wrapped.GetSecretData(ctx, "secret/test")

	// Cancel the context before the next token would be available.
	cancel()

	_, err := wrapped.GetSecretData(ctx, "secret/test")
	if err == nil {
		t.Fatal("expected error after context cancellation, got nil")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got: %v", err)
	}
}
