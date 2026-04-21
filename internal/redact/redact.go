// Package redact provides utilities for scrubbing secret values from strings
// and structured log output before they are written to any sink.
package redact

import "strings"

// Redactor holds a set of sensitive values and can scrub them from text.
type Redactor struct {
	secrets []string
}

// New returns a Redactor loaded with the provided secret values.
// Empty strings are silently ignored.
func New(secrets []string) *Redactor {
	filtered := make([]string, 0, len(secrets))
	for _, s := range secrets {
		if s != "" {
			filtered = append(filtered, s)
		}
	}
	return &Redactor{secrets: filtered}
}

// Scrub replaces every occurrence of a known secret inside text with the
// placeholder string "[REDACTED]".
func (r *Redactor) Scrub(text string) string {
	for _, s := range r.secrets {
		text = strings.ReplaceAll(text, s, "[REDACTED]")
	}
	return text
}

// ScrubMap returns a copy of env where values matching any known secret are
// replaced with "[REDACTED]". Keys are never modified.
func (r *Redactor) ScrubMap(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = r.Scrub(v)
	}
	return out
}

// Add appends additional secret values to the redactor at runtime.
// Empty strings are silently ignored.
func (r *Redactor) Add(secrets ...string) {
	for _, s := range secrets {
		if s != "" {
			r.secrets = append(r.secrets, s)
		}
	}
}

// Len returns the number of secret values currently tracked by the Redactor.
func (r *Redactor) Len() int {
	return len(r.secrets)
}
