package retry_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/your-org/vaultenv/internal/retry"
)

var errTransient = errors.New("transient error")

func fastPolicy() retry.Policy {
	return retry.Policy{
		MaxAttempts: 3,
		BaseDelay:   1 * time.Millisecond,
		MaxDelay:    5 * time.Millisecond,
	}
}

func TestDo_SuccessOnFirstAttempt(t *testing.T) {
	calls := 0
	err := retry.Do(context.Background(), fastPolicy(), func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestDo_RetriesAndSucceeds(t *testing.T) {
	var calls atomic.Int32
	err := retry.Do(context.Background(), fastPolicy(), func() error {
		if calls.Add(1) < 3 {
			return errTransient
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if calls.Load() != 3 {
		t.Fatalf("expected 3 calls, got %d", calls.Load())
	}
}

func TestDo_ExhaustsAttempts(t *testing.T) {
	calls := 0
	p := fastPolicy()
	err := retry.Do(context.Background(), p, func() error {
		calls++
		return errTransient
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, retry.ErrMaxAttempts) {
		t.Fatalf("expected ErrMaxAttempts in chain, got %v", err)
	}
	if !errors.Is(err, errTransient) {
		t.Fatalf("expected wrapped transient error, got %v", err)
	}
	if calls != p.MaxAttempts {
		t.Fatalf("expected %d calls, got %d", p.MaxAttempts, calls)
	}
}

func TestDo_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	calls := 0
	err := retry.Do(ctx, fastPolicy(), func() error {
		calls++
		return errTransient
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if calls != 0 {
		t.Fatalf("expected 0 calls for pre-cancelled context, got %d", calls)
	}
}

func TestDefaultPolicy_Sensible(t *testing.T) {
	p := retry.DefaultPolicy()
	if p.MaxAttempts <= 0 {
		t.Fatalf("MaxAttempts must be positive, got %d", p.MaxAttempts)
	}
	if p.BaseDelay <= 0 {
		t.Fatalf("BaseDelay must be positive, got %v", p.BaseDelay)
	}
	if p.MaxDelay < p.BaseDelay {
		t.Fatalf("MaxDelay must be >= BaseDelay")
	}
}
