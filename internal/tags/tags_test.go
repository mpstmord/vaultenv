package tags_test

import (
	"testing"

	"github.com/your-org/vaultenv/internal/tags"
)

func TestParse_Valid(t *testing.T) {
	pairs := []string{"env=prod", "region=us-east-1", "team=platform"}
	got, err := tags.Parse(pairs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["env"] != "prod" {
		t.Errorf("env: got %q, want %q", got["env"], "prod")
	}
	if got["region"] != "us-east-1" {
		t.Errorf("region: got %q, want %q", got["region"], "us-east-1")
	}
}

func TestParse_MissingEquals(t *testing.T) {
	_, err := tags.Parse([]string{"badpair"})
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
}

func TestParse_EmptyKey(t *testing.T) {
	_, err := tags.Parse([]string{"=value"})
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestParse_EmptySlice(t *testing.T) {
	got, err := tags.Parse(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty tags, got %v", got)
	}
}

func TestMatches_AllPresent(t *testing.T) {
	base, _ := tags.Parse([]string{"env=prod", "region=us-east-1"})
	sub, _ := tags.Parse([]string{"env=prod"})
	if !base.Matches(sub) {
		t.Error("expected Matches to return true")
	}
}

func TestMatches_ValueMismatch(t *testing.T) {
	base, _ := tags.Parse([]string{"env=prod"})
	sub, _ := tags.Parse([]string{"env=staging"})
	if base.Matches(sub) {
		t.Error("expected Matches to return false on value mismatch")
	}
}

func TestMatches_MissingKey(t *testing.T) {
	base, _ := tags.Parse([]string{"env=prod"})
	sub, _ := tags.Parse([]string{"env=prod", "region=us-east-1"})
	if base.Matches(sub) {
		t.Error("expected Matches to return false when key is absent")
	}
}

func TestKeys_Sorted(t *testing.T) {
	pairs := []string{"z=last", "a=first", "m=middle"}
	t2, _ := tags.Parse(pairs)
	keys := t2.Keys()
	if keys[0] != "a" || keys[1] != "m" || keys[2] != "z" {
		t.Errorf("keys not sorted: %v", keys)
	}
}

func TestString_Stable(t *testing.T) {
	pairs := []string{"env=prod", "region=us-east-1"}
	t1, _ := tags.Parse(pairs)
	got := t1.String()
	want := "env=prod,region=us-east-1"
	if got != want {
		t.Errorf("String: got %q, want %q", got, want)
	}
}

func TestGet_Found(t *testing.T) {
	tg, _ := tags.Parse([]string{"key=val"})
	v, ok := tg.Get("key")
	if !ok || v != "val" {
		t.Errorf("Get: got (%q, %v), want (\"val\", true)", v, ok)
	}
}

func TestGet_Missing(t *testing.T) {
	tg, _ := tags.Parse([]string{"key=val"})
	_, ok := tg.Get("missing")
	if ok {
		t.Error("expected Get to return false for missing key")
	}
}
