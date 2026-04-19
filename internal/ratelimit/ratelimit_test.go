package ratelimit

import (
	"testing"
	"time"
)

func TestNew_InvalidRate(t *testing.T) {
	_, err := New(0, 5)
	if err == nil {
		t.Fatal("expected error for zero rate")
	}
}

func TestNew_InvalidBurst(t *testing.T) {
	_, err := New(10, 0)
	if err == nil {
		t.Fatal("expected error for zero burst")
	}
}

func TestAllow_BurstConsumed(t *testing.T) {
	l, err := New(1, 3)
	if err != nil {
		t.Fatal(err)
	}
	// Freeze clock so no refill happens.
	fixed := time.Now()
	l.clock = func() time.Time { return fixed }
	l.lastTick = fixed

	if !l.Allow() {
		t.Error("expected first Allow to succeed")
	}
	if !l.Allow() {
		t.Error("expected second Allow to succeed")
	}
	if !l.Allow() {
		t.Error("expected third Allow to succeed")
	}
	if l.Allow() {
		t.Error("expected fourth Allow to fail (burst exhausted)")
	}
}

func TestAllow_Refill(t *testing.T) {
	l, err := New(2, 1)
	if err != nil {
		t.Fatal(err)
	}
	fixed := time.Now()
	l.clock = func() time.Time { return fixed }
	l.lastTick = fixed

	if !l.Allow() {
		t.Fatal("expected Allow to succeed initially")
	}
	if l.Allow() {
		t.Fatal("expected Allow to fail after burst exhausted")
	}

	// Advance clock by 1 second — should refill 2 tokens, capped at burst=1.
	l.clock = func() time.Time { return fixed.Add(time.Second) }
	if !l.Allow() {
		t.Error("expected Allow to succeed after refill")
	}
}

func TestNew_Valid(t *testing.T) {
	l, err := New(5, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l == nil {
		t.Fatal("expected non-nil limiter")
	}
}
