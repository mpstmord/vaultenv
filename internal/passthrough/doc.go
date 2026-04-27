// Package passthrough implements an environment-variable override layer for
// secret fetching.
//
// When a matching environment variable is present, the upstream Vault fetch is
// skipped entirely and the env value is returned directly. This is useful for
// local development, CI pipelines, and testing scenarios where real Vault
// credentials are unavailable.
package passthrough
