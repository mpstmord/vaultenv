package circuit

import (
	"context"
	"errors"
	"testing"
	"time"
)

type mockFetcher struct {
	data map[string]interface{}
	err  error
}

func (m *mockFetcher) GetSecretData(_ context.Context, _ string) (map[string]interface{}, error) {
	return m.data, m.err
}

func TestGuardedFetcher_AllowsOnSuccess(t *testing.T) {
	b, _ := New(3, time.Second)
	f := &mockFetcher{data: map[string]interface{}{"key": "val"}}
	gf := NewGuardedFetcher(f, b)

	data, err := gf.GetSecretData(context.Background(), "secret/foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data["key"] != "val" {
		t.Errorf("unexpected data: %v", data)
	}
	if b.State() != StateClosed {
		t.Errorf("expected closed after success, got %s", b.State())
	}
}

func TestGuardedFetcher_RecordsFailure(t *testing.T) {
	b, _ := New(2, time.Second)
	f := &mockFetcher{err: errors.New("vault unavailable")}
	gf := NewGuardedFetcher(f, b)

	for i := 0; i < 2; i++ {
		_, _ = gf.GetSecretData(context.Background(), "secret/foo")
	}
	if b.State() != StateOpen {
		t.Errorf("expected open after threshold failures, got %s", b.State())
	}
}

func TestGuardedFetcher_BlocksWhenOpen(t *testing.T) {
	b, _ := New(1, time.Hour)
	b.RecordFailure() // force open
	f := &mockFetcher{data: map[string]interface{}{"k": "v"}}
	gf := NewGuardedFetcher(f, b)

	_, err := gf.GetSecretData(context.Background(), "secret/bar")
	if !errors.Is(err, ErrOpen) {
		t.Errorf("expected ErrOpen, got %v", err)
	}
}

func TestGuardedFetcher_PropagatesUpstreamError(t *testing.T) {
	b, _ := New(5, time.Second)
	upstreamErr := errors.New("not found")
	f := &mockFetcher{err: upstreamErr}
	gf := NewGuardedFetcher(f, b)

	_, err := gf.GetSecretData(context.Background(), "secret/missing")
	if !errors.Is(err, upstreamErr) {
		t.Errorf("expected upstream error, got %v", err)
	}
}
