// Package semaphore provides a counting semaphore for limiting concurrent
// access to shared resources such as Vault secret fetches.
//
// Use New to create a semaphore with a fixed capacity, then call Acquire
// before entering a critical section and Release when leaving.
//
// GuardedFetcher wraps any Fetcher implementation and transparently enforces
// a concurrency limit, making it easy to compose with the middleware chain.
package semaphore
