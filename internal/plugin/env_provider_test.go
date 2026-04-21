package plugin_test

import (
	"strings"
	"testing"

	"vaultenv/internal/plugin"
)

func TestEnvProvider_Name(t *testing.T) {
	p := plugin.NewEnvProvider()
	if p.Name() != "env" {
		t.Errorf("expected name 'env', got %q", p.Name())
	}
}

func TestEnvProvider_GetSecret_Found(t *testing.T) {
	t.Setenv("VAULTENV_TEST_VAR", "supersecret")

	p := plugin.NewEnvProvider()
	val, err := p.GetSecret("ignored/path", "VAULTENV_TEST_VAR")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "supersecret" {
		t.Errorf("expected 'supersecret', got %q", val)
	}
}

func TestEnvProvider_GetSecret_Missing(t *testing.T) {
	p := plugin.NewEnvProvider()
	_, err := p.GetSecret("some/path", "VAULTENV_DEFINITELY_NOT_SET_XYZ")
	if err == nil {
		t.Fatal("expected error for missing variable, got nil")
	}
	if !strings.Contains(err.Error(), "not set") {
		t.Errorf("error message should mention 'not set', got: %v", err)
	}
}

func TestEnvProvider_CanBeRegistered(t *testing.T) {
	r := plugin.New()
	p := plugin.NewEnvProvider()

	if err := r.Register(p); err != nil {
		t.Fatalf("registration failed: %v", err)
	}

	got, err := r.Get("env")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if got.Name() != "env" {
		t.Errorf("expected 'env', got %q", got.Name())
	}
}
