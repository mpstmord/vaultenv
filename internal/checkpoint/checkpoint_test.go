package checkpoint_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultenv/internal/checkpoint"
)

func tempFile(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "checkpoint.json")
}

func TestNew_MissingFile(t *testing.T) {
	s, err := checkpoint.New(tempFile(t))
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil store")
	}
}

func TestSet_And_Get(t *testing.T) {
	s, _ := checkpoint.New(tempFile(t))
	if err := s.Set("secret/foo", 3); err != nil {
		t.Fatalf("Set: %v", err)
	}
	r, ok := s.Get("secret/foo")
	if !ok {
		t.Fatal("expected record to exist")
	}
	if r.Version != 3 {
		t.Errorf("expected version 3, got %d", r.Version)
	}
	if r.Path != "secret/foo" {
		t.Errorf("unexpected path %q", r.Path)
	}
	if r.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}
}

func TestGet_Missing(t *testing.T) {
	s, _ := checkpoint.New(tempFile(t))
	_, ok := s.Get("secret/nope")
	if ok {
		t.Fatal("expected record to be absent")
	}
}

func TestDelete_RemovesRecord(t *testing.T) {
	s, _ := checkpoint.New(tempFile(t))
	_ = s.Set("secret/bar", 1)
	if err := s.Delete("secret/bar"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, ok := s.Get("secret/bar")
	if ok {
		t.Fatal("expected record to be deleted")
	}
}

func TestPersistence_ReloadFromDisk(t *testing.T) {
	f := tempFile(t)
	s1, _ := checkpoint.New(f)
	_ = s1.Set("secret/baz", 7)

	s2, err := checkpoint.New(f)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	r, ok := s2.Get("secret/baz")
	if !ok {
		t.Fatal("expected record after reload")
	}
	if r.Version != 7 {
		t.Errorf("expected version 7, got %d", r.Version)
	}
}

func TestNew_CorruptFile(t *testing.T) {
	f := tempFile(t)
	_ = os.WriteFile(f, []byte("not json{"), 0o600)
	_, err := checkpoint.New(f)
	if err == nil {
		t.Fatal("expected error for corrupt file")
	}
}
