package metrics

import (
	"testing"
	"time"
)

func TestCounter_IncAndGet(t *testing.T) {
	var c Counter
	if c.Get() != 0 {
		t.Fatalf("expected 0, got %d", c.Get())
	}
	c.Inc()
	c.Inc()
	if c.Get() != 2 {
		t.Fatalf("expected 2, got %d", c.Get())
	}
}

func TestNew_Uptime(t *testing.T) {
	col := New()
	time.Sleep(2 * time.Millisecond)
	if col.Uptime() < time.Millisecond {
		t.Fatal("expected uptime >= 1ms")
	}
}

func TestSnapshot_AllKeys(t *testing.T) {
	col := New()
	col.SecretsResolved.Inc()
	col.CacheHits.Inc()
	col.CacheHits.Inc()
	col.Errors.Inc()

	snap := col.Snapshot()
	expected := map[string]uint64{
		"secrets_resolved": 1,
		"cache_hits":       2,
		"cache_misses":     0,
		"renewals":         0,
		"errors":           1,
	}
	for k, want := range expected {
		if got := snap[k]; got != want {
			t.Errorf("key %q: want %d, got %d", k, want, got)
		}
	}
}

func TestSnapshot_IndependentCopy(t *testing.T) {
	col := New()
	snap1 := col.Snapshot()
	col.CacheMisses.Inc()
	snap2 := col.Snapshot()
	if snap1["cache_misses"] != 0 {
		t.Error("snap1 should not reflect later increments")
	}
	if snap2["cache_misses"] != 1 {
		t.Error("snap2 should reflect the increment")
	}
}
