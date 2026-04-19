// Package env provides utilities for mapping HashiCorp Vault secret fields
// to environment variables. It includes:
//
//   - ParseMapping: parses "ENV_VAR=vault/path#field" strings into SecretMapping structs.
//   - Resolver: fetches secret values from Vault and resolves them into a
//     map of environment variable names to values.
//   - InjectIntoEnv: sets resolved key/value pairs as OS environment variables
//     so they are available to child processes launched by vaultenv.
package env
