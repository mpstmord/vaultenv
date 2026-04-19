package env

import (
	"os"
	"testing"
)

func TestParseMapping_Valid(t *testing.T) {
	m, err := ParseMapping("DB_PASSWORD=secret/db#password")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.EnvVar != "DB_PASSWORD" {
		t.Errorf("expected DB_PASSWORD, got %s", m.EnvVar)
	}
	if m.Path != "secret/db" {
		t.Errorf("expected secret/db, got %s", m.Path)
	}
	if m.Field != "password" {
		t.Errorf("expected password, got %s", m.Field)
	}
}

func TestParseMapping_MissingEquals(t *testing.T) {
	_, err := ParseMapping("DB_PASSWORDsecret/db#password")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestParseMapping_MissingHash(t *testing.T) {
	_, err := ParseMapping("DB_PASSWORD=secret/dbpassword")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestParseMapping_EmptyField(t *testing.T) {
	_, err := ParseMapping("DB_PASSWORD=secret/db#")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestInjectIntoEnv(t *testing.T) {
	resolved := map[string]string{
		"TEST_VAR_ONE": "hello",
		"TEST_VAR_TWO": "world",
	}
	if err := InjectIntoEnv(resolved); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v := os.Getenv("TEST_VAR_ONE"); v != "hello" {
		t.Errorf("expected hello, got %s", v)
	}
	if v := os.Getenv("TEST_VAR_TWO"); v != "world" {
		t.Errorf("expected world, got %s", v)
	}
}
