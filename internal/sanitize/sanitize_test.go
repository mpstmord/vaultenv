package sanitize_test

import (
	"testing"

	"github.com/your-org/vaultenv/internal/sanitize"
)

func TestEnvKey_Valid(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"FOO", "FOO"},
		{"foo", "FOO"},
		{"  bar  ", "BAR"},
		{"My_Var_1", "MY_VAR_1"},
		{"_PRIVATE", "_PRIVATE"},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got, err := sanitize.EnvKey(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("EnvKey(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestEnvKey_Invalid(t *testing.T) {
	cases := []string{
		"",
		"   ",
		"1STARTS_WITH_DIGIT",
		"HAS-HYPHEN",
		"HAS SPACE",
		"HAS.DOT",
	}
	for _, input := range cases {
		t.Run(input, func(t *testing.T) {
			_, err := sanitize.EnvKey(input)
			if err == nil {
				t.Errorf("EnvKey(%q): expected error, got nil", input)
			}
		})
	}
}

func TestSecretValue_Valid(t *testing.T) {
	got, err := sanitize.SecretValue("  my-secret  ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "my-secret" {
		t.Errorf("got %q, want %q", got, "my-secret")
	}
}

func TestSecretValue_Empty(t *testing.T) {
	for _, v := range []string{"", "   "} {
		_, err := sanitize.SecretValue(v)
		if err == nil {
			t.Errorf("SecretValue(%q): expected error, got nil", v)
		}
	}
}

func TestEnvPair_Valid(t *testing.T) {
	k, v, err := sanitize.EnvPair(" db_pass ", " s3cr3t ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if k != "DB_PASS" {
		t.Errorf("key = %q, want DB_PASS", k)
	}
	if v != "s3cr3t" {
		t.Errorf("value = %q, want s3cr3t", v)
	}
}

func TestEnvPair_BadKey(t *testing.T) {
	_, _, err := sanitize.EnvPair("bad-key", "value")
	if err == nil {
		t.Error("expected error for bad key, got nil")
	}
}

func TestEnvPair_EmptyValue(t *testing.T) {
	_, _, err := sanitize.EnvPair("GOOD_KEY", "")
	if err == nil {
		t.Error("expected error for empty value, got nil")
	}
}
