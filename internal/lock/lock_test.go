package lock_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/your-org/vaultenv/internal/lock"
)

func TestNew_DefaultTimeout(t *testing.T) {
	l := lock.New(0)
	if l == nil {
		t.Fatal("expected non-nil Locker")
	}
}

func TestLock_AcquireAndRelease(t *testing.T) {
	l := lock.New(time.Second)
	ctx := context.Background()

	if err := l.Lock(ctx, "secret/foo"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	l.Unlock("secret/foo")
}

func TestLock_MutualExclusion(t *testing.T) {
	l := lock.New(2 * time.Second)
	ctx := context.Background()

	var order []int
	var mu sync.Mutex

	if err := l.Lock(ctx, "key"); err != nil {
		t.Fatalf("first lock: %v", err)
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		if err := l.Lock(ctx, "key"); err != nil {
			return
		}
		mu.Lock()
		order = append(order, 2)
		mu.Unlock()
		l.Unlock("key")
	}()

	time.Sleep(20 * time.Millisecond)
	mu.Lock()
	order = append(order, 1)
	mu.Unlock()
	l.Unlock("key")

	<-done

	if len(order) != 2 || order[0] != 1 || order[1] != 2 {
		t.Fatalf("expected [1 2], got %v", order)
	}
}

func TestLock_Timeout(t *testing.T) {
	l := lock.New(50 * time.Millisecond)
	ctx := context.Background()

	if err := l.Lock(ctx, "blocked"); err != nil {
		t.Fatalf("first lock: %v", err)
	}
	defer l.Unlock("blocked")

	err := l.Lock(ctx, "blocked")
	if err != lock.ErrTimeout {
		t.Fatalf("expected ErrTimeout, got %v", err)
	}
}

func TestLock_ContextCancelled(t *testing.T) {
	l := lock.New(5 * time.Second)
	ctx, cancel := context.WithCancel(context.Background())

	if err := l.Lock(ctx, "held"); err != nil {
		t.Fatalf("first lock: %v", err)
	}
	defer l.Unlock("held")

	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()

	err := l.Lock(ctx, "held")
	if err != lock.ErrCancelled {
		t.Fatalf("expected ErrCancelled, got %v", err)
	}
}

func TestUnlock_UnknownKey(t *testing.T) {
	l := lock.New(time.Second)
	// Should not panic.
	l.Unlock("nonexistent")
}

func TestLock_DifferentKeysConcurrent(t *testing.T) {
	l := lock.New(time.Second)
	ctx := context.Background()

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		key := fmt.Sprintf("key/%d", i)
		go func(k string) {
			defer wg.Done()
			if err := l.Lock(ctx, k); err != nil {
				t.Errorf("lock %s: %v", k, err)
				return
			}
			time.Sleep(5 * time.Millisecond)
			l.Unlock(k)
		}(key)
	}
	wg.Wait()
}
