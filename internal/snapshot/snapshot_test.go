package snapshot_test

import (
	"testing"

	"github.com/your-org/vaultenv/internal/snapshot"
)

func TestNew_CopiesData(t *testing.T) {
	original := map[string]string{"A": "1", "B": "2"}
	s := snapshot.New(original)
	original["A"] = "mutated"
	if v, _ := s.Get("A"); v != "1" {
		t.Fatalf("expected '1', got %q", v)
	}
}

func TestGet_MissingKey(t *testing.T) {
	s := snapshot.New(map[string]string{})
	_, ok := s.Get("missing")
	if ok {
		t.Fatal("expected missing key to return false")
	}
}

func TestKeys_Sorted(t *testing.T) {
	s := snapshot.New(map[string]string{"Z": "z", "A": "a", "M": "m"})
	keys := s.Keys()
	expected := []string{"A", "M", "Z"}
	for i, k := range keys {
		if k != expected[i] {
			t.Fatalf("keys[%d] = %q, want %q", i, k, expected[i])
		}
	}
}

func TestDigest_Stable(t *testing.T) {
	a := snapshot.New(map[string]string{"X": "1", "Y": "2"})
	b := snapshot.New(map[string]string{"Y": "2", "X": "1"})
	if a.Digest() != b.Digest() {
		t.Fatal("digests should be equal for same content regardless of insertion order")
	}
}

func TestEqual_SameContent(t *testing.T) {
	a := snapshot.New(map[string]string{"K": "v"})
	b := snapshot.New(map[string]string{"K": "v"})
	if !a.Equal(b) {
		t.Fatal("expected snapshots to be equal")
	}
}

func TestEqual_DifferentContent(t *testing.T) {
	a := snapshot.New(map[string]string{"K": "v1"})
	b := snapshot.New(map[string]string{"K": "v2"})
	if a.Equal(b) {
		t.Fatal("expected snapshots to differ")
	}
}

func TestDiff_DetectsChanges(t *testing.T) {
	old := snapshot.New(map[string]string{"A": "1", "B": "2", "C": "3"})
	new := snapshot.New(map[string]string{"A": "1", "B": "changed", "D": "4"})
	diff := old.Diff(new)
	// B changed, C removed, D added
	expected := map[string]bool{"B": true, "C": true, "D": true}
	if len(diff) != 3 {
		t.Fatalf("expected 3 diffs, got %d: %v", len(diff), diff)
	}
	for _, k := range diff {
		if !expected[k] {
			t.Errorf("unexpected diff key: %q", k)
		}
	}
}
