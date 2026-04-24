package environ

import (
	"sort"
	"testing"
)

func TestFromSlice_Valid(t *testing.T) {
	m := fromSlice([]string{"FOO=bar", "BAZ=qux"})
	if m["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", m["FOO"])
	}
	if m["BAZ"] != "qux" {
		t.Errorf("expected BAZ=qux, got %q", m["BAZ"])
	}
}

func TestFromSlice_NoEquals(t *testing.T) {
	m := fromSlice([]string{"NOEQUALS"})
	if v, ok := m["NOEQUALS"]; !ok || v != "" {
		t.Errorf("expected empty value for NOEQUALS, got %q", v)
	}
}

func TestFromSlice_EmptyKey(t *testing.T) {
	m := fromSlice([]string{"=value"})
	if len(m) != 0 {
		t.Errorf("expected empty map for entry with empty key, got %v", m)
	}
}

func TestToSlice_RoundTrip(t *testing.T) {
	orig := map[string]string{"A": "1", "B": "2"}
	slice := ToSlice(orig)
	back := fromSlice(slice)
	for k, v := range orig {
		if back[k] != v {
			t.Errorf("round-trip mismatch for %s: want %q got %q", k, v, back[k])
		}
	}
}

func TestMerge_OverrideTakesPrecedence(t *testing.T) {
	base := map[string]string{"X": "base", "Y": "keep"}
	override := map[string]string{"X": "new"}
	result := Merge(base, override)
	if result["X"] != "new" {
		t.Errorf("expected X=new, got %q", result["X"])
	}
	if result["Y"] != "keep" {
		t.Errorf("expected Y=keep, got %q", result["Y"])
	}
}

func TestMerge_DoesNotMutateInputs(t *testing.T) {
	base := map[string]string{"A": "1"}
	override := map[string]string{"B": "2"}
	_ = Merge(base, override)
	if _, ok := base["B"]; ok {
		t.Error("Merge mutated base map")
	}
}

func TestDiff_DetectsChanges(t *testing.T) {
	old := map[string]string{"A": "1", "B": "2", "C": "3"}
	next := map[string]string{"A": "1", "B": "changed", "D": "4"}
	changed := Diff(old, next)
	sort.Strings(changed)
	// B changed, C removed, D added
	expected := []string{"B", "C", "D"}
	if len(changed) != len(expected) {
		t.Fatalf("expected %v, got %v", expected, changed)
	}
	for i, k := range expected {
		if changed[i] != k {
			t.Errorf("pos %d: expected %q got %q", i, k, changed[i])
		}
	}
}

func TestDiff_NoChanges(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2"}
	changed := Diff(env, env)
	if len(changed) != 0 {
		t.Errorf("expected no changes, got %v", changed)
	}
}
