package circuit

import (
	"testing"
	"time"
)

func TestNew_InvalidThreshold(t *testing.T) {
	_, err := New(0, time.Second)
	if err == nil {
		t.Fatal("expected error for zero threshold")
	}
}

func TestNew_InvalidResetTimeout(t *testing.T) {
	_, err := New(3, 0)
	if err == nil {
		t.Fatal("expected error for zero reset timeout")
	}
}

func TestNew_Valid(t *testing.T) {
	b, err := New(3, time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.State() != StateClosed {
		t.Errorf("expected closed state, got %s", b.State())
	}
}

func TestBreaker_OpensAfterThreshold(t *testing.T) {
	b, _ := New(3, time.Second)
	for i := 0; i < 3; i++ {
		b.RecordFailure()
	}
	if b.State() != StateOpen {
		t.Errorf("expected open state after threshold, got %s", b.State())
	}
	if err := b.Allow(); err != ErrOpen {
		t.Errorf("expected ErrOpen, got %v", err)
	}
}

func TestBreaker_SuccessResets(t *testing.T) {
	b, _ := New(2, time.Second)
	b.RecordFailure()
	b.RecordSuccess()
	if b.State() != StateClosed {
		t.Errorf("expected closed after success, got %s", b.State())
	}
	if err := b.Allow(); err != nil {
		t.Errorf("expected nil after reset, got %v", err)
	}
}

func TestBreaker_HalfOpenAfterTimeout(t *testing.T) {
	b, _ := New(1, 10*time.Millisecond)
	b.RecordFailure()
	time.Sleep(20 * time.Millisecond)
	if err := b.Allow(); err != nil {
		t.Errorf("expected nil in half-open, got %v", err)
	}
	if b.State() != StateHalfOpen {
		t.Errorf("expected half-open state, got %s", b.State())
	}
}

func TestState_String(t *testing.T) {
	cases := []struct {
		s    State
		want string
	}{
		{StateClosed, "closed"},
		{StateHalfOpen, "half-open"},
		{StateOpen, "open"},
		{State(99), "unknown"},
	}
	for _, c := range cases {
		if got := c.s.String(); got != c.want {
			t.Errorf("State(%d).String() = %q, want %q", c.s, got, c.want)
		}
	}
}
