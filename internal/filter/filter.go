// Package filter provides utilities for selecting and excluding
// environment variables by key patterns before process injection.
package filter

import (
	"path"
	"strings"
)

// Filter holds include and exclude glob patterns for env key filtering.
type Filter struct {
	include []string
	exclude []string
}

// New creates a Filter with the given include and exclude glob patterns.
// An empty include list means all keys are included by default.
func New(include, exclude []string) *Filter {
	return &Filter{include: include, exclude: exclude}
}

// Allow reports whether the given env key should be passed through.
func (f *Filter) Allow(key string) bool {
	upper := strings.ToUpper(key)
	for _, pat := range f.exclude {
		if matchGlob(pat, upper) {
			return false
		}
	}
	if len(f.include) == 0 {
		return true
	}
	for _, pat := range f.include {
		if matchGlob(pat, upper) {
			return true
		}
	}
	return false
}

// Apply returns a new map containing only the entries allowed by the filter.
func (f *Filter) Apply(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if f.Allow(k) {
			out[k] = v
		}
	}
	return out
}

func matchGlob(pattern, key string) bool {
	matched, err := path.Match(strings.ToUpper(pattern), key)
	return err == nil && matched
}
