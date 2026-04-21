package limiter_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/your-org/vaultenv/internal/limiter"
)

func TestNew_InvalidLimit(t *testing.T) {
	_, err := limiter.New(0)
	if err == nil {
		t.Fatal("expected error for limit=0, got nil")
	}
}

func TestNew_Valid(t *testing.T) {
	l, err := limiter.New(3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l.Capacity() != 3 {
		t.Fatalf("expected capacity 3, got %d", l.Capacity())
	}
}

func TestAcquire_And_Release(t *testing.T) {
	l, _ := limiter.New(2)
	ctx := context.Background()

	if err := l.Acquire(ctx); err != nil {
		t.Fatalf("first acquire failed: %v", err)
	}
	if l.InFlight() != 1 {
		t.Fatalf("expected 1 in-flight, got %d", l.InFlight())
	}
	l.Release()
	if l.InFlight() != 0 {
		t.Fatalf("expected 0 in-flight after release, got %d", l.InFlight())
	}
}

func TestAcquire_BlocksWhenFull(t *testing.T) {
	l, _ := limiter.New(1)
	ctx := context.Background()

	if err := l.Acquire(ctx); err != nil {
		t.Fatalf("first acquire failed: %v", err)
	}

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := l.Acquire(ctxTimeout)
	if err == nil {
		t.Fatal("expected ErrLimitExceeded, got nil")
	}
	if err != limiter.ErrLimitExceeded {
		t.Fatalf("expected ErrLimitExceeded, got %v", err)
	}
	l.Release()
}

func TestLimiter_ConcurrentAcquire(t *testing.T) {
	const cap = 4
	l, _ := limiter.New(cap)
	ctx := context.Background()

	var wg sync.WaitGroup
	for i := 0; i < cap; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := l.Acquire(ctx); err != nil {
				t.Errorf("acquire failed: %v", err)
				return
			}
			time.Sleep(10 * time.Millisecond)
			l.Release()
		}()
	}
	wg.Wait()

	if l.InFlight() != 0 {
		t.Fatalf("expected 0 in-flight after all goroutines done, got %d", l.InFlight())
	}
}
