// Package plugin defines the Provider interface and Registry for pluggable
// secret backends in vaultenv.
//
// A Provider retrieves a secret value given a path and field name. Providers
// are registered by name and looked up at runtime based on the mapping
// configuration supplied to the run or exec commands.
//
// Built-in providers:
//
//	env  – reads secrets from environment variables (field = variable name)
package plugin
