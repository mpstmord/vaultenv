// Package tags provides utilities for parsing and matching metadata tags
// attached to secret mappings, enabling fine-grained filtering and labelling.
package tags

import (
	"fmt"
	"sort"
	"strings"
)

// Tags holds a set of key=value metadata pairs.
type Tags map[string]string

// Parse parses a slice of "key=value" strings into a Tags map.
// It returns an error if any entry is malformed.
func Parse(pairs []string) (Tags, error) {
	t := make(Tags, len(pairs))
	for _, p := range pairs {
		k, v, ok := strings.Cut(p, "=")
		if !ok {
			return nil, fmt.Errorf("tags: missing '=' in pair %q", p)
		}
		k = strings.TrimSpace(k)
		if k == "" {
			return nil, fmt.Errorf("tags: empty key in pair %q", p)
		}
		t[k] = strings.TrimSpace(v)
	}
	return t, nil
}

// Get returns the value for key and whether it was present.
func (t Tags) Get(key string) (string, bool) {
	v, ok := t[key]
	return v, ok
}

// Set sets key to value.
func (t Tags) Set(key, value string) {
	t[key] = value
}

// Matches reports whether all entries in subset are present and equal in t.
func (t Tags) Matches(subset Tags) bool {
	for k, v := range subset {
		if t[k] != v {
			return false
		}
	}
	return true
}

// Keys returns all keys in sorted order.
func (t Tags) Keys() []string {
	keys := make([]string, 0, len(t))
	for k := range t {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// String returns a stable "k=v,..." representation.
func (t Tags) String() string {
	parts := make([]string, 0, len(t))
	for _, k := range t.Keys() {
		parts = append(parts, k+"="+t[k])
	}
	return strings.Join(parts, ",")
}
