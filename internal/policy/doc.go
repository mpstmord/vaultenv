// Package policy provides rule-based access control for Vault secret paths.
//
// Rules are evaluated in declaration order; the first matching rule
// determines whether access is allowed or denied. If no rule matches,
// access is denied by default.
//
// Patterns support a single "*" wildcard that matches any non-separator
// segment, and "**" which matches any path unconditionally.
package policy
