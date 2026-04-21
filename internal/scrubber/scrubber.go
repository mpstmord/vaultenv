// Package scrubber provides a middleware-compatible secret scrubber that
// rewrites log output and error messages to replace known secret values
// with a configurable placeholder before they leave the process.
package scrubber

import (
	"strings"
	"sync"
)

const defaultPlaceholder = "[REDACTED]"

// Scrubber holds a set of secret strings and replaces them in text.
type Scrubber struct {
	mu          sync.RWMutex
	secrets     []string
	placeholder string
}

// New returns a Scrubber with the given placeholder string.
// If placeholder is empty, "[REDACTED]" is used.
func New(placeholder string) *Scrubber {
	if placeholder == "" {
		placeholder = defaultPlaceholder
	}
	return &Scrubber{placeholder: placeholder}
}

// Add registers one or more secret values to be scrubbed.
// Empty strings are silently ignored.
func (s *Scrubber) Add(secrets ...string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, v := range secrets {
		if v != "" {
			s.secrets = append(s.secrets, v)
		}
	}
}

// Scrub replaces all registered secret values in text with the placeholder.
func (s *Scrubber) Scrub(text string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, secret := range s.secrets {
		text = strings.ReplaceAll(text, secret, s.placeholder)
	}
	return text
}

// ScrubMap returns a copy of m where every value that contains a registered
// secret has that secret replaced with the placeholder.
func (s *Scrubber) ScrubMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = s.Scrub(v)
	}
	return out
}

// Len returns the number of registered secrets.
func (s *Scrubber) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.secrets)
}
