package process_test

import (
	"testing"

	"github.com/your-org/vaultenv/internal/process"
)

func TestRun_EmptyCommand(t *testing.T) {
	r := process.NewRunner(nil)
	err := r.Run("", nil)
	if err == nil {
		t.Fatal("expected error for empty command")
	}
}

func TestRun_UnknownCommand(t *testing.T) {
	r := process.NewRunner(nil)
	err := r.Run("__no_such_binary_xyz__", nil)
	if err == nil {
		t.Fatal("expected error for unknown command")
	}
}

func TestRun_Success(t *testing.T) {
	r := process.NewRunner([]string{"PATH=/usr/bin:/bin"})
	err := r.Run("true", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMergeEnv_Override(t *testing.T) {
	base := []string{"FOO=old", "BAR=keep"}
	overrides := map[string]string{"FOO": "new"}

	result := process.MergeEnv(base, overrides)

	found := false
	for _, e := range result {
		if e == "FOO=new" {
			found = true
		}
		if e == "FOO=old" {
			t.Fatal("old value should have been replaced")
		}
	}
	if !found {
		t.Fatal("expected FOO=new in result")
	}
}

func TestMergeEnv_Append(t *testing.T) {
	base := []string{"EXISTING=yes"}
	overrides := map[string]string{"NEW_KEY": "value"}

	result := process.MergeEnv(base, overrides)

	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
}

func TestMergeEnv_Empty(t *testing.T) {
	result := process.MergeEnv(nil, map[string]string{"A": "1"})
	if len(result) != 1 || result[0] != "A=1" {
		t.Fatalf("unexpected result: %v", result)
	}
}
