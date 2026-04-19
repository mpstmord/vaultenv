package vault

import (
	"testing"
)

func TestParseSecretPath_Valid(t *testing.T) {
	cases := []struct {
		input string
		mont  string
		path  string
	}{
		{"secret/myapp/db", "secret", "myapp/db"},
		{"/secret/myapp/db", "secret", "myapp/db"},
		{"kv/service/token", "kv", "service/token"},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			sp, err := ParseSecretPath(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if sp.Mount != tc.mont {
				t.Errorf("mount: got %q, want %q", sp.Mount, tc.mont)
			}
			if sp.Path != tc.path {
				t.Errorf("path: got %q, want %q", sp.Path, tc.path)
			}
		})
	}
}

func TestParseSecretPath_Invalid(t *testing.T) {
	cases := []string{
		"",
		"noseparator",
		"/missingpath/",
		"mount/",
	}
	for _, tc := range cases {
		t.Run(tc, func(t *testing.T) {
			_, err := ParseSecretPath(tc)
			if err == nil {
				t.Errorf("expected error for input %q, got nil", tc)
			}
		})
	}
}
