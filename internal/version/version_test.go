package version_test

import (
	"strings"
	"testing"

	"github.com/your-org/vaultenv/internal/version"
)

func TestGet_ReturnsInfo(t *testing.T) {
	info := version.Get()
	if info.Major == "" {
		t.Error("expected non-empty Major")
	}
	if info.Minor == "" {
		t.Error("expected non-empty Minor")
	}
	if info.Patch == "" {
		t.Error("expected non-empty Patch")
	}
}

func TestInfo_String_ContainsSemver(t *testing.T) {
	info := version.Info{
		Major:     "1",
		Minor:     "2",
		Patch:     "3",
		Commit:    "abc1234",
		BuildDate: "2024-06-01",
	}
	s := info.String()
	if !strings.Contains(s, "v1.2.3") {
		t.Errorf("expected semver in string, got: %s", s)
	}
	if !strings.Contains(s, "abc1234") {
		t.Errorf("expected commit hash in string, got: %s", s)
	}
	if !strings.Contains(s, "2024-06-01") {
		t.Errorf("expected build date in string, got: %s", s)
	}
}

func TestInfo_Short_OnlySemver(t *testing.T) {
	info := version.Info{
		Major:  "2",
		Minor:  "0",
		Patch:  "1",
		Commit: "deadbeef",
	}
	short := info.Short()
	if short != "v2.0.1" {
		t.Errorf("expected v2.0.1, got: %s", short)
	}
	if strings.Contains(short, "deadbeef") {
		t.Errorf("Short() should not contain commit hash, got: %s", short)
	}
}

func TestInfo_String_DefaultValues(t *testing.T) {
	info := version.Get()
	s := info.String()
	// Should always produce a non-empty, parseable string.
	if !strings.HasPrefix(s, "v") {
		t.Errorf("expected string to start with 'v', got: %s", s)
	}
}
