package watch

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

type mockFetcher struct {
	calls  int32
	result map[string]string
	err    error
}

func (m *mockFetcher) FetchSecrets(_ context.Context) (map[string]string, error) {
	atomic.AddInt32(&m.calls, 1)
	return m.result, m.err
}

func TestNew_DefaultInterval(t *testing.T) {
	w := New(&mockFetcher{}, func(_, _ map[string]string) error { return nil }, 0)
	if w.interval != 30*time.Second {
		t.Fatalf("expected 30s default, got %v", w.interval)
	}
}

func TestWatcher_DetectsChange(t *testing.T) {
	fetcher := &mockFetcher{result: map[string]string{"KEY": "v1"}}
	changed := make(chan struct{}, 1)

	w := New(fetcher, func(old, new map[string]string) error {
		changed <- struct{}{}
		return nil
	}, 20*time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	go w.Start(ctx) //nolint:errcheck

	// Allow first poll to settle, then change value
	time.Sleep(40 * time.Millisecond)
	fetcher.result = map[string]string{"KEY": "v2"}

	select {
	case <-changed:
		// success
	case <-time.After(150 * time.Millisecond):
		t.Fatal("expected change handler to be called")
	}
}

func TestWatcher_NoChangeNoHandler(t *testing.T) {
	fetcher := &mockFetcher{result: map[string]string{"KEY": "stable"}}
	handlerCalls := int32(0)

	w := New(fetcher, func(_, _ map[string]string) error {
		atomic.AddInt32(&handlerCalls, 1)
		return nil
	}, 20*time.Millisecond)

	w.last = map[string]string{"KEY": "stable"}

	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()
	go w.Start(ctx) //nolint:errcheck

	time.Sleep(100 * time.Millisecond)
	if atomic.LoadInt32(&handlerCalls) != 0 {
		t.Fatalf("handler should not be called when secrets unchanged")
	}
}

func TestWatcher_FetchError_StopsLoop(t *testing.T) {
	fetchErr := errors.New("vault unavailable")
	fetcher := &mockFetcher{err: fetchErr}

	w := New(fetcher, func(_, _ map[string]string) error { return nil }, 10*time.Millisecond)

	ctx := context.Background()
	errCh := make(chan error, 1)
	go func() { errCh <- w.poll(ctx) }()

	select {
	case err := <-errCh:
		if !errors.Is(err, fetchErr) {
			t.Fatalf("expected fetchErr, got %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for poll error")
	}
}
