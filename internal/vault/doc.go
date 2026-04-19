// Package vault provides a thin wrapper around the HashiCorp Vault client
// for use by vaultenv. It exposes helpers for creating an authenticated
// client and reading KV v2 secrets by path.
//
// Secret paths are expected in the format "<mount>/<secret-path>", e.g.
// "secret/myapp/database". The mount point is the first path segment and
// must match the KV v2 engine mount configured in Vault.
package vault
