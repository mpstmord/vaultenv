package backoff_test

import (
	"testing"
	"time"

	"github.com/your-org/vaultenv/internal/backoff"
)

func TestDefault_ReturnsNonZeroValues(t *testing.T) {
	s := backoff.Default()
	if s.Base == 0 {
		t.Error("expected non-zero Base")
	}
	if s.Max == 0 {
		t.Error("expected non-zero Max")
	}
	if s.Factor == 0 {
		t.Error("expected non-zero Factor")
	}
}

func TestDelay_ZeroAttempt(t *testing.T) {
	s := backoff.Strategy{
		Base:   100 * time.Millisecond,
		Max:    10 * time.Second,
		Factor: 2.0,
		Jitter: false,
	}
	d := s.Delay(0)
	if d != 100*time.Millisecond {
		t.Errorf("expected 100ms, got %v", d)
	}
}

func TestDelay_Increases(t *testing.T) {
	s := backoff.Strategy{
		Base:   100 * time.Millisecond,
		Max:    10 * time.Second,
		Factor: 2.0,
		Jitter: false,
	}
	prev := s.Delay(0)
	for i := 1; i <= 4; i++ {
		curr := s.Delay(i)
		if curr <= prev {
			t.Errorf("attempt %d: expected delay > %v, got %v", i, prev, curr)
		}
		prev = curr
	}
}

func TestDelay_CappedAtMax(t *testing.T) {
	s := backoff.Strategy{
		Base:   1 * time.Second,
		Max:    2 * time.Second,
		Factor: 10.0,
		Jitter: false,
	}
	for i := 0; i < 10; i++ {
		d := s.Delay(i)
		if d > s.Max {
			t.Errorf("attempt %d: delay %v exceeds max %v", i, d, s.Max)
		}
	}
}

func TestDelay_NegativeAttemptClamped(t *testing.T) {
	s := backoff.Strategy{
		Base:   50 * time.Millisecond,
		Max:    5 * time.Second,
		Factor: 2.0,
		Jitter: false,
	}
	d := s.Delay(-3)
	if d != 50*time.Millisecond {
		t.Errorf("expected base delay for negative attempt, got %v", d)
	}
}

func TestAttempts_Length(t *testing.T) {
	s := backoff.Default()
	delays := s.Attempts(5)
	if len(delays) != 5 {
		t.Errorf("expected 5 delays, got %d", len(delays))
	}
}

func TestDelay_WithJitter_WithinBounds(t *testing.T) {
	s := backoff.Strategy{
		Base:   100 * time.Millisecond,
		Max:    10 * time.Second,
		Factor: 2.0,
		Jitter: true,
	}
	// With jitter the delay should still be positive and <= Max.
	for i := 0; i < 20; i++ {
		d := s.Delay(2)
		if d <= 0 {
			t.Errorf("expected positive delay, got %v", d)
		}
		if d > s.Max {
			t.Errorf("delay %v exceeds max %v", d, s.Max)
		}
	}
}
