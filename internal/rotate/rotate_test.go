package rotate

import (
	"context"
	"errors"
	"testing"
	"time"
)

type mockFetcher struct {
	calls int
	data  map[string]interface{}
	err   error
}

func (m *mockFetcher) GetSecretData(_ context.Context, _ string) (map[string]interface{}, error) {
	m.calls++
	return m.data, m.err
}

func TestNew_DefaultInterval(t *testing.T) {
	r := New(&mockFetcher{}, 0, func(_ string, _ map[string]interface{}) error { return nil }, nil)
	if r.interval != 30*time.Second {
		t.Fatalf("expected 30s, got %v", r.interval)
	}
}

func TestRotator_HandlerCalledOnChange(t *testing.T) {
	f := &mockFetcher{data: map[string]interface{}{"key": "v1"}}
	var called int
	h := func(_ string, _ map[string]interface{}) error { called++; return nil }
	r := New(f, time.Minute, h, []string{"secret/app"})
	ctx := context.Background()
	if err := r.poll(ctx); err != nil {
		t.Fatal(err)
	}
	if called != 1 {
		t.Fatalf("expected 1 call, got %d", called)
	}
	// Same data — no additional call.
	if err := r.poll(ctx); err != nil {
		t.Fatal(err)
	}
	if called != 1 {
		t.Fatalf("expected still 1 call, got %d", called)
	}
	// Changed data.
	f.data = map[string]interface{}{"key": "v2"}
	if err := r.poll(ctx); err != nil {
		t.Fatal(err)
	}
	if called != 2 {
		t.Fatalf("expected 2 calls, got %d", called)
	}
}

func TestRotator_FetchError(t *testing.T) {
	f := &mockFetcher{err: errors.New("vault down")}
	r := New(f, time.Minute, func(_ string, _ map[string]interface{}) error { return nil }, []string{"secret/x"})
	if err := r.poll(context.Background()); err == nil {
		t.Fatal("expected error")
	}
}

func TestRotator_StartCancels(t *testing.T) {
	f := &mockFetcher{data: map[string]interface{}{}}
	r := New(f, 10*time.Millisecond, func(_ string, _ map[string]interface{}) error { return nil }, []string{})
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Millisecond)
	defer cancel()
	if err := r.Start(ctx); err != nil {
		t.Fatal(err)
	}
}
