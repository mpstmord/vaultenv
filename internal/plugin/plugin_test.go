package plugin_test

import (
	"errors"
	"testing"

	"vaultenv/internal/plugin"
)

// stubProvider is a minimal Provider for testing.
type stubProvider struct{ name string }

func (s *stubProvider) Name() string { return s.name }
func (s *stubProvider) GetSecret(path, field string) (string, error) {
	return "stub-value", nil
}

func TestRegister_And_Get(t *testing.T) {
	r := plugin.New()
	p := &stubProvider{name: "env"}

	if err := r.Register(p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := r.Get("env")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if got.Name() != "env" {
		t.Errorf("expected name 'env', got %q", got.Name())
	}
}

func TestRegister_AlreadyRegistered(t *testing.T) {
	r := plugin.New()
	p := &stubProvider{name: "dup"}

	_ = r.Register(p)
	err := r.Register(p)
	if !errors.Is(err, plugin.ErrAlreadyRegistered) {
		t.Errorf("expected ErrAlreadyRegistered, got %v", err)
	}
}

func TestGet_NotFound(t *testing.T) {
	r := plugin.New()
	_, err := r.Get("missing")
	if !errors.Is(err, plugin.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestNames_ReturnsAll(t *testing.T) {
	r := plugin.New()
	_ = r.Register(&stubProvider{name: "a"})
	_ = r.Register(&stubProvider{name: "b"})

	names := r.Names()
	if len(names) != 2 {
		t.Errorf("expected 2 names, got %d", len(names))
	}
}

func TestUnregister_Success(t *testing.T) {
	r := plugin.New()
	_ = r.Register(&stubProvider{name: "x"})

	if err := r.Unregister("x"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := r.Get("x"); !errors.Is(err, plugin.ErrNotFound) {
		t.Errorf("expected ErrNotFound after unregister, got %v", err)
	}
}

func TestUnregister_NotFound(t *testing.T) {
	r := plugin.New()
	err := r.Unregister("ghost")
	if !errors.Is(err, plugin.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}
