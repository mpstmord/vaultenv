// Package warmup resolves a set of Vault secret mappings concurrently
// before the child process is launched, ensuring all required environment
// variables are available at startup.
//
// Usage:
//
//	r := warmup.New(vaultClient, envStore)
//	if err := r.Run(ctx, mappings); err != nil {
//		log.Fatal(err)
//	}
package warmup
