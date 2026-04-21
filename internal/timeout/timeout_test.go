package timeout_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/your-org/vaultenv/internal/timeout"
)

// mockFetcher implements timeout.Fetcher for testing.
type mockFetcher struct {
	delay time.Duration
	data  map[string]interface{}
	err   error
}

func (m *mockFetcher) GetSecretData(ctx context.Context, path string) (map[string]interface{}, error) {
	if m.delay > 0 {
		select {
		case <-time.After(m.delay):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	return m.data, m.err
}

func TestNew_NilUpstream(t *testing.T) {
	_, err := timeout.New(nil, 5*time.Second)
	if err == nil {
		t.Fatal("expected error for nil upstream, got nil")
	}
}

func TestNew_NegativeDurationUsesDefault(t *testing.T) {
	f := &mockFetcher{}
	g, err := timeout.New(f, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if g.Duration() != timeout.DefaultTimeout {
		t.Errorf("expected default timeout %s, got %s", timeout.DefaultTimeout, g.Duration())
	}
}

func TestGetSecretData_Success(t *testing.T) {
	want := map[string]interface{}{"api_key": "abc123"}
	f := &mockFetcher{data: want}
	g, _ := timeout.New(f, 5*time.Second)

	got, err := g.GetSecretData(context.Background(), "secret/data/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["api_key"] != want["api_key"] {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestGetSecretData_Timeout(t *testing.T) {
	f := &mockFetcher{delay: 200 * time.Millisecond}
	g, _ := timeout.New(f, 20*time.Millisecond)

	_, err := g.GetSecretData(context.Background(), "secret/data/slow")
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
	if !errors.Is(err, timeout.ErrDeadlineExceeded) {
		t.Errorf("expected ErrDeadlineExceeded, got: %v", err)
	}
}

func TestGetSecretData_UpstreamError(t *testing.T) {
	upstreamErr := errors.New("vault unavailable")
	f := &mockFetcher{err: upstreamErr}
	g, _ := timeout.New(f, 5*time.Second)

	_, err := g.GetSecretData(context.Background(), "secret/data/app")
	if !errors.Is(err, upstreamErr) {
		t.Errorf("expected upstream error, got: %v", err)
	}
}

func TestGetSecretData_ParentContextCancelled(t *testing.T) {
	f := &mockFetcher{delay: 500 * time.Millisecond}
	g, _ := timeout.New(f, 5*time.Second)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := g.GetSecretData(ctx, "secret/data/app")
	if err == nil {
		t.Fatal("expected error from cancelled context, got nil")
	}
}
