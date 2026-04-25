// Package fanout provides two strategies for querying multiple secret
// fetchers:
//
//   - Fanout queries all fetchers concurrently and merges their results,
//     with later fetchers overriding keys from earlier ones.
//
//   - PriorityFetcher queries fetchers sequentially and returns the first
//     successful response, falling back to the next fetcher on error.
//
// Both types satisfy the same Fetcher interface used throughout vaultenv.
package fanout
