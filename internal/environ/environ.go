// Package environ provides utilities for capturing and diffing
// the current OS environment as a key-value map.
package environ

import (
	"os"
	"strings"
)

// Snapshot captures the current OS environment and returns it as a map.
func Snapshot() map[string]string {
	return fromSlice(os.Environ())
}

// fromSlice converts a slice of "KEY=VALUE" strings into a map.
// Entries without an "=" separator are stored with an empty value.
func fromSlice(pairs []string) map[string]string {
	m := make(map[string]string, len(pairs))
	for _, p := range pairs {
		key, val, _ := strings.Cut(p, "=")
		if key != "" {
			m[key] = val
		}
	}
	return m
}

// ToSlice converts a map back into a slice of "KEY=VALUE" strings
// suitable for use with os/exec.Cmd.Env.
func ToSlice(env map[string]string) []string {
	out := make([]string, 0, len(env))
	for k, v := range env {
		out = append(out, k+"="+v)
	}
	return out
}

// Merge returns a new map that is the union of base and override.
// Keys present in override take precedence over keys in base.
func Merge(base, override map[string]string) map[string]string {
	result := make(map[string]string, len(base)+len(override))
	for k, v := range base {
		result[k] = v
	}
	for k, v := range override {
		result[k] = v
	}
	return result
}

// Diff returns the keys whose values differ between old and next.
// It includes keys that were added, removed, or changed.
func Diff(old, next map[string]string) []string {
	seen := make(map[string]struct{})
	var changed []string

	for k, v := range next {
		seen[k] = struct{}{}
		if prev, ok := old[k]; !ok || prev != v {
			changed = append(changed, k)
		}
	}
	for k := range old {
		if _, ok := seen[k]; !ok {
			changed = append(changed, k)
		}
	}
	return changed
}
