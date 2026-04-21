package store

import (
	"testing"
)

func TestSet_And_Get(t *testing.T) {
	s := New()
	s.Set("FOO", "bar")
	v, err := s.Get("FOO")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != "bar" {
		t.Errorf("expected bar, got %s", v)
	}
}

func TestGet_NotFound(t *testing.T) {
	s := New()
	_, err := s.Get("MISSING")
	if err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestSet_Overwrite(t *testing.T) {
	s := New()
	s.Set("KEY", "first")
	s.Set("KEY", "second")
	v, _ := s.Get("KEY")
	if v != "second" {
		t.Errorf("expected second, got %s", v)
	}
}

func TestDelete_RemovesKey(t *testing.T) {
	s := New()
	s.Set("DEL", "val")
	s.Delete("DEL")
	_, err := s.Get("DEL")
	if err != ErrNotFound {
		t.Fatalf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestDelete_Noop(t *testing.T) {
	s := New()
	// Should not panic on missing key
	s.Delete("NONEXISTENT")
}

func TestLen(t *testing.T) {
	s := New()
	if s.Len() != 0 {
		t.Fatalf("expected 0, got %d", s.Len())
	}
	s.Set("A", "1")
	s.Set("B", "2")
	if s.Len() != 2 {
		t.Fatalf("expected 2, got %d", s.Len())
	}
}

func TestKeys_ReturnsAll(t *testing.T) {
	s := New()
	s.Set("X", "1")
	s.Set("Y", "2")
	keys := s.Keys()
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(keys))
	}
}

func TestSnapshot_IsCopy(t *testing.T) {
	s := New()
	s.Set("K", "v")
	snap := s.Snapshot()
	snap["K"] = "mutated"
	v, _ := s.Get("K")
	if v != "v" {
		t.Error("snapshot mutation affected original store")
	}
}
