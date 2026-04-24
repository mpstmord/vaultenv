// Package policy evaluates access rules for secret paths.
package policy

import (
	"errors"
	"strings"
)

// Rule defines a single access control entry.
type Rule struct {
	// Path is a glob-style pattern, e.g. "secret/app/*".
	Path string
	// Allow indicates whether matching paths are permitted.
	Allow bool
}

// Policy holds an ordered list of rules evaluated top-to-bottom.
type Policy struct {
	rules []Rule
}

// New creates a Policy from the provided rules.
// Rules are evaluated in order; the first match wins.
// If no rule matches, access is denied by default.
func New(rules []Rule) (*Policy, error) {
	for _, r := range rules {
		if strings.TrimSpace(r.Path) == "" {
			return nil, errors.New("policy: rule path must not be empty")
		}
	}
	return &Policy{rules: rules}, nil
}

// Allow reports whether the given secret path is permitted by the policy.
func (p *Policy) Allow(path string) bool {
	for _, r := range p.rules {
		if matchGlob(r.Path, path) {
			return r.Allow
		}
	}
	return false
}

// Len returns the number of rules in the policy.
func (p *Policy) Len() int { return len(p.rules) }

// matchGlob performs simple glob matching where "*" matches any
// sequence of non-separator characters and "**" matches everything.
func matchGlob(pattern, s string) bool {
	if pattern == "**" {
		return true
	}
	parts := strings.SplitN(pattern, "*", 2)
	if len(parts) == 1 {
		return pattern == s
	}
	prefix, suffix := parts[0], parts[1]
	if !strings.HasPrefix(s, prefix) {
		return false
	}
	rest := s[len(prefix):]
	if suffix == "" {
		return !strings.Contains(rest, "/")
	}
	idx := strings.Index(rest, suffix)
	if idx < 0 {
		return false
	}
	return !strings.Contains(rest[:idx], "/")
}
