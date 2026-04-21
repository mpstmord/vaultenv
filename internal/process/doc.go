// Package process provides utilities for executing child processes with
// environment variables injected from HashiCorp Vault secrets.
//
// It supports merging the current process environment with resolved secret
// values and spawning the target command with the combined environment.
//
// # Environment Merging
//
// When building the environment for a child process, the current process
// environment is used as a base. Secret values resolved from Vault are
// merged on top, with Vault-provided values taking precedence over any
// existing environment variables of the same name.
//
// # Process Execution
//
// The child process is executed using syscall.Exec (on Unix systems),
// which replaces the current process image. This means vaultenv itself
// does not remain resident once the child process starts.
package process
