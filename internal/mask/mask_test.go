package mask

import (
	"testing"
)

func TestMask_NoSecrets(t *testing.T) {
	m := New()
	got := m.Mask("hello world")
	if got != "hello world" {
		t.Errorf("expected unchanged string, got %q", got)
	}
}

func TestMask_SingleSecret(t *testing.T) {
	m := New()
	m.Add("s3cr3t")
	got := m.Mask("password is s3cr3t, remember s3cr3t")
	expected := "password is ***REDACTED***, remember ***REDACTED***"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestMask_MultipleSecrets(t *testing.T) {
	m := New()
	m.Add("alpha")
	m.Add("beta")
	got := m.Mask("alpha and beta")
	if got == "alpha and beta" {
		t.Error("expected secrets to be redacted")
	}
	if got != "***REDACTED*** and ***REDACTED***" {
		t.Errorf("unexpected result: %q", got)
	}
}

func TestMask_EmptySecretIgnored(t *testing.T) {
	m := New()
	m.Add("")
	if len(m.secrets) != 0 {
		t.Error("empty secret should not be registered")
	}
}

func TestMaskEnv_RedactsSensitiveKeys(t *testing.T) {
	m := New()
	env := []string{"HOME=/root", "DB_PASS=supersecret", "USER=alice"}
	sensitive := map[string]struct{}{"DB_PASS": {}}
	got := m.MaskEnv(env, sensitive)
	for _, entry := range got {
		if entry == "DB_PASS=supersecret" {
			t.Error("sensitive value was not redacted")
		}
	}
	if got[0] != "HOME=/root" {
		t.Errorf("non-sensitive entry changed: %q", got[0])
	}
	if got[1] != "DB_PASS=***REDACTED***" {
		t.Errorf("expected redacted value, got %q", got[1])
	}
}

func TestMaskEnv_NoSensitiveKeys(t *testing.T) {
	m := New()
	env := []string{"FOO=bar", "BAZ=qux"}
	got := m.MaskEnv(env, map[string]struct{}{})
	for i, entry := range got {
		if entry != env[i] {
			t.Errorf("entry %d changed unexpectedly: %q", i, entry)
		}
	}
}
