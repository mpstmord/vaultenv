// Package process provides utilities for executing child processes with
// environment variables injected from HashiCorp Vault secrets.
//
// It supports merging the current process environment with resolved secret
// values and spawning the target command with the combined environment.
package process
