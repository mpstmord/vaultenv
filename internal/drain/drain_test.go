package drain_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/your-org/vaultenv/internal/drain"
)

func TestNew_DefaultTimeout(t *testing.T) {
	d := drain.New(0)
	if d == nil {
		t.Fatal("expected non-nil Drainer")
	}
}

func TestAcquire_BeforeDrain(t *testing.T) {
	d := drain.New(time.Second)
	if !d.Acquire() {
		t.Fatal("expected Acquire to succeed before Drain")
	}
	d.Release()
}

func TestAcquire_AfterDrain(t *testing.T) {
	d := drain.New(time.Second)

	// Drain with no in-flight ops — completes immediately.
	if err := d.Drain(context.Background()); err != nil {
		t.Fatalf("unexpected drain error: %v", err)
	}

	if d.Acquire() {
		t.Fatal("expected Acquire to fail after Drain")
	}
}

func TestDrain_WaitsForRelease(t *testing.T) {
	d := drain.New(2 * time.Second)

	if !d.Acquire() {
		t.Fatal("Acquire failed")
	}

	var wg sync.WaitGroup
	wg.Add(1)
	var drainErr error
	go func() {
		defer wg.Done()
		drainErr = d.Drain(context.Background())
	}()

	// Give goroutine time to block on Drain.
	time.Sleep(50 * time.Millisecond)
	if !d.Closed() {
		t.Error("expected Drainer to be closed after Drain called")
	}

	d.Release()
	wg.Wait()

	if drainErr != nil {
		t.Fatalf("expected nil error, got: %v", drainErr)
	}
}

func TestDrain_Timeout(t *testing.T) {
	d := drain.New(50 * time.Millisecond)

	if !d.Acquire() {
		t.Fatal("Acquire failed")
	}
	// Never call Release — timeout should fire.
	defer d.Release()

	err := d.Drain(context.Background())
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
}

func TestClosed_FalseInitially(t *testing.T) {
	d := drain.New(time.Second)
	if d.Closed() {
		t.Fatal("expected Closed() == false before Drain")
	}
}

func TestDrain_ConcurrentAcquire(t *testing.T) {
	d := drain.New(2 * time.Second)
	const workers = 10

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if d.Acquire() {
				time.Sleep(10 * time.Millisecond)
				d.Release()
			}
		}()
	}

	wg.Wait()
	if err := d.Drain(context.Background()); err != nil {
		t.Fatalf("unexpected drain error: %v", err)
	}
}
