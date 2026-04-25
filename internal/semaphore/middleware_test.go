package semaphore_test

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/vaultenv/vaultenv/internal/semaphore"
)

type mockFetcher struct {
	delay time.Duration
	err   error
	calls int32
}

func (m *mockFetcher) GetSecretData(_ context.Context, _ string) (map[string]interface{}, error) {
	atomic.AddInt32(&m.calls, 1)
	if m.delay > 0 {
		time.Sleep(m.delay)
	}
	if m.err != nil {
		return nil, m.err
	}
	return map[string]interface{}{"key": "value"}, nil
}

func TestGuardedFetcher_AllowsWithinLimit(t *testing.T) {
	upstream := &mockFetcher{}
	gf := semaphore.NewGuardedFetcher(upstream, 3)

	data, err := gf.GetSecretData(context.Background(), "secret/test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data["key"] != "value" {
		t.Fatalf("unexpected data: %v", data)
	}
}

func TestGuardedFetcher_PropagatesUpstreamError(t *testing.T) {
	expected := errors.New("vault unavailable")
	upstream := &mockFetcher{err: expected}
	gf := semaphore.NewGuardedFetcher(upstream, 2)

	_, err := gf.GetSecretData(context.Background(), "secret/test")
	if !errors.Is(err, expected) {
		t.Fatalf("expected upstream error, got %v", err)
	}
}

func TestGuardedFetcher_LimitsConcurrency(t *testing.T) {
	var active int32
	var maxActive int32
	var mu sync.Mutex

	upstream := &slowCountingFetcher{active: &active, maxActive: &maxActive, mu: &mu, delay: 30 * time.Millisecond}
	gf := semaphore.NewGuardedFetcher(upstream, 2)

	var wg sync.WaitGroup
	for i := 0; i < 6; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = gf.GetSecretData(context.Background(), "secret/x")
		}()
	}
	wg.Wait()

	if maxActive > 2 {
		t.Fatalf("expected max 2 concurrent, got %d", maxActive)
	}
}

type slowCountingFetcher struct {
	active    *int32
	maxActive *int32
	mu        *sync.Mutex
	delay     time.Duration
}

func (s *slowCountingFetcher) GetSecretData(_ context.Context, _ string) (map[string]interface{}, error) {
	cur := atomic.AddInt32(s.active, 1)
	s.mu.Lock()
	if cur > atomic.LoadInt32(s.maxActive) {
		atomic.StoreInt32(s.maxActive, cur)
	}
	s.mu.Unlock()
	time.Sleep(s.delay)
	atomic.AddInt32(s.active, -1)
	return map[string]interface{}{}, nil
}

func TestGuardedFetcher_ContextCancelled(t *testing.T) {
	upstream := &mockFetcher{delay: 500 * time.Millisecond}
	gf := semaphore.NewGuardedFetcher(upstream, 1)

	// Fill the slot
	_ = context.Background()
	go func() { _, _ = gf.GetSecretData(context.Background(), "secret/fill") }()
	time.Sleep(5 * time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := gf.GetSecretData(ctx, "secret/blocked")
	if err == nil {
		t.Fatal("expected error on cancelled context")
	}
}
