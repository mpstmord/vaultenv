// Package quota provides per-path request rate limiting for secret fetches.
//
// A Limiter tracks how many times each Vault secret path has been accessed
// within a sliding time window and rejects requests that exceed the configured
// limit. GuardedFetcher wraps any SecretFetcher with quota enforcement so
// callers do not need to manage limit checks explicitly.
package quota
