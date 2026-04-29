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
//
// Each Result in the returned slice corresponds positionally to the input
// Request slice. Callers should check Result.Err before using Result.Value:
//
//	for i, r := range results {
//		if r.Err != nil {
//			log.Printf("failed to fetch %s: %v", results[i].Path, r.Err)
//			continue
//		}
//		fmt.Println(r.Value)
//	}
//
// The worker pool size passed to New controls the maximum number of concurrent
// requests made to Vault. A value of 0 defaults to runtime.NumCPU().
package batch
