// Package batch provides concurrent, ordered fetching of multiple Vault secrets
// in a single call. It dispatches requests across a configurable pool of worker
// goroutines and collects results in the original request order.
//
// Basic usage:
//
//	b := batch.New(vaultClient, 8)
//	results := b.FetchAll(ctx, []batch.Request{
//		{Path: "secret/db", Field: "password"},
//		{Path: "secret/api", Field: "key"},
//	})
package batch
