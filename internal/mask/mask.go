// Package mask provides utilities for redacting sensitive values
// before they are written to logs or displayed in output.
package mask

import "strings"

const redacted = "***REDACTED***"

// Masker holds a set of sensitive values and can redact them from strings.
type Masker struct {
	secrets map[string]struct{}
}

// New returns a new Masker with no registered secrets.
func New() *Masker {
	return &Masker{secrets: make(map[string]struct{})}
}

// Add registers a sensitive value that should be redacted.
// Empty strings are ignored.
func (m *Masker) Add(secret string) {
	if secret == "" {
		return
	}
	m.secrets[secret] = struct{}{}
}

// Mask replaces all registered secret values in s with the redaction marker.
func (m *Masker) Mask(s string) string {
	for secret := range m.secrets {
		s = strings.ReplaceAll(s, secret, redacted)
	}
	return s
}

// MaskEnv takes a slice of KEY=VALUE environment strings and redacts any
// value whose key is in the provided set of sensitive keys.
func (m *Masker) MaskEnv(env []string, sensitiveKeys map[string]struct{}) []string {
	out := make([]string, len(env))
	for i, entry := range env {
		parts := strings.SplitN(entry, "=", 2)
		if len(parts) == 2 {
			if _, ok := sensitiveKeys[parts[0]]; ok {
				out[i] = parts[0] + "=" + redacted
				continue
			}
		}
		out[i] = entry
	}
	return out
}
