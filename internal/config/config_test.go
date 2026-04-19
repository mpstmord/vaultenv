package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultenv/internal/config"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	return p
}

func TestLoadFile_Valid(t *testing.T) {
	p := writeTemp(t, `
vault:
  address: http://127.0.0.1:8200
  token: root
secrets:
  - mapping: DB_PASS=secret/db#password
`)
	cfg, err := config.LoadFile(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Vault.Address != "http://127.0.0.1:8200" {
		t.Errorf("address mismatch: %q", cfg.Vault.Address)
	}
	if len(cfg.Secrets) != 1 {
		t.Errorf("expected 1 secret, got %d", len(cfg.Secrets))
	}
}

func TestLoadFile_Missing(t *testing.T) {
	_, err := config.LoadFile("/nonexistent/path.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadFile_UnknownField(t *testing.T) {
	p := writeTemp(t, `vault:\n  address: http://localhost\n  token: t\n  unknown: bad\n`)
	_, err := config.LoadFile(p)
	if err == nil {
		t.Fatal("expected error for unknown field")
	}
}

func TestValidate_MissingAddress(t *testing.T) {
	cfg := &config.Config{
		Vault:   config.VaultConfig{Token: "t"},
		Secrets: []config.SecretEntry{{Mapping: "X=p#f"}},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for missing address")
	}
}

func TestValidate_NoSecrets(t *testing.T) {
	cfg := &config.Config{
		Vault: config.VaultConfig{Address: "http://x", Token: "t"},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for empty secrets")
	}
}

func TestValidate_OK(t *testing.T) {
	cfg := &config.Config{
		Vault:   config.VaultConfig{Address: "http://x", Token: "t"},
		Secrets: []config.SecretEntry{{Mapping: "X=p#f"}},
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
