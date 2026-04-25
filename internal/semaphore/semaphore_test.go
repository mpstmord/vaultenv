package semaphore_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/vaultenv/vaultenv/internal/semaphore"
)

func TestNew_PanicsOnZeroCapacity(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for zero capacity")
		}
	}()
	semaphore.New(0)
}

func TestNew_Valid(t *testing.T) {
	s := semaphore.New(3)
	if s.Capacity() != 3 {
		t.Fatalf("expected capacity 3, got %d", s.Capacity())
	}
	if s.Available() != 3 {
		t.Fatalf("expected 3 available slots, got %d", s.Available())
	}
}

func TestAcquire_And_Release(t *testing.T) {
	s := semaphore.New(2)
	ctx := context.Background()

	if err := s.Acquire(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Available() != 1 {
		t.Fatalf("expected 1 available, got %d", s.Available())
	}
	s.Release()
	if s.Available() != 2 {
		t.Fatalf("expected 2 available after release, got %d", s.Available())
	}
}

func TestAcquire_BlocksWhenFull(t *testing.T) {
	s := semaphore.New(1)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_ = s.Acquire(context.Background())
	err := s.Acquire(ctx)
	if err == nil {
		t.Fatal("expected error when semaphore full")
	}
}

func TestAcquireTimeout_TimesOut(t *testing.T) {
	s := semaphore.New(1)
	_ = s.Acquire(context.Background())

	err := s.AcquireTimeout(20 * time.Millisecond)
	if err != semaphore.ErrTimeout {
		t.Fatalf("expected ErrTimeout, got %v", err)
	}
}

func TestAcquireTimeout_Success(t *testing.T) {
	s := semaphore.New(2)
	err := s.AcquireTimeout(100 * time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s.Release()
}

func TestRelease_PanicsWithoutAcquire(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on excess Release")
		}
	}()
	s := semaphore.New(1)
	s.Release()
}

func TestSemaphore_ConcurrentAcquire(t *testing.T) {
	s := semaphore.New(3)
	var wg sync.WaitGroup
	for i := 0; i < 9; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = s.Acquire(context.Background())
			time.Sleep(10 * time.Millisecond)
			s.Release()
		}()
	}
	wg.Wait()
	if s.Available() != 3 {
		t.Fatalf("expected all slots free after concurrent use, got %d", s.Available())
	}
}
