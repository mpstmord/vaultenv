package warmup_test

import (
	"testing"

	"github.com/your-org/vaultenv/internal/warmup"
)

func TestParseMappings_Valid(t *testing.T) {
	specs := []string{
		"DB_PASSWORD=secret/db#password",
		"API_KEY=secret/svc/api#key",
	}
	mappings, err := warmup.ParseMappings(specs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(mappings) != 2 {
		t.Fatalf("expected 2 mappings, got %d", len(mappings))
	}

	if mappings[0].EnvKey != "DB_PASSWORD" {
		t.Errorf("expected EnvKey=DB_PASSWORD, got %q", mappings[0].EnvKey)
	}
	if mappings[0].Path != "secret/db" {
		t.Errorf("expected Path=secret/db, got %q", mappings[0].Path)
	}
	if mappings[0].Field != "password" {
		t.Errorf("expected Field=password, got %q", mappings[0].Field)
	}
}

func TestParseMappings_MissingEquals(t *testing.T) {
	_, err := warmup.ParseMappings([]string{"DB_PASSWORDsecret/db#password"})
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
}

func TestParseMappings_MissingHash(t *testing.T) {
	_, err := warmup.ParseMappings([]string{"DB_PASSWORD=secret/dbpassword"})
	if err == nil {
		t.Fatal("expected error for missing '#'")
	}
}

func TestParseMappings_EmptyField(t *testing.T) {
	_, err := warmup.ParseMappings([]string{"DB_PASSWORD=secret/db#"})
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestParseMappings_EmptyPath(t *testing.T) {
	_, err := warmup.ParseMappings([]string{"DB_PASSWORD=#password"})
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestParseMappings_EmptySlice(t *testing.T) {
	mappings, err := warmup.ParseMappings(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(mappings) != 0 {
		t.Errorf("expected 0 mappings, got %d", len(mappings))
	}
}
