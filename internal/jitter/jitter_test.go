package jitter_test

import (
	"testing"
	"time"

	"github.com/your-org/vaultenv/internal/jitter"
)

// fixed returns a Source that always produces v.
func fixed(v float64) jitter.Source {
	return func() float64 { return v }
}

func TestJitter_ZeroFactor(t *testing.T) {
	base := 10 * time.Second
	got := jitter.Jitter(base, 0, fixed(0.9))
	if got != base {
		t.Fatalf("expected %v, got %v", base, got)
	}
}

func TestJitter_NegativeFactor(t *testing.T) {
	base := 5 * time.Second
	got := jitter.Jitter(base, -1, fixed(1.0))
	if got != base {
		t.Fatalf("expected %v, got %v", base, got)
	}
}

func TestJitter_FactorClampedToOne(t *testing.T) {
	base := 4 * time.Second
	// factor=2 should be clamped to 1; with src=0.5 delta = base*1*0.5 = 2s
	got := jitter.Jitter(base, 2.0, fixed(0.5))
	want := base + 2*time.Second
	if got != want {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestJitter_Normal(t *testing.T) {
	base := 10 * time.Second
	// factor=0.5, src=1.0 => delta = 10s * 0.5 * 1.0 = 5s
	got := jitter.Jitter(base, 0.5, fixed(1.0))
	want := 15 * time.Second
	if got != want {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestDefault_ReturnsAtLeastBase(t *testing.T) {
	base := 1 * time.Second
	for i := 0; i < 20; i++ {
		got := jitter.Default(base, 0.3)
		if got < base {
			t.Fatalf("Default returned %v which is less than base %v", got, base)
		}
	}
}

func TestFull_ZeroBase(t *testing.T) {
	got := jitter.Full(0, fixed(0.9))
	if got != 0 {
		t.Fatalf("expected 0, got %v", got)
	}
}

func TestFull_Distribution(t *testing.T) {
	base := 10 * time.Second
	// src=0.75 => 2*base*0.75 = 15s
	got := jitter.Full(base, fixed(0.75))
	want := 15 * time.Second
	if got != want {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestEqual_ZeroBase(t *testing.T) {
	got := jitter.Equal(0, fixed(0.5))
	if got != 0 {
		t.Fatalf("expected 0, got %v", got)
	}
}

func TestEqual_Distribution(t *testing.T) {
	base := 10 * time.Second
	// half=5s; src=1.0 => 5s + 5s*1.0 = 10s
	got := jitter.Equal(base, fixed(1.0))
	if got != base {
		t.Fatalf("expected %v, got %v", base, got)
	}
	// src=0.0 => 5s + 0 = 5s
	got = jitter.Equal(base, fixed(0.0))
	if got != 5*time.Second {
		t.Fatalf("expected 5s, got %v", got)
	}
}
