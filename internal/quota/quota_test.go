package quota_test

import (
	"testing"
	"time"

	"github.com/your-org/vaultenv/internal/quota"
)

func TestNew_InvalidLimit(t *testing.T) {
	_, err := quota.New(0, time.Minute)
	if err == nil {
		t.Fatal("expected error for zero limit")
	}
}

func TestNew_InvalidWindow(t *testing.T) {
	_, err := quota.New(10, 0)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestAllow_WithinLimit(t *testing.T) {
	l, err := quota.New(3, time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 0; i < 3; i++ {
		if err := l.Allow("secret/data/foo"); err != nil {
			t.Fatalf("request %d: unexpected error: %v", i+1, err)
		}
	}
}

func TestAllow_ExceedsLimit(t *testing.T) {
	l, _ := quota.New(2, time.Minute)
	_ = l.Allow("secret/data/bar")
	_ = l.Allow("secret/data/bar")
	err := l.Allow("secret/data/bar")
	if err == nil {
		t.Fatal("expected quota exceeded error")
	}
	var qe *quota.ErrQuotaExceeded
	if ok := errorAs(err, &qe); !ok {
		t.Fatalf("expected *ErrQuotaExceeded, got %T", err)
	}
	if qe.Path != "secret/data/bar" {
		t.Errorf("unexpected path: %s", qe.Path)
	}
}

func TestAllow_WindowReset(t *testing.T) {
	l, _ := quota.New(1, 50*time.Millisecond)
	if err := l.Allow("secret/data/baz"); err != nil {
		t.Fatalf("first request failed: %v", err)
	}
	if err := l.Allow("secret/data/baz"); err == nil {
		t.Fatal("expected quota exceeded before window reset")
	}
	time.Sleep(60 * time.Millisecond)
	if err := l.Allow("secret/data/baz"); err != nil {
		t.Fatalf("after window reset: %v", err)
	}
}

func TestAllow_IndependentPaths(t *testing.T) {
	l, _ := quota.New(1, time.Minute)
	if err := l.Allow("secret/data/a"); err != nil {
		t.Fatalf("path a: %v", err)
	}
	if err := l.Allow("secret/data/b"); err != nil {
		t.Fatalf("path b: %v", err)
	}
}

func TestReset_ClearsState(t *testing.T) {
	l, _ := quota.New(1, time.Minute)
	_ = l.Allow("secret/data/x")
	l.Reset("secret/data/x")
	if err := l.Allow("secret/data/x"); err != nil {
		t.Fatalf("after reset: %v", err)
	}
}

func TestStats_ReturnsCount(t *testing.T) {
	l, _ := quota.New(5, time.Minute)
	_ = l.Allow("secret/data/q")
	_ = l.Allow("secret/data/q")
	count, end := l.Stats("secret/data/q")
	if count != 2 {
		t.Errorf("expected count 2, got %d", count)
	}
	if end.IsZero() {
		t.Error("expected non-zero window end")
	}
}

func TestStats_MissingPath(t *testing.T) {
	l, _ := quota.New(5, time.Minute)
	count, end := l.Stats("secret/data/missing")
	if count != 0 || !end.IsZero() {
		t.Errorf("expected zero stats, got count=%d end=%v", count, end)
	}
}

// errorAs is a minimal stand-in for errors.As without importing errors.
func errorAs(err error, target **quota.ErrQuotaExceeded) bool {
	if e, ok := err.(*quota.ErrQuotaExceeded); ok {
		*target = e
		return true
	}
	return false
}
