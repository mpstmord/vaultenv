package diff_test

import (
	"testing"

	"github.com/your-org/vaultenv/internal/diff"
)

func TestCompare_NoChanges(t *testing.T) {
	c := diff.New()
	prev := map[string]string{"A": "1", "B": "2"}
	next := map[string]string{"A": "1", "B": "2"}
	changes := c.Compare(prev, next)
	if len(changes) != 0 {
		t.Fatalf("expected 0 changes, got %d", len(changes))
	}
}

func TestCompare_Added(t *testing.T) {
	c := diff.New()
	prev := map[string]string{}
	next := map[string]string{"X": "hello"}
	changes := c.Compare(prev, next)
	if len(changes) != 1 || changes[0].Type != diff.Added || changes[0].Key != "X" {
		t.Fatalf("expected Added X, got %+v", changes)
	}
}

func TestCompare_Removed(t *testing.T) {
	c := diff.New()
	prev := map[string]string{"Y": "val"}
	next := map[string]string{}
	changes := c.Compare(prev, next)
	if len(changes) != 1 || changes[0].Type != diff.Removed || changes[0].Key != "Y" {
		t.Fatalf("expected Removed Y, got %+v", changes)
	}
}

func TestCompare_Changed(t *testing.T) {
	c := diff.New()
	prev := map[string]string{"Z": "old"}
	next := map[string]string{"Z": "new"}
	changes := c.Compare(prev, next)
	if len(changes) != 1 || changes[0].Type != diff.Changed || changes[0].Key != "Z" {
		t.Fatalf("expected Changed Z, got %+v", changes)
	}
}

func TestHasChanges_True(t *testing.T) {
	c := diff.New()
	if !c.HasChanges(map[string]string{"a": "1"}, map[string]string{"a": "2"}) {
		t.Fatal("expected HasChanges to be true")
	}
}

func TestHasChanges_False(t *testing.T) {
	c := diff.New()
	if c.HasChanges(map[string]string{"a": "1"}, map[string]string{"a": "1"}) {
		t.Fatal("expected HasChanges to be false")
	}
}
