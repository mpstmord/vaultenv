package expire_test

import (
	"context"
	"testing"
	"time"

	"github.com/your-org/vaultenv/internal/expire"
)

func TestSet_And_Get(t *testing.T) {
	tr := expire.New()
	val := map[string]any{"password": "s3cr3t"}
	tr.Set("secret/db", val, 10*time.Second)

	e, err := tr.Get("secret/db")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.Value["password"] != "s3cr3t" {
		t.Errorf("expected s3cr3t, got %v", e.Value["password"])
	}
}

func TestGet_NotFound(t *testing.T) {
	tr := expire.New()
	e, err := tr.Get("secret/missing")
	if err != nil {
		t.Fatalf("unexpected error for missing key: %v", err)
	}
	if e.Value != nil {
		t.Errorf("expected nil value for missing key")
	}
}

func TestGet_Expired(t *testing.T) {
	tr := expire.New()
	tr.Set("secret/old", map[string]any{"k": "v"}, -1*time.Second)

	_, err := tr.Get("secret/old")
	if err != expire.ErrExpired {
		t.Fatalf("expected ErrExpired, got %v", err)
	}
}

func TestEntry_TTL_Positive(t *testing.T) {
	tr := expire.New()
	tr.Set("secret/ttl", map[string]any{}, 5*time.Second)

	e, err := tr.Get("secret/ttl")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.TTL() <= 0 {
		t.Errorf("expected positive TTL, got %v", e.TTL())
	}
}

func TestDelete_RemovesEntry(t *testing.T) {
	tr := expire.New()
	tr.Set("secret/del", map[string]any{"x": "y"}, 10*time.Second)
	tr.Delete("secret/del")

	e, err := tr.Get("secret/del")
	if err != nil {
		t.Fatalf("unexpected error after delete: %v", err)
	}
	if e.Value != nil {
		t.Errorf("expected nil value after delete")
	}
}

func TestPurge_RemovesExpired(t *testing.T) {
	tr := expire.New()
	tr.Set("secret/expired", map[string]any{}, -1*time.Second)
	tr.Set("secret/live", map[string]any{}, 10*time.Second)

	n := tr.Purge(context.Background())
	if n != 1 {
		t.Errorf("expected 1 purged entry, got %d", n)
	}

	_, err := tr.Get("secret/live")
	if err != nil {
		t.Errorf("live entry should not be purged: %v", err)
	}
}
