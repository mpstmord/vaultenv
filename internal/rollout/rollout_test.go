package rollout

import (
	"testing"
)

func TestNew_InvalidPercentage(t *testing.T) {
	tests := []struct {
		name string
		pct  int
	}{
		{"negative", -1},
		{"over 100", 101},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := New(Config{Percentage: tc.pct})
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

func TestNew_Valid(t *testing.T) {
	g, err := New(Config{Percentage: 50})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if g == nil {
		t.Fatal("expected non-nil gate")
	}
}

func TestAllow_ZeroPercent(t *testing.T) {
	g, _ := New(Config{Percentage: 0})
	for _, key := range []string{"a", "b", "secret/db", "secret/api"} {
		if g.Allow(key) {
			t.Errorf("Allow(%q) = true, want false at 0%%", key)
		}
	}
}

func TestAllow_FullPercent(t *testing.T) {
	g, _ := New(Config{Percentage: 100})
	for _, key := range []string{"a", "b", "secret/db", "secret/api"} {
		if !g.Allow(key) {
			t.Errorf("Allow(%q) = false, want true at 100%%", key)
		}
	}
}

func TestAllow_HashIsDeterministic(t *testing.T) {
	g, _ := New(Config{Percentage: 50, Strategy: StrategyHash})
	key := "secret/payments/api-key"
	first := g.Allow(key)
	for i := 0; i < 10; i++ {
		if g.Allow(key) != first {
			t.Fatalf("Allow(%q) returned different results across calls", key)
		}
	}
}

func TestAllow_DistributionApproximate(t *testing.T) {
	g, _ := New(Config{Percentage: 50, Strategy: StrategyHash})
	allowed := 0
	total := 1000
	for i := 0; i < total; i++ {
		key := string(rune('a' + i%26))
		// Build varied keys.
		key = key + string(rune(i))
		if g.Allow(key) {
			allowed++
		}
	}
	// Expect roughly 50% ± 15%.
	if allowed < 350 || allowed > 650 {
		t.Errorf("distribution out of range: %d/1000 allowed", allowed)
	}
}

func TestSetPercentage_UpdatesGate(t *testing.T) {
	g, _ := New(Config{Percentage: 0})
	if err := g.SetPercentage(100); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !g.Allow("any-key") {
		t.Error("expected Allow to return true after setting 100%")
	}
}

func TestSetPercentage_Invalid(t *testing.T) {
	g, _ := New(Config{Percentage: 50})
	if err := g.SetPercentage(150); err == nil {
		t.Error("expected error for invalid percentage")
	}
}
