package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunCmd_MissingMapping(t *testing.T) {
	rootCmd.SetArgs([]string{"run", "--", "env"})
	var buf bytes.Buffer
	rootCmd.SetErr(&buf)

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error when --mapping flag is missing")
	}
}

func TestRunCmd_InvalidMapping(t *testing.T) {
	rootCmd.SetArgs([]string{
		"run",
		"--mapping", "INVALID_NO_EQUALS_OR_HASH",
		"--vault-addr", "http://127.0.0.1:8200",
		"--vault-token", "test-token",
		"--", "env",
	})
	var buf bytes.Buffer
	rootCmd.SetErr(&buf)

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for invalid mapping")
	}
	if !strings.Contains(err.Error(), "parse mapping") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestRunCmd_UsageContainsMapping(t *testing.T) {
	usage := runCmd.UsageString()
	if !strings.Contains(usage, "--mapping") {
		t.Error("usage should mention --mapping flag")
	}
}
