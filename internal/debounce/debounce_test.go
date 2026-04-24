package debounce_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/your-org/vaultenv/internal/debounce"
)

func TestNewDebouncer_NilFnPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil fn")
		}
	}()
	debounce.NewDebouncer(50*time.Millisecond, nil)
}

func TestNewDebouncer_ZeroWaitUsesDefault(t *testing.T) {
	var called int32
	d := debounce.NewDebouncer(0, func() { atomic.AddInt32(&called, 1) })
	defer d.Stop()
	// Just verify it was created without panic.
	if d == nil {
		t.Fatal("expected non-nil debouncer")
	}
}

func TestDebouncer_FiresAfterWait(t *testing.T) {
	var called int32
	wait := 40 * time.Millisecond
	d := debounce.NewDebouncer(wait, func() { atomic.AddInt32(&called, 1) })
	defer d.Stop()

	d.Trigger()

	time.Sleep(wait + 30*time.Millisecond)

	if n := atomic.LoadInt32(&called); n != 1 {
		t.Fatalf("expected fn called once, got %d", n)
	}
}

func TestDebouncer_CoalescesRapidTriggers(t *testing.T) {
	var called int32
	wait := 60 * time.Millisecond
	d := debounce.NewDebouncer(wait, func() { atomic.AddInt32(&called, 1) })
	defer d.Stop()

	for i := 0; i < 5; i++ {
		d.Trigger()
		time.Sleep(10 * time.Millisecond)
	}

	time.Sleep(wait + 40*time.Millisecond)

	if n := atomic.LoadInt32(&called); n != 1 {
		t.Fatalf("expected exactly 1 call after coalescing, got %d", n)
	}
}

func TestDebouncer_PendingBeforeAndAfterFire(t *testing.T) {
	wait := 60 * time.Millisecond
	d := debounce.NewDebouncer(wait, func() {})
	defer d.Stop()

	if d.Pending() {
		t.Fatal("should not be pending before first trigger")
	}

	d.Trigger()
	if !d.Pending() {
		t.Fatal("should be pending after trigger")
	}

	time.Sleep(wait + 30*time.Millisecond)
	if d.Pending() {
		t.Fatal("should not be pending after fn fired")
	}
}

func TestDebouncer_StopCancelsPending(t *testing.T) {
	var called int32
	wait := 60 * time.Millisecond
	d := debounce.NewDebouncer(wait, func() { atomic.AddInt32(&called, 1) })

	d.Trigger()
	d.Stop()

	time.Sleep(wait + 30*time.Millisecond)

	if n := atomic.LoadInt32(&called); n != 0 {
		t.Fatalf("expected fn not called after Stop, got %d calls", n)
	}
	if d.Pending() {
		t.Fatal("should not be pending after Stop")
	}
}
