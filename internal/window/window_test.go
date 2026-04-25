package window

import (
	"testing"
	"time"
)

func TestNew_InvalidWindow(t *testing.T) {
	_, err := New(-time.Second, 10)
	if err == nil {
		t.Fatal("expected error for non-positive window")
	}
}

func TestNew_InvalidBuckets(t *testing.T) {
	_, err := New(time.Second, 0)
	if err == nil {
		t.Fatal("expected error for zero buckets")
	}
}

func TestNew_Valid(t *testing.T) {
	c, err := New(time.Minute, 6)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil counter")
	}
}

func TestTotal_StartsAtZero(t *testing.T) {
	c, _ := New(time.Minute, 6)
	if got := c.Total(); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestAdd_And_Total(t *testing.T) {
	c, _ := New(time.Minute, 6)
	c.Add(3)
	c.Add(7)
	if got := c.Total(); got != 10 {
		t.Fatalf("expected 10, got %d", got)
	}
}

func TestReset_ClearsCount(t *testing.T) {
	c, _ := New(time.Minute, 6)
	c.Add(5)
	c.Reset()
	if got := c.Total(); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestTotal_ExcludesExpiredBuckets(t *testing.T) {
	// Use a very short window so we can force expiry.
	c, _ := New(50*time.Millisecond, 2)
	c.Add(100)
	// Wait longer than the full window so all buckets expire.
	time.Sleep(80 * time.Millisecond)
	if got := c.Total(); got != 0 {
		t.Fatalf("expected 0 after window expiry, got %d", got)
	}
}

func TestAdd_MultipleBuckets(t *testing.T) {
	// Two buckets each 25 ms wide; window = 50 ms.
	c, _ := New(50*time.Millisecond, 2)
	c.Add(10)
	time.Sleep(30 * time.Millisecond) // advance past first bucket boundary
	c.Add(5)
	total := c.Total()
	// Both buckets should still be within the window at this point.
	if total < 5 {
		t.Fatalf("expected at least 5, got %d", total)
	}
}
