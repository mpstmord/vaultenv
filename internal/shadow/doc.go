// Package shadow implements a dual-read pattern for secret fetching.
//
// A Shadower wraps a primary and a secondary Fetcher. Every call to
// GetSecretData is forwarded to both sources concurrently. The result from the
// primary is always returned to the caller; the secondary result is used only
// for comparison. Any discrepancy between the two responses — or any error from
// the secondary — is written to the configured Logger.
//
// This is useful when migrating between secret backends: traffic can be
// shadowed to the new backend before cutting over, giving operators confidence
// that the new source returns equivalent data.
package shadow
