// Package labels provides helpers for attaching and filtering
// arbitrary key-value metadata to secret fetch requests.
package labels

import (
	"fmt"
	"strings"
)

// Labels is an immutable set of key=value metadata pairs.
type Labels map[string]string

// Parse parses a slice of "key=value" strings into a Labels map.
// It returns an error if any entry is malformed.
func Parse(pairs []string) (Labels, error) {
	l := make(Labels, len(pairs))
	for _, p := range pairs {
		key, val, ok := strings.Cut(p, "=")
		if !ok {
			return nil, fmt.Errorf("labels: missing '=' in pair %q", p)
		}
		key = strings.TrimSpace(key)
		if key == "" {
			return nil, fmt.Errorf("labels: empty key in pair %q", p)
		}
		l[key] = strings.TrimSpace(val)
	}
	return l, nil
}

// Get returns the value for key and whether it was present.
func (l Labels) Get(key string) (string, bool) {
	v, ok := l[key]
	return v, ok
}

// Matches reports whether all pairs in filter are present with equal
// values in l. An empty filter always matches.
func (l Labels) Matches(filter Labels) bool {
	for k, v := range filter {
		if l[k] != v {
			return false
		}
	}
	return true
}

// Merge returns a new Labels that is the union of l and other.
// Keys in other take precedence over keys in l.
func (l Labels) Merge(other Labels) Labels {
	out := make(Labels, len(l)+len(other))
	for k, v := range l {
		out[k] = v
	}
	for k, v := range other {
		out[k] = v
	}
	return out
}

// String returns a deterministic comma-separated representation.
func (l Labels) String() string {
	if len(l) == 0 {
		return ""
	}
	parts := make([]string, 0, len(l))
	for k, v := range l {
		parts = append(parts, k+"="+v)
	}
	// stable sort for deterministic output
	for i := 1; i < len(parts); i++ {
		for j := i; j > 0 && parts[j] < parts[j-1]; j-- {
			parts[j], parts[j-1] = parts[j-1], parts[j]
		}
	}
	return strings.Join(parts, ",")
}
