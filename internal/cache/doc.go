// Package cache implements a thread-safe, TTL-based in-memory cache
// for Vault secret payloads. It reduces redundant Vault API calls
// when multiple environment mappings reference the same secret path.
package cache
