// Package namespace provides a Fetcher wrapper that scopes every secret path
// under a fixed prefix string. This enables multi-tenant and
// environment-aware secret isolation when multiple teams or deployment stages
// share a single Vault cluster.
//
// Example usage:
//
//	ns, err := namespace.New("prod/myteam", vaultClient)
//	if err != nil { ... }
//	data, err := ns.GetSecretData(ctx, "service/api-key")
//	// resolves to: prod/myteam/service/api-key
package namespace
