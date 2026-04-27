package clock_test

import (
	"testing"
	"time"

	"github.com/your-org/vaultenv/internal/clock"
)

func TestReal_Now_IsRecent(t *testing.T) {
	c := clock.New()
	before := time.Now()
	now := c.Now()
	after := time.Now()

	if now.Before(before) || now.After(after) {
		t.Errorf("Real.Now() = %v, want between %v and %v", now, before, after)
	}
}

func TestReal_Since_IsNonNegative(t *testing.T) {
	c := clock.New()
	past := time.Now().Add(-time.Second)
	if c.Since(past) < 0 {
		t.Error("Real.Since returned negative duration")
	}
}

func TestFake_Now_ReturnsStart(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	f := clock.NewFake(start)

	if got := f.Now(); !got.Equal(start) {
		t.Errorf("Fake.Now() = %v, want %v", got, start)
	}
}

func TestFake_Advance_MovesTime(t *testing.T) {
	start := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	f := clock.NewFake(start)

	f.Advance(5 * time.Minute)

	want := start.Add(5 * time.Minute)
	if got := f.Now(); !got.Equal(want) {
		t.Errorf("after Advance: got %v, want %v", got, want)
	}
}

func TestFake_Set_SetsAbsoluteTime(t *testing.T) {
	f := clock.NewFake(time.Now())
	target := time.Date(2030, 6, 15, 9, 0, 0, 0, time.UTC)
	f.Set(target)

	if got := f.Now(); !got.Equal(target) {
		t.Errorf("Fake.Set: got %v, want %v", got, target)
	}
}

func TestFake_Since_ReflectsAdvance(t *testing.T) {
	start := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	f := clock.NewFake(start)

	f.Advance(10 * time.Second)

	got := f.Since(start)
	if got != 10*time.Second {
		t.Errorf("Fake.Since: got %v, want %v", got, 10*time.Second)
	}
}

func TestFake_After_AdvancesAndFires(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	f := clock.NewFake(start)

	ch := f.After(30 * time.Second)

	select {
	case fired := <-ch:
		want := start.Add(30 * time.Second)
		if !fired.Equal(want) {
			t.Errorf("After fired at %v, want %v", fired, want)
		}
	default:
		t.Error("After channel did not fire immediately on Fake")
	}

	if got := f.Now(); !got.Equal(start.Add(30*time.Second)) {
		t.Errorf("clock not advanced after After(): got %v", got)
	}
}
