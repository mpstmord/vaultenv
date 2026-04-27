package hedge_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/your-org/vaultenv/internal/hedge"
)

type mockFetcher struct {
	calls atomic.Int32
	delay time.Duration
	data  map[string]interface{}
	err   error
}

func (m *mockFetcher) GetSecretData(ctx context.Context, path string) (map[string]interface{}, error) {
	m.calls.Add(1)
	if m.delay > 0 {
		select {
		case <-time.After(m.delay):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	return m.data, m.err
}

func TestNew_PanicsOnNilUpstream(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil upstream")
		}
	}()
	hedge.New(nil, 0)
}

func TestNew_DefaultDelay(t *testing.T) {
	m := &mockFetcher{data: map[string]interface{}{"k": "v"}}
	h := hedge.New(m, 0) // 0 triggers default
	data, err := h.GetSecretData(context.Background(), "secret/foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data["k"] != "v" {
		t.Fatalf("unexpected data: %v", data)
	}
}

func TestGetSecretData_FastFirstResponse(t *testing.T) {
	m := &mockFetcher{data: map[string]interface{}{"secret": "value"}}
	h := hedge.New(m, 200*time.Millisecond)

	data, err := h.GetSecretData(context.Background(), "secret/fast")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data["secret"] != "value" {
		t.Fatalf("unexpected data: %v", data)
	}
	// Only one call expected because first replied before hedge delay.
	if n := m.calls.Load(); n != 1 {
		t.Fatalf("expected 1 call, got %d", n)
	}
}

func TestGetSecretData_HedgeFiresSecondRequest(t *testing.T) {
	// First call is slow; hedge fires a second which also returns.
	m := &mockFetcher{
		delay: 150 * time.Millisecond,
		data:  map[string]interface{}{"x": "1"},
	}
	h := hedge.New(m, 20*time.Millisecond)

	start := time.Now()
	data, err := h.GetSecretData(context.Background(), "secret/slow")
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data["x"] != "1" {
		t.Fatalf("unexpected data: %v", data)
	}
	if elapsed >= 150*time.Millisecond {
		t.Fatalf("hedge did not shorten latency: %v", elapsed)
	}
	if n := m.calls.Load(); n < 2 {
		t.Fatalf("expected hedge to fire second request, calls=%d", n)
	}
}

func TestGetSecretData_PropagatesError(t *testing.T) {
	sentinel := errors.New("vault unavailable")
	m := &mockFetcher{err: sentinel}
	h := hedge.New(m, 10*time.Millisecond)

	_, err := h.GetSecretData(context.Background(), "secret/err")
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
}

func TestGetSecretData_ContextCancelled(t *testing.T) {
	m := &mockFetcher{delay: 5 * time.Second}
	h := hedge.New(m, 10*time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()

	_, err := h.GetSecretData(ctx, "secret/cancel")
	if err == nil {
		t.Fatal("expected context error")
	}
}
