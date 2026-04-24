package labels_test

import (
	"strings"
	"testing"

	"github.com/yourusername/vaultenv/internal/labels"
)

func TestParse_Valid(t *testing.T) {
	l, err := labels.Parse([]string{"env=prod", "region=us-east-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, ok := l.Get("env"); !ok || v != "prod" {
		t.Errorf("expected env=prod, got %q ok=%v", v, ok)
	}
	if v, ok := l.Get("region"); !ok || v != "us-east-1" {
		t.Errorf("expected region=us-east-1, got %q ok=%v", v, ok)
	}
}

func TestParse_MissingEquals(t *testing.T) {
	_, err := labels.Parse([]string{"noequalssign"})
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
}

func TestParse_EmptyKey(t *testing.T) {
	_, err := labels.Parse([]string{"=value"})
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestParse_EmptySlice(t *testing.T) {
	l, err := labels.Parse(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(l) != 0 {
		t.Errorf("expected empty labels, got %v", l)
	}
}

func TestMatches_AllPresent(t *testing.T) {
	l, _ := labels.Parse([]string{"env=prod", "team=platform"})
	filter, _ := labels.Parse([]string{"env=prod"})
	if !l.Matches(filter) {
		t.Error("expected match")
	}
}

func TestMatches_WrongValue(t *testing.T) {
	l, _ := labels.Parse([]string{"env=prod"})
	filter, _ := labels.Parse([]string{"env=staging"})
	if l.Matches(filter) {
		t.Error("expected no match")
	}
}

func TestMatches_EmptyFilterAlwaysMatches(t *testing.T) {
	l, _ := labels.Parse([]string{"env=prod"})
	if !l.Matches(labels.Labels{}) {
		t.Error("empty filter should always match")
	}
}

func TestMerge_OtherTakesPrecedence(t *testing.T) {
	base, _ := labels.Parse([]string{"env=prod", "owner=alice"})
	override, _ := labels.Parse([]string{"env=staging"})
	merged := base.Merge(override)
	if v, _ := merged.Get("env"); v != "staging" {
		t.Errorf("expected env=staging, got %q", v)
	}
	if v, _ := merged.Get("owner"); v != "alice" {
		t.Errorf("expected owner=alice, got %q", v)
	}
}

func TestString_ContainsPairs(t *testing.T) {
	l, _ := labels.Parse([]string{"env=prod", "region=eu"})
	s := l.String()
	if !strings.Contains(s, "env=prod") {
		t.Errorf("String() missing env=prod: %q", s)
	}
	if !strings.Contains(s, "region=eu") {
		t.Errorf("String() missing region=eu: %q", s)
	}
}

func TestString_EmptyLabels(t *testing.T) {
	var l labels.Labels
	if s := l.String(); s != "" {
		t.Errorf("expected empty string, got %q", s)
	}
}
